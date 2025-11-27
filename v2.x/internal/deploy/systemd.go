package deploy

import (
	"fmt"

	"github.com/manasa-bhagwat/vm-deploy-engine/v2/pkg/sshutil"
)

type SystemdOps struct {
	Client *sshutil.SSHClient
}

func (s *SystemdOps) WriteEnvFile(path, dbURL, dbUser, dbPass string, port int) error {
	cmd := fmt.Sprintf(`
sudo bash -c 'cat > %[1]s <<EOF
SPRING_DATASOURCE_URL=%[2]s
SPRING_DATASOURCE_USERNAME=%[3]s
SPRING_DATASOURCE_PASSWORD=%[4]s
SERVER_PORT=%[5]d
SERVER_ADDRESS=0.0.0.0
EOF'
`, path, dbURL, dbUser, dbPass, port)

	_, stderr, err := s.Client.Run(cmd)
	if err != nil {
		return fmt.Errorf("env file write failed: %v (stderr=%s)", err, stderr)
	}
	return nil
}

func (s *SystemdOps) WriteService(path, name, repoDir, artifact string) error {
	cmd := fmt.Sprintf(`
sudo bash -c 'cat > %[1]s <<EOF
[Unit]
Description=Spring Boot App: %[2]s
After=network.target

[Service]
User=root
EnvironmentFile=/etc/%[2]s.env
ExecStart=/usr/bin/java -jar %[3]s/%[4]s
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF'
`, path, name, repoDir, artifact)

	_, stderr, err := s.Client.Run(cmd)
	if err != nil {
		return fmt.Errorf("service file write failed: %v (stderr=%s)", err, stderr)
	}
	return nil
}

func (s *SystemdOps) ReloadAndRestart(name string) error {
	cmd := fmt.Sprintf(`
sudo systemctl daemon-reload &&
sudo systemctl enable %[1]s &&
sudo systemctl restart %[1]s
`, name)

	_, stderr, err := s.Client.Run(cmd)
	if err != nil {
		return fmt.Errorf("systemd restart failed: %v (stderr=%s)", err, stderr)
	}
	return nil
}
