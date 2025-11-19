# vm-deploy-engine

A universal VM application deployment engine that automates provisioning, building, and running applications across multiple runtimes and databases — starting from Bash (v1.x) and evolving into a full Go-based CLI (v2.x+).

This project is my long-term journey to understand how real infrastructure tools are built under the hood.  
It begins with simple Bash automation and grows into CNCF-style tooling with Go, modular design, concurrency, SSH agents, secret management, and multi-runtime support.

---

## 🚀 Vision

`vm-deploy-engine` aims to become a lightweight, developer-friendly deployment orchestrator capable of:

- Deploying any language runtime (Java, Node, Python, Go, .NET, etc.)
- Supporting multiple databases (MySQL, Postgres, MongoDB)
- Managing VM lifecycles (build → deploy → start → logs → rollback)
- Automating systemd service creation
- Providing a clean CLI interface
- Running over SSH across multiple servers
- Evolving toward GitOps-style reconciliation patterns

It starts simple.  
It grows powerful.  
It stays transparent.

---

# 📦 Version Overview

### **v1.0 — Baseline Deployment Script**
A single Bash script that:
- Installs Java, MySQL, Git
- Creates DB + user (idempotent)
- Waits for DB readiness
- Clones a GitHub repo
- Builds Spring Boot app (mvnw)
- Creates systemd service
- Starts the app reliably on reboot

👉 Focus: **Learn the full VM deployment lifecycle**

---

### **v1.1 — Config-Driven Multi-DB Deployment**
Enhancements:
- Config file (`/etc/app-deploy.conf`)
- Secure prompts (GitHub PAT, DB password)
- Multi-DB support (MySQL/Postgres/MongoDB)
- Better readiness checks
- Cleaner folder structure & permissions

👉 Focus: **Declarative config + security + flexibility**

---

### **v2.0 — Go CLI (In Progress)**
This is the next major milestone where the project becomes a proper tool.

Planned features:
- Go-based CLI (`vmde` or `vmctl`)
- Concurrency (goroutines)
- Runtime detection
- Multi-language builders
- Structured logs (JSON/text)
- DB drivers for provisioning
- SSH deploy mode
- Versioned releases
- Rollback directories
- Audit logs

---

# 🔥 Why v2.0 Moves to Go

Bash is perfect for learning how deployments work on Linux.  
But the next requirements need real engineering:

- concurrency  
- SSH to multiple servers  
- runtime detection  
- secrets  
- structured logs  
- clear error handling  
- testability  
- interfaces  
- modular design  

Go is the language used for almost all CNCF projects (Kubernetes, Helm, ArgoCD, containerd, etc.)  
The v2.x versions bring the project closer to that ecosystem.

---

# 🧪 Testing Strategy

### v1.x  
- No unit tests (Bash)
- Functional tests done manually on VMs

### v2.x  
- Go tests (`go test`)
- Module-level tests for:
  - config parsing  
  - DB provisioning  
  - runtime detection  
  - service generation  
  - SSH executor  

---

# 📜 License  
MIT License

---

# 🙌 Contributions  
This is a personal learning + engineering exploration project, but contributions are welcome.

---

# 👩‍💻 Author  
**Manasa** — Software Engineer, Infra Tools Enthusiast, CNCF Learner