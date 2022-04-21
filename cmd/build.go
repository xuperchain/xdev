package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xuperchain/xdev/mkfile"
)

var (
	defaultCxxFlags = []string{
		"-std=c++11",
		"-I/usr/local/include",
		"-Isrc",
		"-Werror=vla",
	}
	defaultLDFlags = []string{
		"-s TOTAL_STACK=256KB",
		"-s DETERMINISTIC=1",
		"-s EXPORTED_RUNTIME_METHODS=[\"stackAlloc\"]",
		"-L/usr/local/lib",
		"-L/emsdk/upstream/emscripten/cache/sysroot/lib/",
		"-lprotobuf-lite",
		"-lpthread",
	}
)

const (
	buildModeDebug   = "debug"
	buildModeRelease = "release"
)

var (
	debugBuildFlags = []string{
		"-fsanitize=undefined",
		"-fsanitize=address",
		"-O0",
	}
	debugLinkFlags = []string{
		"-fsanitize=undefined",
		"-fsanitize=address",

		"-s TOTAL_MEMORY=64MB",
		"-O0",
		"-s ALLOW_MEMORY_GROWTH=1",
		"-s MAXIMUM_MEMORY=128MB"}
)

var (
	releaseBuildFlags = []string{"-Os"}
	releaseLinkFlags  = []string{"-s TOTAL_MEMORY=1MB", "-Oz"}
)
var (
	ccImageRelease = "xuper/emcc:0.1.0"
	ccImageDebug   = "xuper/emcc:llvm_backend"
)

type buildCommand struct {
	cxxFlags            []string
	ldflags             []string
	builder             *mkfile.Builder
	entryPkg            *mkfile.Package
	UsingPrecompiledSDK bool
	NoEntry             bool
	xdevRoot            string
	buildMod            string
	ccImage             string

	genCompileCommand bool
	makeFileOnly      bool
	output            string
	compiler          string
	makeFlags         string
	submodules        []string
}

func newBuildCommand() *cobra.Command {
	c := &buildCommand{
		ldflags:  defaultLDFlags,
		cxxFlags: defaultCxxFlags,
	}

	cmd := &cobra.Command{
		Use:   "build",
		Short: "build command builds a project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if c.UsingPrecompiledSDK {
				c.ldflags = append(c.ldflags, fmt.Sprintf("-L%s/lib", mkfile.DefaultXROOT), "-lxchain", "-lprotobuf-lite")
				c.ldflags = append(c.ldflags, fmt.Sprintf("--js-library %s/src/xchain/exports.js", mkfile.DefaultXROOT))
				c.cxxFlags = append(c.cxxFlags, fmt.Sprintf("-I%s/src", mkfile.DefaultXROOT))
			} else {
				xroot := os.Getenv("XDEV_ROOT")
				c.xdevRoot = xroot
				c.ldflags = append(c.ldflags, fmt.Sprintf("--js-library %s/src/xchain/exports.js", xroot))
			}
			if c.NoEntry {
				c.ldflags = append(c.ldflags, "--no-entry")
			}
			// CCImage 优先级：环境变量 > 默认值
			// 1. 如果是debug 模式，则采用debugImage
			// 2. 如果有环境变量，则以环境变量为准

			c.ccImage = ccImageRelease
			if c.buildMod == buildModeDebug {
				c.ccImage = ccImageDebug
			}

			if image := os.Getenv("XDEV_CC_IMAGE"); image != "" {
				c.ccImage = image
			}

			if c.buildMod == buildModeDebug {
				c.cxxFlags = append(c.cxxFlags, debugBuildFlags...)
				c.ldflags = append(c.ldflags, debugLinkFlags...)
			} else if c.buildMod == buildModeRelease {
				c.cxxFlags = append(c.cxxFlags, releaseBuildFlags...)
				c.ldflags = append(c.ldflags, releaseLinkFlags...)
			}
			return c.build(args)
		},
	}
	cmd.Flags().BoolVarP(&c.makeFileOnly, "makefile", "m", false, "generate makefile and exit")
	cmd.Flags().BoolVarP(&c.genCompileCommand, "compile_command", "p", false, "generate compile_commands.json for IDE")
	cmd.Flags().StringVarP(&c.output, "output", "o", "", "output file name")
	cmd.Flags().StringVarP(&c.compiler, "compiler", "", "docker", "compiler env docker|host")
	cmd.Flags().StringVarP(&c.makeFlags, "mkflags", "", "", "extra flags passing to make command")
	cmd.Flags().StringSliceVarP(&c.submodules, "submodule", "s", nil, "build submodules")
	cmd.Flags().BoolVarP(&c.UsingPrecompiledSDK, "using-precompiled-sdk", "", true, "using precompiled sdk")
	cmd.Flags().BoolVarP(&c.NoEntry, "no-entry", "", true, "do not output any entry point")
	cmd.Flags().StringVarP(&c.buildMod, "build-mode", "", buildModeRelease, "build mode, may be debug or release")
	// cmd.Flags().StringVarP(&c.ccImage, "cc-image", "", ccImageRelease, "")
	return cmd
}

