package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	addCommand(newLintCommand)
}

type lintCommand struct {
}

func (lintCommand) lint(args []string) error {
	cmd := exec.Command("golangci-lint", "run", "./...")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
func newLintCommand() *cobra.Command {
	c := &lintCommand{}
	cmd := &cobra.Command{
		Use:   "lint",
		Short: "lint contract code",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.lint(os.Args)
		},
	}
	return cmd
}
