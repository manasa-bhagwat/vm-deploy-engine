package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/manasa-bhagwat/vm-deploy-engine/v2/internal/config"
	"github.com/manasa-bhagwat/vm-deploy-engine/v2/pkg/sshutil"
)

// main is the entrypoint of the vmdeploy CLI library
func main() {

	if len(os.Args) < 2 {
		fmt.Println("vmdeploy: missing command.")
		fmt.Println("Usage: vmdeploy <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "deploy":
		handleDeploy()
	default:
		fmt.Println("Unknown command: ", os.Args[1])
		os.Exit(1)
	}
}

// handleDeploy is the function that handles the "vmdeploy deploy" command.
func handleDeploy() {
	fmt.Println("Starting deployment (v2.0) ...")

	// Load AppConfig.yaml
	appCfg, err := config.LoadAppConfig("appconfig.yaml")
	if err != nil {
		fmt.Println("[ERROR] AppConfig: ", err)
		os.Exit(1)
	}
	fmt.Println("[INFO] App config loaded: ", appCfg.AppName)

	// Load VMConfig.yaml
	vmCfg, err := config.LoadVMConfig("vmconfig.yaml")
	if err != nil {
		fmt.Println("[ERROR] VMConfig: ", err)
		os.Exit(1)
	}
	fmt.Println("[INFO] VM config loaded: ", vmCfg.Host)

	// Prompt for SSH password(not hidden yet - will be fixed later)
	sshPass := readLine("SSH Password (visible for now): ")

	// Build SSHConfig
	sshCfg := sshutil.SSHConfig{
		Host:     vmCfg.Host,
		Port:     vmCfg.Port,
		User:     vmCfg.User,
		Password: sshPass,
	}

	client := sshutil.NewSSHClient(sshCfg)

	// Ping
	fmt.Println("[INFO] Pinging VM...")
	if err := client.Ping(); err != nil {
		fmt.Println("[ERROR] Ping failed:", err)
		os.Exit(1)
	}

	// Connect
	fmt.Println("[INFO] Connecting via SSH...")
	if err := client.Connect(); err != nil {
		fmt.Println("[ERROR] SSH connect failed:", err)
		os.Exit(1)
	}
	defer client.Close()

	// Remote test
	fmt.Println("[INFO] Running remote connectivity check -> uname -a")
	stdout, stderr, err := client.Run("uname -a")

	fmt.Println("\n================ REMOTE OUTPUT ================")

	if strings.TrimSpace(stdout) != "" {
		fmt.Println("[REMOTE STDOUT]: ", stdout)
	}
	if strings.TrimSpace(stderr) != "" {
		fmt.Println("[REMOTE STDERR]:", stderr)
	}
	fmt.Println("===============================================\n")

	if err != nil {
		fmt.Println("[ERROR] Remote command failed:", err)
		os.Exit(1)
	}

	fmt.Println("[INFO] Remote system information received successfully.")

}

func readLine(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("[ERROR] input error: ", err)
		os.Exit(1)
	}
	return strings.TrimSpace(text)
}
