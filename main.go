package main

import (
	"fmt"
	"os"
)

const version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("agent-builder %s\n", version)
		return
	}

	fmt.Println("Hello from Agent Builder!")
	fmt.Println("ADK Multi-Agent Builder CLI")
	fmt.Println()
	fmt.Println("This tool helps you build ADK (Agent Development Kit) agents.")
	fmt.Println()
	fmt.Println("Run 'agent-builder version' to see the version.")
}
