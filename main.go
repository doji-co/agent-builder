package main

import "github.com/doji-co/agent-builder/cmd"

var version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
