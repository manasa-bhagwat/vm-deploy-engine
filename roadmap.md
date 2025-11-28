# vm-deploy-engine — Roadmap

This document tracks the evolution of the tool across versions.

---

## 🎯 High-Level Goals
- Understand how real infra tools are designed
- Build a universal VM deployment engine
- Grow into a Go-based CLI (CNCF-style architecture)
- Add support for multiple runtimes and databases
- Add SSH deploy mode
- Add rollback + audit logs
- Eventually explore GitOps-style reconciliation

---

## v1.x — Bash Era (Learning Foundations)

### ✔ v1.0 — Baseline Deployment
- Java installation
- MySQL setup
- DB creation
- Repo cloning
- Maven build
- systemd service creation
- Basic readiness checks

### ✔ v1.1 — Config-Driven Multi-DB
- Config file added
- Secure prompts (DB + PAT)
- MySQL/Postgres/MongoDB
- Better safety
- Cleaner logs & permissions

---

## v2.x — Go CLI Era (Real Tooling Begins)

### v2.0 — CLI Scaffolding
- Implement CLI structure (`cmd/vmde`)
- Basic deploy command
- Logging module
- Config loader (YAML/JSON)
- Modular folder layout

### v2.1 — Java runtime builder
- Detect Java project
- Build using mvnw/gradle
- Produce deployable artifact

### v2.2 — Node.js runtime builder
- Detect Node project
- npm/yarn install
- pm2/systemd support

### v2.3 — Python runtime builder
- Detect Python project
- virtualenv builder
- gunicorn/uvicorn runner

### v2.4 — .NET builder
- dotnet install
- dotnet publish
- systemd integration

### v2.5 — SSH Deploy Mode
- Deploy to remote VMs
- Parallel deploy support
- Timeout/retry logic

---

##  v3.x — Secrets & Vault-Lite
- Encrypted secrets
- Local keyring integration
- Encrypted config files

---

## v4.x — GitOps & Reconciliation Lite
- State file
- Drift detection
- Auto-sync mode

---

## Long-Term Ideas
- Web UI dashboard
- Buildpack-inspired architecture
- Plugin model for runtimes
- Kubernetes operator (far future)