# Changelog

All notable changes to this project will be documented in this file.

The format follows **Keep a Changelog** and this project adheres to **Semantic Versioning**.

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
