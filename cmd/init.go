package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
)

var descTpl = `[package]
name = "main"
`

var codeTpl = `#include "xchain/xchain.h"

struct Hello : public xchain::Contract {};

DEFINE_METHOD(Hello, initialize) {
    xchain::Context* ctx = self.context();
    ctx->ok("initialize succeed");
}

DEFINE_METHOD(Hello, hello) {
    xchain::Context* ctx = self.context();
    ctx->ok("hello world");
}
`

var testTpl = `
var assert = require("assert");

Test("hello", function (t) {
    var contract;
    t.Run("deploy", function (tt) {
        contract = xchain.Deploy({
            name: "hello",
            code: "../hello.wasm",
            lang: "c",
            init_args: {}
        })
    });

    t.Run("invoke", function (tt) {
        resp = contract.Invoke("hello", {});
        assert.equal(resp.Body, "hello world");
    })
})
`

type initCommand struct {
	lang     string
	contract string
}

func newInitCommand() *cobra.Command {
	c := &initCommand{}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init initializes a new project",
		RunE: func(cmd *cobra.Command, args []string) error {
			var root string
			//if len(args) == 1 {
			//	root = args[0]
			//}
			return c.init(root)
		},
	}
	cmd.Flags().StringVarP(&c.lang, "lang", "q", "cpp", "language of contract")
	cmd.Flags().StringVarP(&c.contract, "contract", "c", "counter", "contract name")
	return cmd
}

func (c *initCommand) init(root string) error {
	reppURL := fmt.Sprintf("https://hub.fastgit.org/xuperchain/contract-sdk-%s.git", c.lang)
	dest := os.TempDir()
	defer os.RemoveAll(dest)
	if _, err := os.Stat(dest); err != nil {
		//no error check
		os.RemoveAll(dest)
	}

	cmd := exec.Command("git", "clone", reppURL, dest+"/contracts")
	fmt.Println(filepath.Base("git") == "git")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	//no cp syscall exist,so just use cp command
	//TODO
	// use filepath
	// process cc contract
	contract := c.contract
	if c.lang == "cpp" {
		contract += ".cc"
	}
	cmd2 := exec.Command("cp", "-r", dest+"contracts/example/"+contract, ".")
	cmd2.Stderr = os.Stderr
	cmd2.Stdout = os.Stdout
	cmd2.Start()
	return cmd2.Wait()
	//if root != "" {
	//	err := os.MkdirAll(root, 0755)
	//	if err != nil {
	//		return err
	//	}
	//	os.Chdir(root)
	//}
	//pkgfile := mkfile.PkgDescFile
	//if _, err := os.Stat(pkgfile); err == nil {
	//	return errors.New("project already initialized")
	//}
	//err := ioutil.WriteFile(pkgfile, []byte(descTpl), 0644)
	//if err != nil {
	//	return err
	//}
	//maindir := filepath.Join("src")
	//err = os.MkdirAll(maindir, 0755)
	//if err != nil {
	//	return err
	//}
	//mainfile := filepath.Join(maindir, "main.cc")
	//err = ioutil.WriteFile(mainfile, []byte(codeTpl), 0644)
	//if err != nil {
	//	return err
	//}
	//
	//testdir := filepath.Join("test")
	//err = os.MkdirAll(testdir, 0755)
	//if err != nil {
	//	return err
	//}
	//testfile := filepath.Join(testdir, "hello.test.js")
	//err = ioutil.WriteFile(testfile, []byte(testTpl), 0644)
	//if err != nil {
	//	return err
	//}
	//return nil
}

func init() {
	addCommand(newInitCommand)
}
