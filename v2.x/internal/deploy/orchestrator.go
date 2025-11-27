package deploy

import (
	"fmt"

	"github.com/manasa-bhagwat/vm-deploy-engine/v2/internal/config"
	"github.com/manasa-bhagwat/vm-deploy-engine/v2/pkg/sshutil"
)

// Deployer is the top-level orchestrator for the VM deployment flow
type Deployer struct {
	Client *sshutil.SSHClient
}

// RunFullDeployment runs ALL deployment steps in correct order.
// Think of this like an Ansible playbook executed step-by-step.
func (d *Deployer) RunFullDeployment(app *config.AppConfig, dbUser, dbPass, pat string) error {

	// instantiate all modules
	inst := Installer{Client: d.Client}
	gitOps := GitOps{Client: d.Client}
	mvn := MavenOps{Client: d.Client}
	sd := SystemdOps{Client: d.Client}
	dbInst := DBInstaller{Client: d.Client}

	// ========== STEP 1: Install Base Packages ==========
	fmt.Println("[STEP] Installing base packages (Java + Git + curl)...")
	if err := inst.InstallBasePackages(); err != nil {
		return err
	}

	// ========== STEP 2: Install & Configure Database ==========
	fmt.Printf("[STEP] Installing & provisioning database (%s)...\n", app.DBType)
	if err := dbInst.InstallAndProvision(app, dbUser, dbPass); err != nil {
		return err
	}

	fmt.Println("[STEP] Preparing /opt/app permissions...")

	if _, _, err := d.Client.Run("sudo mkdir -p /opt/app"); err != nil {
		return fmt.Errorf("failed to create /opt/app: %w", err)
	}

	if _, _, err := d.Client.Run("sudo chown -R ubuntu:ubuntu /opt/app"); err != nil {
		return fmt.Errorf("failed to change ownership for /opt/app: %w", err)
	}

	// ========== STEP 3: Clone Git Repository ==========
	fmt.Println("[STEP] Cloning application repository...")
	if err := gitOps.CloneAsUser(app.RepoURL, app.Branch, "/opt/app", pat, "ubuntu"); err != nil {
		return err
	}

	// ========== STEP 4: Build App (Maven Wrapper) ==========
	fmt.Println("[STEP] Building application via Maven...")
	if err := mvn.Build("/opt/app", app.ArtifactName); err != nil {
		return err
	}

	// ========== STEP 5: Build DB URL based on DB Type ==========
	dbURL, err := buildDBURL(app, dbUser, dbPass)
	if err != nil {
		return err
	}

	// ========== STEP 6: Write Env File ==========
	envPath := "/etc/" + app.ServiceName + ".env"
	fmt.Println("[STEP] Writing environment file:", envPath)
	if err := sd.WriteEnvFile(envPath, dbURL, dbUser, dbPass, app.AppPort); err != nil {
		return err
	}

	// ========== STEP 7: Systemd Service ==========
	servicePath := "/etc/systemd/system/" + app.ServiceName + ".service"
	fmt.Println("[STEP] Writing systemd service:", servicePath)
	if err := sd.WriteService(servicePath, app.ServiceName, "/opt/app", app.ArtifactName); err != nil {
		return err
	}

	// ========== STEP 8: Restart Service ==========
	fmt.Println("[STEP] Reloading systemd & restarting service...")
	if err := sd.ReloadAndRestart(app.ServiceName); err != nil {
		return err
	}

	return nil
}

// buildDBURL generates the correct DB connection URL based on the database type.
func buildDBURL(app *config.AppConfig, user, pass string) (string, error) {
	switch app.DBType {
	case "mysql":
		return fmt.Sprintf("jdbc:mysql://%s:%d/%s", app.DBHost, app.DBPort, app.DBName), nil

	case "postgres":
		return fmt.Sprintf("jdbc:postgresql://%s:%d/%s", app.DBHost, app.DBPort, app.DBName), nil

	case "mongodb":
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
			user, pass, app.DBHost, app.DBPort, app.DBName,
		), nil

	default:
		return "", fmt.Errorf("unsupported DB type: %s", app.DBType)
	}
}
