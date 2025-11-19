#!/bin/bash
set -e

### CONFIG ###
REPO_URL="<insert>"
REPO_DIR="/opt/java-app"
SERVICE_NAME="java-app"
BRANCH="main"
JAR_NAME="<insert>"
GITHUB_PAT="<insert>"
APP_PORT=8080

DB_HOST="localhost"
DB_PORT=3306
DB_USER="<insert>"
DB_PASS="<insert>"
DB_NAME="<insert>"
################

echo "[INFO] Installing Java 17 + MySQL"
sudo apt update -y
sudo apt install -y openjdk-17-jdk mysql-server mysql-client git

echo "[INFO] Starting MySQL service"
sudo systemctl enable mysql
sudo systemctl start mysql

echo "[INFO] Creating database and user (if not exists)"
sudo mysql <<EOF
CREATE DATABASE IF NOT EXISTS $DB_NAME;
CREATE USER IF NOT EXISTS '$DB_USER'@'localhost' IDENTIFIED BY '$DB_PASS';
GRANT ALL PRIVILEGES ON $DB_NAME.* TO '$DB_USER'@'localhost';
FLUSH PRIVILEGES;
EOF

echo "[INFO] Waiting for MySQL to be ready..."
until mysql -u$DB_USER -p$DB_PASS -h $DB_HOST -P $DB_PORT -e "SELECT 1" >/dev/null 2>&1; do
  echo " - MySQL not ready yet, retrying..."
  sleep 2
done
echo "[INFO] MySQL is ready."

echo "[INFO] Cloning repo"
sudo rm -rf $REPO_DIR
sudo git clone -b $BRANCH https://$GITHUB_PAT@${REPO_URL#https://} $REPO_DIR

echo "[INFO] Fixing mvnw permissions"
sudo chmod +x $REPO_DIR/mvnw

echo "[INFO] Building app"
cd $REPO_DIR
./mvnw clean package -DskipTests

echo "[INFO] Copying JAR"
sudo cp target/*.jar $REPO_DIR/$JAR_NAME

echo "[INFO] Writing environment file"
sudo bash -c "cat > /etc/$SERVICE_NAME.env" <<EOF
SPRING_DATASOURCE_URL=jdbc:mysql://$DB_HOST:$DB_PORT/$DB_NAME
SPRING_DATASOURCE_USERNAME=$DB_USER
SPRING_DATASOURCE_PASSWORD=$DB_PASS
SERVER_PORT=$APP_PORT
SERVER_ADDRESS=0.0.0.0
EOF

echo "[INFO] Creating systemd service"
sudo bash -c "cat > /etc/systemd/system/$SERVICE_NAME.service" <<EOF
[Unit]
Description=Spring Boot App
After=syslog.target mysql.service

[Service]
User=root
EnvironmentFile=/etc/$SERVICE_NAME.env
ExecStart=/usr/bin/java -jar $REPO_DIR/$JAR_NAME
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

echo "[INFO] Starting service"
sudo systemctl daemon-reload
sudo systemctl enable $SERVICE_NAME
sudo systemctl restart $SERVICE_NAME

echo "[DONE] Spring Boot App deployed on port $APP_PORT"