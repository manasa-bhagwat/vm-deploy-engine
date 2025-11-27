#!/bin/bash
set -euo pipefail

### ============================================================
### UNIVERSAL SPRING BOOT DEPLOY SCRIPT (SECURE v2.7)
### ALWAYS PROMPTS:
###   - GitHub PAT (hidden)
###   - DB Username
###   - DB Password (hidden)
### No secrets stored anywhere.
### Works for MySQL / Postgres / MongoDB
### ============================================================

CONFIG_FILE="/etc/app-deploy.conf"
REPO_DIR="/opt/app"

### ============================================================
### 0. Load Config File FIRST
### ============================================================
if [[ ! -f "$CONFIG_FILE" ]]; then
  echo "[ERROR] Config file not found: $CONFIG_FILE"
  exit 1
fi

source "$CONFIG_FILE"

echo "[INFO] Starting deployment for app: $APP_NAME"


### ============================================================
### 1. ALWAYS ASK FOR GITHUB TOKEN (NEVER STORED)
### ============================================================
read -sp "Enter GitHub Personal Access Token (hidden): " GITHUB_PAT
echo


### ============================================================
### 2. ALWAYS PROMPT FOR DB USER & PASSWORD
### ============================================================
read -p "Enter database username: " DB_USER
read -sp "Enter database password (hidden): " DB_PASS
echo


### ============================================================
### 3. Install Dependencies
### ============================================================
echo "[INFO] Installing dependencies..."
sudo apt update -y
sudo apt install -y openjdk-17-jdk git curl

case "$DB_TYPE" in
  mysql)
    sudo apt install -y mysql-server mysql-client
    ;;
  postgres)
    sudo apt install -y postgresql postgresql-client
    ;;
  mongodb)
    sudo apt install -y mongodb
    ;;
  *)
    echo "[ERROR] Unsupported DB_TYPE: $DB_TYPE"
    exit 1
    ;;
esac


### ============================================================
### 4. Start Database & Create DB/User
### ============================================================
echo "[INFO] Configuring database: $DB_TYPE"

if [[ "$DB_TYPE" == "mysql" ]]; then
  sudo systemctl enable mysql
  sudo systemctl start mysql

  sudo mysql <<EOF
CREATE DATABASE IF NOT EXISTS $DB_NAME;
CREATE USER IF NOT EXISTS '$DB_USER'@'localhost' IDENTIFIED BY '$DB_PASS';
GRANT ALL PRIVILEGES ON $DB_NAME.* TO '$DB_USER'@'localhost';
FLUSH PRIVILEGES;
EOF

elif [[ "$DB_TYPE" == "postgres" ]]; then
  sudo systemctl enable postgresql
  sudo systemctl start postgresql

  sudo -u postgres psql <<EOF
CREATE DATABASE $DB_NAME;
CREATE USER $DB_USER WITH ENCRYPTED PASSWORD '$DB_PASS';
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;
EOF

elif [[ "$DB_TYPE" == "mongodb" ]]; then
  sudo systemctl enable mongodb
  sudo systemctl start mongodb

  mongo <<EOF
use $DB_NAME
db.createUser({
  user: "$DB_USER",
  pwd: "$DB_PASS",
  roles: ["readWrite"]
})
EOF
fi


### ============================================================
### 5. Wait for DB Readiness
### ============================================================
echo "[INFO] Waiting for DB to become ready..."

for i in {1..10}; do
  if [[ "$DB_TYPE" == "mysql" ]]; then
    mysql -u"$DB_USER" -p"$DB_PASS" -e "SELECT 1" >/dev/null 2>&1 && break
  elif [[ "$DB_TYPE" == "postgres" ]]; then
    PGPASSWORD="$DB_PASS" psql -h localhost -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1" >/dev/null 2>&1 && break
  elif [[ "$DB_TYPE" == "mongodb" ]]; then
    mongo --eval "db.runCommand({ ping: 1 })" >/dev/null 2>&1 && break
  fi

  echo " - DB not ready yet… retrying"
  sleep 2
done

echo "[INFO] Database is ready."


### ============================================================
### 6. Clone Repository
### ============================================================
echo "[INFO] Cloning repository..."
sudo rm -rf "$REPO_DIR"
sudo git clone -b "$BRANCH" https://$GITHUB_PAT@${REPO_URL#https://} "$REPO_DIR"


### ============================================================
### 7. Build Spring Boot App
### ============================================================
echo "[INFO] Building application..."
sudo chmod +x $REPO_DIR/mvnw

cd "$REPO_DIR"
./mvnw clean package -DskipTests

sudo cp target/*.jar "$REPO_DIR/$ARTIFACT_NAME"


### ============================================================
### 8. Generate Environment File
### ============================================================
ENV_FILE="/etc/$SERVICE_NAME.env"

echo "[INFO] Writing environment variables..."

if [[ "$DB_TYPE" == "mysql" ]]; then
  DB_URL="jdbc:mysql://$DB_HOST:$DB_PORT/$DB_NAME"
elif [[ "$DB_TYPE" == "postgres" ]]; then
  DB_URL="jdbc:postgresql://$DB_HOST:$DB_PORT/$DB_NAME"
elif [[ "$DB_TYPE" == "mongodb" ]]; then
  DB_URL="mongodb://$DB_USER:$DB_PASS@$DB_HOST:$DB_PORT/$DB_NAME"
fi

sudo bash -c "cat > $ENV_FILE" <<EOF
SPRING_DATASOURCE_URL=$DB_URL
SPRING_DATASOURCE_USERNAME=$DB_USER
SPRING_DATASOURCE_PASSWORD=$DB_PASS
SERVER_PORT=$APP_PORT
SERVER_ADDRESS=0.0.0.0
EOF

sudo chmod 600 $ENV_FILE
sudo chown root:root $ENV_FILE


### ============================================================
### 9. Create & Start systemd Service
### ============================================================
SERVICE_PATH="/etc/systemd/system/$SERVICE_NAME.service"

echo "[INFO] Creating systemd service..."

sudo bash -c "cat > $SERVICE_PATH" <<EOF
[Unit]
Description=Spring Boot App: $APP_NAME
After=network.target

[Service]
User=root
EnvironmentFile=$ENV_FILE
ExecStart=/usr/bin/java -jar $REPO_DIR/$ARTIFACT_NAME
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable "$SERVICE_NAME"
sudo systemctl restart "$SERVICE_NAME"

echo "[SUCCESS] Application deployed on port $APP_PORT."