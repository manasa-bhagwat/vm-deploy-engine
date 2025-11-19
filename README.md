
# 🚀 VM Deploy Engine

*A lightweight, script-powered deployment engine for provisioning and deploying apps on Linux VMs.*

This project is my attempt to understand **how real infrastructure tools work under the hood** — service lifecycle, builds, config management, DB provisioning, systemd automation, and secure scripting.

It currently includes:

* **v1.0** – baseline deployment script
* **v1.1** – config-driven, secure multi-DB deploy engine
* **v2.0 (upcoming)** – Go CLI rewrite

📦 Repo: [https://github.com/manasa-bhagwat/vm-deploy-engine](https://github.com/manasa-bhagwat/vm-deploy-engine)

---

# 📜 Table of Contents

* [Project Overview](#project-overview)
* [How It Works](#how-it-works)
* [Architecture](#architecture)
* [Features](#features)
* [Versioning](#versioning)
* [Common Errors & Fixes](#common-errors--fixes)
* [Contributing](#contributing)
* [Roadmap](#roadmap)
* [License](#license)

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

## ⚙️ How It Works

At a high level, the engine performs:

1. **VM Preparation**

   * Install Java, Git, DB runtimes
   * Enable/start DB services

2. **Database Provisioning**

   * MySQL, PostgreSQL, MongoDB support
   * DB/user creation
   * Idempotent provisioning
   * Readiness checks

3. **App Acquisition**

   * Clone GitHub repo via PAT
   * Use Maven Wrapper for consistent builds

4. **Build & Artifact Handling**

   * Clean + package
   * Normalize JAR name
   * Copy into predictable path

5. **Configuration**

   * Generate `/etc/<service>.env`
   * Store DB + server settings
   * Apply strict permissions

6. **Service Management**

   * Create systemd unit
   * Register service
   * Ensure autostart
   * Restart on failure

7. **Observability**

   * Viewable logs via `journalctl`
   * Service state via `systemctl status`

This is the same lifecycle followed by production tools—just simplified and explicit.

---

# 🏗 Architecture

```
vm-deploy-engine/
│
├── scripts/
│   ├── v1.0/              # baseline deploy script
│   └── v1.1/              # config-driven deploy engine
│
├── docs/
│   ├── v1.0-overview.md
│   └── v1.1-overview.md
│
├── roadmap.md
└── README.md
```


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

---

## 🐛 Common Errors & Fixes


### ❗ *mvnw: Permission denied*

Fix:

```bash
chmod +x mvnw
```

---

### ❗ *MySQL not ready: Communications link failure*

Fix:

* Use readiness loops (v1.0 & v1.1 handle this)
* Increase sleep time if needed

---

### ❗ *JAR not found in target/***.jar*

Usually means build failed.

Fix:

```bash
./mvnw clean package
```

---

### ❗ *Service keeps exiting*

Check logs:

```bash
journalctl -u <service> -f
```

---

## 🤝 Contributors Guide

 🔧 **Development Workflow**

1. Fork the repository
2. Create a feature branch:

```bash
git checkout -b feature/<name>
```

3. Follow versioned structure:

```
scripts/vX.Y/
docs/vX.Y-overview.md
```

4. Write clear commit messages:

```
feat: add new db provisioning flow  
fix: correct env permissions  
docs: update v1.1 overview  
```

5. Submit PR → Wait for review

---

### 📐 Code Style

* Bash → `set -euo pipefail` mandatory
* No secrets in scripts
* Config split from logic
* Use functions for repeatable logic
* Follow v1.0 → v1.1 → v2.0 branch patterns

---

## Roadmap

### v2.0 — Go CLI

* Deploy
* Logs
* Status
* Rollback
* YAML config
* DB provisioning
* Concurrency
* SSH agent mode
* Multi-runtime Support
  * Java
  * Python
  * Node
  * .NET

### v3.0 — Mini Vault

* Secrets storage
* Encryption
* Rotation

### v4.0 — Mini CI/CD

* Pipelines
* Steps
* Deploy triggers

---

## 📄 License

MIT License — free for anyone to study, fork, and build upon.

---

