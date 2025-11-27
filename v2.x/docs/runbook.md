# 🧭 vmdeploy RUNBOOK
Operational guide for deploying, debugging, and operating **vmdeploy v2.0.5**.  
This runbook is written for Infra/Platform engineers to execute, diagnose, and recover vmdeploy deployments safely.

---

# 📌 1. Overview
`vmdeploy` is a Go-based CLI tool that deploys applications to Linux VMs over SSH.

It performs:

1. SSH connectivity checks  
2. Installation of base packages  
3. Git clone of the application  
4. Remote Maven build  
5. Database configuration  
6. systemd service creation  
7. Application startup & verification  

This runbook documents **how to operate**, **troubleshoot**, and **recover** vmdeploy.

---

# 📁 2. Prerequisites

### ✔️ Local machine
- Go 1.22+
- SSH private key (`.pem` / `.rsa`)
- SSH agent optional (`ssh-add <key>`)

### ✔️ Remote VM
- Ubuntu/Debian recommended
- OpenSSH enabled
- Ports open:
```

22/tcp
8080 (or app port)

```
- User has sudo privileges

### ✔️ Required files
```

appconfig.yaml
vmconfig.yaml

````

### ✔️ Required secrets
- DB Username
- DB Password
- GitHub PAT with repo access

---

# 📦 3. Installation

### Option A — From source
```bash
go build -o vmdeploy ./cmd/vmdeploy
sudo mv vmdeploy /usr/local/bin/
````

### Option B — From GitHub Release

Download the asset:

```
vmdeploy-linux-amd64
vmdeploy-macos-arm64
vmdeploy-windows-amd64.exe
```

Make executable & move to PATH:

```bash
chmod +x vmdeploy
sudo mv vmdeploy /usr/local/bin/
```

Verify installation:

```bash
vmdeploy version
```

---

# ⚙️ 4. Configuration Files

## appconfig.yaml

Defines WHAT to deploy.

```yaml
app_name: bankapp
repo_url: "https://github.com/user/repo.git"
branch: "main"
artifact_name: "app.jar"

db_type: mysql
db_host: localhost
db_port: 3306
db_name: bankdb

service_name: bankapp
app_port: 8080
```

## vmconfig.yaml

Defines WHERE to deploy.

```yaml
host: 13.222.xx.xx
port: 22
user: ubuntu
ssh_key_path: "~/.ssh/mykey.pem"
ssh_passphrase: ""
use_ssh_agent: true
```

---

# 🚀 5. Deploying an Application

Run:

```bash
vmdeploy deploy \
  --app-config appconfig.yaml \
  --vm-config vmconfig.yaml
```

You will be prompted for:

```
DB Username:
DB Password:
GitHub PAT:
```

If successful, you will see:

```
[SUCCESS] Deployment completed successfully.
```

Validate service:

```bash
curl http://<vm-ip>:8080/actuator/health
```

---

# 🧔 6. Operational Procedures

## Restart the app

```bash
ssh ubuntu@<vm> "sudo systemctl restart bankapp"
```

## View logs

```bash
ssh ubuntu@<vm> "sudo journalctl -u bankapp -f"
```

## Check service status

```bash
ssh ubuntu@<vm> "systemctl status bankapp"
```

## View deployed JAR

```bash
ssh ubuntu@<vm> "ls -lah /opt/app/"
```

---

# 🐛 7. Troubleshooting Guide

## 7.1 SSH Failures

### ❌ Error:

```
ssh: handshake failed: no supported methods remain
```

### ✅ Fix:

* Ensure `PasswordAuthentication` is NO.
* Ensure key exists & permissions correct:

```bash
chmod 400 ~/.ssh/key.pem
```

* If using SSH Agent:

```bash
ssh-add ~/.ssh/key.pem
```

---

## 7.2 Permission Issues During Maven Build

### ❌ Error:

```
Cannot create resource output directory: /opt/app/target/classes
```

