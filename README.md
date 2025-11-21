
# 🚀 VM Deploy Engine

*A lightweight, script-powered deployment engine for provisioning and deploying apps on Linux VMs.*

This project is my attempt to understand **how real infrastructure tools work under the hood** — service lifecycle, builds, config management, DB provisioning, systemd automation, and secure scripting.

It currently includes:

* **v1.0** – baseline deployment script
* **v1.1** – config-driven, secure multi-DB deploy engine
* **v2.0 (upcoming)** – Go CLI rewrite

📦 Repo: [https://github.com/manasa-bhagwat/vm-deploy-engine](https://github.com/manasa-bhagwat/vm-deploy-engine)

---

# 🧭 Project Overview

The goal of **VM Deploy Engine** is simple:

> Recreate the fundamentals behind real-world infra tools — manually, intentionally, and with complete understanding.

Instead of relying on Ansible, Terraform, or Kubernetes from day one, this project builds the *mechanics* from scratch:

* Install runtimes
* Configure databases
* Build artifacts
* Generate environment files
* Create services
* Automate deployments
* Enforce secure patterns

It mirrors how early platform teams build internal deployment tools — step-by-step, scaling complexity over time.

---

# ⭐ Features

### ✔ v1.0 — Baseline Deployment Script

* Java + MySQL installation
* DB creation
* Git clone using PAT
* Build using Maven Wrapper
* Normalize JAR
* systemd service
* env file generation

### ✔ v1.1 — Secure, Config-Driven Multi-DB Engine

* Central config file
* Secure prompts (PAT + DB password hidden)
* MySQL / PostgreSQL / MongoDB support
* DB readiness detection
* Strict bash safety (`set -euo pipefail`)
* Cleaner structure
* Root-owned env file with 600 perms

### v2.0 — (Upcoming) Go CLI

* Deploy / Status / Logs / Rollback commands
* YAML config
* SSH agent mode
* Parallel deployments
* Structured logs
* Built-in audit logs
* Cross-platform binaries


## 📄 License

MIT License — free for anyone to study, fork, and build upon.

---

