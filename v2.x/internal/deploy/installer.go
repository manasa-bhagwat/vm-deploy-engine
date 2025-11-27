package deploy

import (
	"fmt"

	"github.com/manasa-bhagwat/vm-deploy-engine/v2/pkg/sshutil"
)

type Installer struct {
	Client *sshutil.SSHClient
}

// InstallBasePackages() installs Java, Git, Curl

func (i *Installer) InstallBasePackages() error {
	cmd := `sudo apt update -y && sudo apt install -y openjdk-17-jdk git curl`
	_, stderr, err := i.Client.Run(cmd)
	if err != nil {
		return fmt.Errorf("base package install failed: %v (stderr=%s)", err, stderr)
	}
	return nil
}