### ✅ Fix:

Run on VM:

```bash
sudo mkdir -p /opt/app
sudo chown -R ubuntu:ubuntu /opt/app
```

---

## 7.3 Git Clone Errors

### ❌ Error:

```
authentication failed for GitHub
```

### Checklist:

* PAT correct?
* PAT has **repo** scope?
* Repo private or public?
* Correct URL in appconfig?

---

## 7.4 systemd Service Failing

### Diagnose:

```bash
ssh ubuntu@vm "systemctl status bankapp"
```

### Logs:

```bash
ssh ubuntu@vm "sudo journalctl -u bankapp -n 50"
```

Common issues:

* Wrong JAR path
* Wrong DB credentials
* Wrong port binding
* Java not installed

---

## 7.5 Deployment Pipeline Stops Early

Check SSH error messages in logs:

Example:

```
[ERROR] Deployment failed: maven build failed
```

SSH into VM:

```bash
cd /opt/app
./mvnw clean package -DskipTests -X
```

---

# 🔁 8. Recovery Procedures

## Roll back to previous JAR

```bash
ssh ubuntu@vm "sudo cp /opt/app/backup/app.jar /opt/app/app.jar"
sudo systemctl restart bankapp
```

## Re-run full deployment

```bash
vmdeploy deploy
```

## Reset deployment directory

```bash
ssh ubuntu@vm "sudo rm -rf /opt/app && sudo mkdir -p /opt/app"
```

---

# 📘 9. Log Collection (for debugging / filing issues)

Run:

```bash
ssh ubuntu@vm "
  echo '=== systemd status ===';
  systemctl status bankapp;
  echo '=== journal logs ===';
  sudo journalctl -u bankapp -n 200;
  echo '=== app directory ===';
  ls -lah /opt/app;
"
```

---

# 🧪 10. Testing vmdeploy (Manual)

### ✔ Connectivity test

```bash
ssh ubuntu@<vm> uname -a
```

### ✔ Dry-run (simulated)

> Coming in v2.1.x

### ✔ Validate configs

Make sure YAML parses correctly:

```bash
yq e '.' appconfig.yaml
yq e '.' vmconfig.yaml
```

---

# 📥 11. Upgrading vmdeploy

Update binary and re-run:

```bash
vmdeploy version
vmdeploy deploy
```

### Breaking changes from Bash v1.x → Go v2.x

| Feature    | v1.x                 | v2.x                  |
| ---------- | -------------------- | --------------------- |
| Config     | Inline inside script | YAML based            |
| Auth       | Password/PAT only    | SSH key + agent       |
| Structure  | Bash scripts         | CNCF-style Go modules |
| Deployment | Single file          | Modular orchestrator  |

---

# 📌 12. Maintenance Tasks

### Clear old logs

```bash
journalctl --vacuum-time=7d
```

### Rotate logs

> Coming in v2.1.x

### Update base packages

```bash
ssh ubuntu@vm "sudo apt update -y && sudo apt upgrade -y"
```

---

# 🧱 13. Known Limitations (as of v2.0.5)

* No unit tests (coming in v2.1.x)
* No rollback mechanism
* No dry-run mode
* Single-VM only
* No structured JSON logging

---

# 🏁 14. Appendices

## A. Port Reference

* SSH → 22
* Application → from `appconfig.yaml`
* MySQL → 3306
* PostgreSQL → 5432
* MongoDB → 27017

## B. Default Paths

* App directory → `/opt/app`
* Env file → `/etc/<service>.env`
* systemd file → `/etc/systemd/system/<service>.service`

---

# 🎯 Final Notes

This runbook evolves with vmdeploy itself.
As vmdeploy becomes more cloud-native (v2.1.x+), this document will expand to cover:

* Observability
* Automated testing
* Dry-run mode
* Multi-VM deployments
* Chaos checks
* Rollbacks
* Secret managers (Vault, AWS Secrets Manager)

---