func (c *buildCommand) parsePackage(root, xcache string) error {
	absroot, err := filepath.Abs(root)
	if err != nil {
		return err
	}

	addons, err := c.addonModules(absroot)
	if err != nil {
		return err
	}
	if c.submodules != nil {
		addons = append(addons, mkfile.DependencyDesc{
			Name:    "self",
			Modules: c.submodules,
		})
	}

	loader := mkfile.NewLoader().WithXROOT(c.xdevRoot)
	pkg, err := loader.Load(absroot, addons)
	if err != nil {
		return err
	}

	output := c.output
	// 如果没有指定输出，且为main package，则用package目录名+wasm后缀作为输出名字
	if output == "" && pkg.Name == mkfile.MainPackage {
		output = filepath.Base(absroot) + ".wasm"
	}

	if output != "" {
		c.output, err = filepath.Abs(output)
		if err != nil {
			return err
		}
	}

	b := mkfile.NewBuilder().
		WithCxxFlags(c.cxxFlags).
		WithLDFlags(c.ldflags).
		WithCacheDir(xcache).
		WithOutput(c.output)

	err = b.Parse(pkg)
	if err != nil {
		return err
	}
	c.builder = b
	c.entryPkg = pkg
	return nil
}

func (c *buildCommand) xdevCacheDir() (string, error) {
	xcache := os.Getenv("XDEV_CACHE")
	if xcache != "" {
		return filepath.Abs(xcache)
	}
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homedir, ".xdev-cache"), nil
}

func (c *buildCommand) build(args []string) error {
	var err error
	if c.output != "" && !filepath.IsAbs(c.output) {
		c.output, err = filepath.Abs(c.output)
		if err != nil {
			return err
		}
	}

	if len(args) == 0 {
		root, err := findPackageRoot()
		if err != nil {
			return err
		}
		return c.buildPackage(root)
	}

	return c.buildFiles(args)
}

func xchainModule(xroot string) mkfile.DependencyDesc {
	return mkfile.DependencyDesc{
		Name: "xchain",
		Path: xroot,
		Modules: []string{
			"xchain",
		},
	}
}

func (c *buildCommand) addonModules(pkgpath string) ([]mkfile.DependencyDesc, error) {
	desc, err := mkfile.ParsePackageDesc(pkgpath)
	if err != nil {
		return nil, err
	}
	if desc.Package.Name != mkfile.MainPackage {
		return nil, nil
	}
	if !c.UsingPrecompiledSDK {
		return []mkfile.DependencyDesc{xchainModule(c.xdevRoot)}, nil
	}
	return []mkfile.DependencyDesc{}, nil
}

func (c *buildCommand) buildPackage(root string) error {
	wd, _ := os.Getwd()
	err := os.Chdir(root)
	if err != nil {
		return err
	}
	defer os.Chdir(wd)
	xcache, err := c.xdevCacheDir()
	if err != nil {
		return err
	}

	err = os.MkdirAll(xcache, 0755)
	if err != nil {
		return err
	}

	err = c.parsePackage(".", xcache)
	if err != nil {
		return err
	}

	if c.makeFileOnly {
		return c.builder.GenerateMakeFile(os.Stdout)
	}

	if c.genCompileCommand {
		cfile, err := os.Create("compile_commands.json")
		if err != nil {
			return err
		}
		defer cfile.Close()
		return c.builder.GenerateCompileCommands(cfile)
	}

	makefile, err := os.Create(".Makefile")
	if err != nil {
		return err
	}
	err = c.builder.GenerateMakeFile(makefile)
	if err != nil {
		makefile.Close()
		return err
	}
	makefile.Close()
	defer os.Remove(".Makefile")

	runner := mkfile.NewRunner(c.ccImage).
		WithEntry(c.entryPkg).
		WithCacheDir(xcache).
		WithXROOT(c.xdevRoot).
		WithOutput(c.output).
		WithMakeFlags(strings.Fields(c.makeFlags)).
		WithLogger(logger)

	if c.compiler != "docker" {
		runner = runner.WithoutDocker()
	}

	if !c.UsingPrecompiledSDK {
		runner = runner.WithoutPrecompiledSDK()
	}

	err = runner.Make(".Makefile")
	if err != nil {
		return err
	}
	return nil
}

func convertWasmFileName(fname string) string {
	idx := strings.LastIndex(fname, ".")
	if idx == -1 {
		return fname + ".wasm"
	}
	return fname[:idx] + ".wasm"
}

// 拷贝文件构造一个工程的目录结构，编译工程
func (c *buildCommand) buildFiles(files []string) error {
	basedir, err := ioutil.TempDir("", "xdev-build")
	if err != nil {
		return err
	}
	defer os.RemoveAll(basedir)

	if c.output == "" {
		c.output, err = filepath.Abs(convertWasmFileName(filepath.Base(files[0])))
		if err != nil {
			return err
		}
	}

	pkgDescFile := filepath.Join(basedir, mkfile.PkgDescFile)
	err = ioutil.WriteFile(pkgDescFile, []byte(`[package]
	name = "main"
	`), 0644)
	if err != nil {
		return err
	}

	srcdir := filepath.Join(basedir, "src")
	err = os.Mkdir(srcdir, 0755)
	if err != nil {
		return err
	}

	for _, file := range files {
		destfile := filepath.Join(srcdir, filepath.Base(file))
		err = cpfile(destfile, file)
		if err != nil {
			return err
		}
	}

	return c.buildPackage(basedir)
}

func cpfile(dest, src string) error {
	srcf, err := os.Open(src)
	if err != nil {
		return err
	}

	destf, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destf.Close()

	_, err = io.Copy(destf, srcf)
	return err
}

func findPackageRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if wd == "/" {
			return "", errors.New("can't find " + mkfile.PkgDescFile)
		}
		xcfile := filepath.Join(wd, mkfile.PkgDescFile)
		if _, err := os.Stat(xcfile); err == nil {
			return wd, nil
		}
		wd = filepath.Dir(wd)
	}
}

func init() {
	addCommand(newBuildCommand)
}
