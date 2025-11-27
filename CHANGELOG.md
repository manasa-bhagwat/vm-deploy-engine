# Changelog

All notable changes to this project will be documented in this file.

The format follows **Keep a Changelog** and this project adheres to **Semantic Versioning**.

---

## [v2.0.4] - 2025-11-27
### Added
- Introduced complete deployment orchestration pipeline in `internal/deploy/steps.go`.
- Added modular operators for:
  - Base package installation (`Installer`)
  - Multi-DB provisioning (MySQL, PostgreSQL, MongoDB) via `DBInstaller`
  - Repository cloning using Git with PAT injection (`GitOps`)
  - Maven-based application build (`MavenOps`)
  - systemd environment file + service generation (`SystemdOps`)
- Implemented `Deployer.RunFullDeployment()` to execute full install → clone → build → configure → service lifecycle.
- Added key-based SSH authentication support:
  - `SSHConfig.SSHKeyPath`
  - `SSHConfig.SSHPassphrase`
  - `SSHConfig.UseSSHAgent`
- Implemented SSH private key parsing:
  - Unencrypted key parsing
  - Encrypted key parsing (`ParsePrivateKeyWithPassphrase`)
  - ssh-agent signer discovery (`SSH_AUTH_SOCK`)
- Improved VM config structure (`vmconfig.yaml`) to support modern infra-style SSH workflows.

### Changed
- Migrated SSH authentication away from password-based auth to key-based workflows.
- Updated `main.go` to:
  - Prompt only for non-secret runtime inputs (DB creds + PAT)
  - Build new SSH client with key-based config
  - Invoke full deployment pipeline instead of standalone commands
- Updated `appconfig.yaml` and `vmconfig.yaml` loading logic to align with new orchestrator.

### Fixed
- Ensured correct remote Maven permissions by guaranteeing `/opt/app` creation and ownership.
- Fixed remote execution output formatting and handling of empty stderr/stdout streams.

### Removed
- Removed legacy password-based SSH authentication paths.

---

## [v2.0.3] - 2025-11-22
### Added
- Introduced initial `vmdeploy` CLI with the `deploy` command.
- Added separate `appconfig.yaml` and `vmconfig.yaml` loading.
- Added SSH password prompt (not stored).
- Implemented SSH ping + connect workflow.
- Introduced first remote execution command (`uname -a`) for connectivity validation.
- Added remote output formatting and improved CLI messages.

### Changed
- Refined project directory layout (`cmd/`, `internal/`, `pkg/`).
- Updated deployment entrypoint to consume new config modules.

---

## [v2.0.2] - 2025-11-22
### Added
- Added `pkg/sshutil` package with:
  - `SSHConfig` structure
  - `SSHClient` constructor
  - Basic `Ping()` method for TCP reachability
  - `Connect()` establishing SSH handshake
- Introduced internal structure for future multi-node deployment.

---

## [v2.0.1] - 2025-11-21
### Added
- Created initial CLI skeleton under `cmd/vmdeploy`.
- Added argument parsing and top-level command routing.
- Implemented `vmdeploy deploy` entrypoint handler.
- Added helper for interactive input (`readLine`).

---

## [v2.0.0] - 2025-11-21
### Added
- Initialized Go module:  
  `go mod init github.com/manasa-bhagwat/vm-deploy-engine/v2`
- First versioned folder `v2.0/` added for Go rewrite of the Bash engine.
- Created initial project structure based on CNCF-style separation.

---

## [v1.1.0] - 2025-11-20
### Added
- Config-driven deployment using `/etc/app-deploy.conf`.
- Support for MySQL, PostgreSQL, MongoDB provisioning.
- Secure prompts for GitHub PAT and DB credentials.
- DB readiness detection across all supported engines.
- Normalized artifact handling and improved script structure.
- Enforced safe bash mode (`set -euo pipefail`).

---

## [v1.0.0] - 2025-11-18
### Added
- Baseline VM deployment script.
- Automated installation of Java, MySQL, Git.
- MySQL provisioning + readiness checks.
- Repo cloning, Maven build, and JAR normalization.
- systemd service generation and lifecycle wiring.
- Fully automated end-to-end Spring Boot deployment.
