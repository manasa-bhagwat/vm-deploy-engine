package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("vmdeploy: missing command.")
		fmt.Println("Usage: vmdeploy <command>")
	}

	switch os.Args[1] {
	case "deploy":
		handleDeploy()
	default:
		fmt.Println("Unknown command: ", os.Args[1])
		os.Exit(1)
	}
}

func handleDeploy() {
	fmt.Println("Starting deployment (v2.0) ...")
}
