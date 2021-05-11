package main

import "github.com/xuperchain/xdev/lib/cmd"

var (
	buildVersion = ""
	buildDate    = ""
	commitHash   = ""
)

func main() {
	cmd.SetVersion(buildVersion, buildDate, commitHash)
	cmd.Main()
}
