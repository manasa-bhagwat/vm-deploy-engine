package deploy

import (
	"fmt"
	"strings"

	"github.com/manasa-bhagwat/vm-deploy-engine/v2/pkg/sshutil"
)

type GitOps struct {
	Client *sshutil.SSHClient
}

// CloneAsUser clones the repo using sudo -u <user>
// so that Maven has write permissions to /opt/app
func (g *GitOps) CloneAsUser(repoURL, branch, targetDir, pat, user string) error {

	// strip "https://" for PAT injection
	repoWithoutHTTPS := strings.TrimPrefix(repoURL, "https://")

	// final clone command
	cmd := fmt.Sprintf(
		"sudo -u %s git clone -b %s https://%s@%s %s",
		user, branch, pat, repoWithoutHTTPS, targetDir,
	)

	stdout, stderr, err := g.Client.Run(cmd)
	if err != nil {
		return fmt.Errorf("git clone failed: %v (stderr: %s)", err, stderr)
	}

	if stdout != "" {
		fmt.Println("[GIT STDOUT]:", stdout)
	}
	if stderr != "" {
		fmt.Println("[GIT STDERR]:", stderr)
	}

	return nil
}
