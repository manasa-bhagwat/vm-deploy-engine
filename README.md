# VM Deploy Engine

A lightweight deployment engine that provisions Linux VMs, installs runtimes, configures databases, builds applications, and creates systemd-managed services — all from a predictable CLI workflow.

Designed as both a practical tool and a study of how real infrastructure automation systems work under the hood.

---

# Purpose

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

# Architecture Evolution

## v1.x — Shell-Based Deployment Engine
A complete lifecycle automation script that performs:
- Install Java, Git, MySQL
- Create DB + user
- Clone repo using PAT
- Build Spring Boot app with Maven Wrapper
- Normalize JAR
- Generate env file
- Generate systemd service  
- Launch the service predictably
- Central configuration file (`/etc/app-deploy.conf`)
- Supports **MySQL**, **PostgreSQL**, **MongoDB**
- Hidden prompts for PAT + DB credentials
- DB readiness checks
- Root-owned env file with 600 permissions  
- Strict Bash safety — `set -euo pipefail`
  
This version demonstrates how traditional automation is built before infra-as-code existed.

---

## v2.0.x — Go CLI Rewrite (Platform-Oriented) 
The entire engine re-implemented using modern Go engineering patterns:
- Structured Go module layout (`cmd/`, `internal/`, `pkg/`)
- YAML configs (`appconfig.yaml`, `vmconfig.yaml`)
- SSH automation via:
  - private key auth
  - passphrase support
  - ssh-agent support
- Remote command execution framework
- Deploy orchestrator replicating every v1.x feature  
- Now powered by **Cobra** (v2.0.5)

This version is the foundation for evolving VM Deploy Engine into a full internal deployment tool.

---

# Current Feature Set

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

# Roadmap (v2.1 and beyond)

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

Every version teaches one real concept.

If you're curious about how infra tools actually work, this project might feel like home.
