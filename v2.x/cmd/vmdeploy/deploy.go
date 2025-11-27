package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"

	"github.com/manasa-bhagwat/vm-deploy-engine/v2/internal/config"
	"github.com/manasa-bhagwat/vm-deploy-engine/v2/internal/deploy"
	"github.com/manasa-bhagwat/vm-deploy-engine/v2/pkg/sshutil"
	"github.com/spf13/cobra"
)

// deployCmd implements: vmdeploy deploy
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy application to a remote VM over SSH",
	Long: `Deploys the configured application to a remote VM.

It will:
  - Load app & VM configs
  - Connect to the VM over SSH (agent/key)
  - Run a remote connectivity check
  - Execute the full deployment pipeline (install, clone, build, systemd).`,
	RunE: runDeploy, // Cobra will call this when "vmdeploy deploy" is executed.
}

// runDeploy contains the end-to-end deployment logic.
func runDeploy(cmd *cobra.Command, args []string) error {
	fmt.Println("=== vmdeploy v2.0.5 :: deploy ===")

	// 1) Load AppConfig
	appCfg, err := config.LoadAppConfig(appConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load app config from %s: %w", appConfigPath, err)
	}
	fmt.Println("[INFO] App config loaded:", appCfg.AppName)

	// 2) Load VMConfig
	vmCfg, err := config.LoadVMConfig(vmConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load vm config from %s: %w", vmConfigPath, err)
	}
	fmt.Println("[INFO] VM config loaded.")

	// 3) Prompt for secrets (DB + GitHub PAT)
	dbUser := readLine("DB Username: ")
	dbPass := readSecret("DB Password: ")
	githubPAT := readSecret("GitHub PAT: ")

	// 4) Build SSHConfig from VMConfig
	sshCfg := sshutil.SSHConfig{
		Host:          vmCfg.Host,
		Port:          vmCfg.Port,
		User:          vmCfg.User,
		SSHKeyPath:    vmCfg.SSHKeyPath,
		SSHPassphrase: vmCfg.SSHPassphrase,
		UseSSHAgent:   vmCfg.UseSSHAgent,
	}

	client := sshutil.NewSSHClient(sshCfg)

	// 5) Ping
	fmt.Println("[INFO] Pinging VM over TCP...")
	if err := client.Ping(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	// 6) Connect over SSH
	fmt.Println("[INFO] Connecting via SSH...")
	if err := client.Connect(); err != nil {
		return fmt.Errorf("ssh connect failed: %w", err)
	}
	defer client.Close()

	// 7) Remote connectivity test
	fmt.Println("[INFO] Checking remote system -> running 'uname -a'")

	stdout, stderr, err := client.Run("uname -a")

	fmt.Println("\n----- Remote Output -----")
	if strings.TrimSpace(stdout) != "" {
		fmt.Println("Output:", stdout)
	}
	if strings.TrimSpace(stderr) != "" {
		fmt.Println("Error Output:", stderr)
	}
	fmt.Println("-------------------------")

	if err != nil {
		return fmt.Errorf("remote check failed: %w", err)
	}

	fmt.Println("[INFO] Remote system is reachable and responding.")

	// 8) Full deployment pipeline
	deployer := deploy.Deployer{Client: client}

	fmt.Println("[INFO] Running full deployment pipeline...")
	if err := deployer.RunFullDeployment(appCfg, dbUser, dbPass, githubPAT); err != nil {
		return fmt.Errorf("deployment failed: %w", err)
	}

	fmt.Println("[SUCCESS] Deployment completed successfully.")
	return nil
}

// It's used for lightweight prompts (DB user, DB pass, PAT, etc.).
func readLine(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("[ERROR] input error:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(text)
}

func readSecret(prompt string) string {
	fmt.Print(prompt)
	byteValue, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("[ERROR] failed to read secret:", err)
		os.Exit(1)
	}
	fmt.Println() // move to next line after input
	return strings.TrimSpace(string(byteValue))
}
