# VM Deploy Engine

A lightweight, CNCF-inspired deployment engine that provisions Linux VMs, installs runtimes, configures databases, builds applications, and creates systemd-managed services — all from a predictable CLI workflow.

This project is not just a tool — it is a **learning ground** for understanding how infrastructure automation works behind the scenes.

---

# Why This Project Exists

Before using Terraform, Ansible, Kubernetes, or ArgoCD…  
I wanted to understand:

- **How does a VM get prepared for deployment?**  
- **How does a service actually start under systemd?**  
- **How do tools provision databases under the hood?**  
- **How does SSH automation really work?**  
- **How do internal deployment platforms evolve over time?**

The **VM Deploy Engine** is my way of rebuilding those fundamentals from scratch, intentionally and systematically.

This repo now contains:

- **v1.x** → Shell-based deploy engine (complete lifecycle automation)  
- **v2.x** → Go-based CLI rewrite (extensible, testable, platform-grade)

---

# Roadmap Overview

### ✔ **v1.0.0 – Baseline Automation Script**
Shell script that performs:

- Install Java, Git, MySQL
- Create DB + user
- Clone repo using PAT
- Build Spring Boot app with Maven Wrapper
- Normalize JAR
- Generate env file
- Generate systemd service  
- Launch the service predictably

---

### ✔ **v1.1.0 – Secure Config-Driven Multi-DB Engine**
Improved Bash engine:

- Central configuration file (`/etc/app-deploy.conf`)
- Supports **MySQL**, **PostgreSQL**, **MongoDB**
- Hidden prompts for PAT + DB credentials
- DB readiness checks
- Root-owned env file with 600 permissions  
- Strict Bash safety — `set -euo pipefail`

---

### ✔ **v2.0.x – Go CLI Rewrite (CNCF Style)**  
Rebuilt the entire engine using:

- Structured Go module layout (`cmd/`, `internal/`, `pkg/`)
- YAML configs (`appconfig.yaml`, `vmconfig.yaml`)
- SSH automation via:
  - private key auth
  - passphrase support
  - ssh-agent support
- Remote command execution framework
- Deploy orchestrator replicating every v1.x feature  
- Now powered by **Cobra** (v2.0.5)

The goal is to evolve this into a complete internal platform tool.

---

# Features (Current State)

### **Deployment Engine**
- Install base packages
- Provision DBs (MySQL/Postgres/MongoDB)
- Clone repository with Git PAT
- Maven-based build pipeline
- systemd service generation
- env file generation + permissions
- Restart + lifecycle automation

### **Go CLI**
- `vmdeploy deploy`
- YAML config loading
- Fully automated SSH workflow
- Remote execution framework
- Preflight system checks

### **Security**
- No secrets stored locally  
- Supports SSH agent / PEM / passphrase  
- Root-only env file permissions

---

# Upcoming (v2.1+)

- Multi-server deployments  
- Parallel deployments  
- Rollback mechanism  
- DB schema migrations  
- Structured logging  
- Built-in audit trail  
- Test suite + CI/CD  
- Installable Homebrew + APT package  
- Plugin architecture  

---

# Philosophy

This is **not** a production orchestrator (yet).
This is my attempt to:

- learn infra engineering deeply  
- think like a platform architect  
- design like CNCF projects  
- iterate with real-world discipline  

Every version teaches one real concept:
deployment lifecycle, SSH, DB provisioning, packaging, or CLI design.

---

If you're curious about how infra tools actually work,  
this project might feel like home.
