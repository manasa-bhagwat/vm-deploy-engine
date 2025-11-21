package main

import (
	"fmt"
	"os"

	"github.com/manasa-bhagwat/vm-deploy-engine/v2/internal/config"
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

	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Println("Config error: ", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded config for app: %s\n", cfg.AppName)
}
