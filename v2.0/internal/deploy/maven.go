package deploy

import (
	"fmt"

	"github.com/manasa-bhagwat/vm-deploy-engine/v2/pkg/sshutil"
)

type MavenOps struct {
	Client *sshutil.SSHClient
}

func (m *MavenOps) Build(repoDir, artifact string) error {
	cmd := fmt.Sprintf(`
		cd %[1]s &&
		sudo chmod +x mvnw &&
		./mvnw clean package -DskipTests &&
		sudo cp target/*.jar %[1]s/%[2]s
	`, repoDir, artifact)

	_, stderr, err := m.Client.Run(cmd)
	if err != nil {
		return fmt.Errorf("maven build failed: %v (stderr=%s)", err, stderr)
	}
	return nil
}
