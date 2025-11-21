package sshutil

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHConfig holds all the information needed to connect to a remote vm over SSH
type SSHConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

type SSHClient struct {
	cfg    SSHConfig
	client *ssh.Client
}

// NewSSHClient is a simple constructor that returns a pointer to SSHClient
// It does NOT connect yet. It only stores the settings
func NewSSHClient(cfg SSHConfig) *SSHClient {
	return &SSHClient{
		cfg: cfg,
	}
}

// Connect dails the remote VM and creates an SSH connection
// It must be called before Run()
func (c *SSHClient) Connect() error {

	// 1. Build the address string like "1.2.3.4:22"
	addr := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)

	// 2. Prepare the SSH authentication method
	// for now we only support auth: ssh.Password()
	auth := ssh.Password(c.cfg.Password)

	// 3. Create the SSH client config
	SSHConfig := &ssh.ClientConfig{
		User:            c.cfg.User,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// 4. Use ssh.Dial to open a TCP connection + handshake SSH.
	client, err := ssh.Dial("tcp", addr, SSHConfig)
	if err != nil {
		return fmt.Errorf("ssh connect failed to %s: %w", addr, err)
	}

	// 5. Save the connected client on our struct so Run() can use it.
	c.client = client
	return nil
}

// Close() closes the underlying SSH connection.
// We should always close when we're done to avoid leaking resources
func (c *SSHClient) Close() error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}

// Run executes a single command on the remote VM and returns:
// - stdout (normal output)
// - stderr (error output)
// - error (if something went wrong with SSH or execution)
func (c *SSHClient) Run(command string) (string, string, error) {
	if c.client == nil {
		return "", "", fmt.Errorf("ssh client is not connected!")
	}

	// 1. Create a new SSH session
	session, err := c.client.NewSession()
	if err != nil {
		return "", "", fmt.Errorf("failed to create ssh session: %w", err)
	}
	// 2. Always close session when this function finishes.
	defer session.Close()

	// 3. Prepare buffers to capture stdout and stderr.
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	// 4. Tell the session to write outputs into those buffers.
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	// 5. Run the command on the remote machine.
	// Note: even if the command returns non-zero exit code,
	//       err will be non-nil. So we still want stdout/stderr.
	if err := session.Run(command); err != nil {
		return stdoutBuf.String(), stderrBuf.String(),
			fmt.Errorf("remote command failed: %w", err)
	}

	stdout := stdoutBuf.String()
	stderr := stderrBuf.String()

	if err != nil {
		// Wrap error with stderr to help debugging later.
		return stdout, stderr, fmt.Errorf("remote command failed: %w (stderr: %s)", err, stderr)
	}

	// If everything is OK, return the outputs and nil error.
	return stdout, stderr, nil
}

// Ping tries to open a TCP connection to the SSH host:port without doing full SSH handshake.
// This can be used later as a "pre-flight" connectivity check before Connect().
func (c *SSHClient) Ping() error {
	addr := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return fmt.Errorf("tcp ping to %s failed: %w", addr, err)
	}
	defer conn.Close()
	return nil
}
