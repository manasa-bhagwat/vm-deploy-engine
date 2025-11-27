package deploy

import (
	"fmt"

	"github.com/manasa-bhagwat/vm-deploy-engine/v2/internal/config"
	"github.com/manasa-bhagwat/vm-deploy-engine/v2/pkg/sshutil"
)

// DBInstaller is responsible for installing the DB server
// and provisioning database + user on the remote VM.
type DBInstaller struct {
	Client *sshutil.SSHClient
}

// InstallAndProvision installs the selected DB type (mysql/postgres/mongodb)
// and creates DB + user with the given credentials.
func (d *DBInstaller) InstallAndProvision(app *config.AppConfig, dbUser, dbPass string) error {
	switch app.DBType {
	case "mysql":
		return d.installMySQL(app, dbUser, dbPass)
	case "postgres":
		return d.installPostgres(app, dbUser, dbPass)
	case "mongodb":
		return d.installMongo(app, dbUser, dbPass)
	default:
		return fmt.Errorf("unsupported DB type: %s", app.DBType)
	}
}

func (d *DBInstaller) installMySQL(app *config.AppConfig, dbUser, dbPass string) error {
	cmd := fmt.Sprintf(`
sudo apt update -y &&
sudo apt install -y mysql-server mysql-client &&
sudo systemctl enable mysql &&
sudo systemctl start mysql &&
sudo mysql <<EOF
CREATE DATABASE IF NOT EXISTS %s;
CREATE USER IF NOT EXISTS '%s'@'localhost' IDENTIFIED BY '%s';
GRANT ALL PRIVILEGES ON %s.* TO '%s'@'localhost';
FLUSH PRIVILEGES;
EOF
`, app.DBName, dbUser, dbPass, app.DBName, dbUser)

	_, stderr, err := d.Client.Run(cmd)
	if err != nil {
		return fmt.Errorf("mysql setup failed: %v (stderr=%s)", err, stderr)
	}
	return nil
}

func (d *DBInstaller) installPostgres(app *config.AppConfig, dbUser, dbPass string) error {
	cmd := fmt.Sprintf(`
sudo apt update -y &&
sudo apt install -y postgresql postgresql-client &&
sudo systemctl enable postgresql &&
sudo systemctl start postgresql &&
sudo -u postgres psql <<EOF
DO
$$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = '%[1]s') THEN
      CREATE DATABASE %[1]s;
   END IF;
END
$$;
DO
$$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '%[2]s') THEN
      CREATE ROLE %[2]s LOGIN PASSWORD '%[3]s';
   END IF;
END
$$;
GRANT ALL PRIVILEGES ON DATABASE %[1]s TO %[2]s;
EOF
`, app.DBName, dbUser, dbPass)

	_, stderr, err := d.Client.Run(cmd)
	if err != nil {
		return fmt.Errorf("postgres setup failed: %v (stderr=%s)", err, stderr)
	}
	return nil
}

func (d *DBInstaller) installMongo(app *config.AppConfig, dbUser, dbPass string) error {
	cmd := fmt.Sprintf(`
sudo apt update -y &&
sudo apt install -y mongodb &&
sudo systemctl enable mongodb &&
sudo systemctl start mongodb &&
mongo <<EOF
use %s
db.createUser({
  user: "%s",
  pwd: "%s",
  roles: ["readWrite"]
})
EOF
`, app.DBName, dbUser, dbPass)

	_, stderr, err := d.Client.Run(cmd)
	if err != nil {
		return fmt.Errorf("mongodb setup failed: %v (stderr=%s)", err, stderr)
	}
	return nil
}
