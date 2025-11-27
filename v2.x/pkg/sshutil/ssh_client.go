package sshutil

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSHConfig holds all the information needed to connect to a remote VM via SSH.
type SSHConfig struct {
	Host string
	Port int
	User string

	SSHKeyPath    string
	SSHPassphrase string
	UseSSHAgent   bool
}

// SSHClient is a small wrapper around *ssh.Client
type SSHClient struct {
	cfg    SSHConfig
	client *ssh.Client
}

// NewSSHClient constructs a new SSHClient.
func NewSSHClient(cfg SSHConfig) *SSHClient {
	return &SSHClient{cfg: cfg}
}

// Ping tries to open a TCP connection to host:port (no SSH handshake yet).
func (c *SSHClient) Ping() error {
	addr := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return fmt.Errorf("tcp ping to %s failed: %w", addr, err)
	}
	defer conn.Close()
	return nil
}

// buildAuthMethods prepares SSH authentication methods: ssh-agent and/or key file.
func (c *SSHClient) buildAuthMethods() ([]ssh.AuthMethod, error) {
	var methods []ssh.AuthMethod

	// 1) ssh-agent (if enabled)
	if c.cfg.UseSSHAgent {
		if sock := os.Getenv("SSH_AUTH_SOCK"); sock != "" {
			conn, err := net.Dial("unix", sock)
			if err == nil {
				ag := agent.NewClient(conn)
				methods = append(methods, ssh.PublicKeysCallback(ag.Signers))
			}
		}
	}

	// 2) Key file (if provided)
	if c.cfg.SSHKeyPath != "" {
		keyPath := os.ExpandEnv(c.cfg.SSHKeyPath)
		keyPath = filepath.Clean(keyPath)

		key, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("cannot read ssh key %s: %w", keyPath, err)
		}

		// Try without passphrase
		signer, err := ssh.ParsePrivateKey(key)
		if err == nil {
			methods = append(methods, ssh.PublicKeys(signer))
			return methods, nil
		}

		// Try with passphrase if configured
		if c.cfg.SSHPassphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(c.cfg.SSHPassphrase))
			if err == nil {
				methods = append(methods, ssh.PublicKeys(signer))
				return methods, nil
			}
			return nil, fmt.Errorf("failed to parse private key with passphrase: %w", err)
		}

		return nil, fmt.Errorf("failed to parse private key (possibly encrypted): %w", err)
	}

	if len(methods) == 0 {
		return nil, fmt.Errorf("no SSH authentication methods available")
	}

	return methods, nil
}

// Connect establishes an SSH connection.
func (c *SSHClient) Connect() error {
	authMethods, err := c.buildAuthMethods()
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)

	sshCfg := &ssh.ClientConfig{
		User:            c.cfg.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	client, err := ssh.Dial("tcp", addr, sshCfg)
	if err != nil {
		return fmt.Errorf("ssh connect failed to %s: %w", addr, err)
	}

	c.client = client
	return nil
}

// Close closes the underlying SSH connection.
func (c *SSHClient) Close() error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}

// Run executes a single command on the remote VM and returns stdout, stderr, error.
func (c *SSHClient) Run(command string) (string, string, error) {
	if c.client == nil {
		return "", "", fmt.Errorf("ssh client is not connected")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return "", "", fmt.Errorf("failed to create ssh session: %w", err)
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	err = session.Run(command)

	stdout := stdoutBuf.String()
	stderr := stderrBuf.String()

	if err != nil {
		return stdout, stderr, fmt.Errorf("remote command failed: %w (stderr: %s)", err, stderr)
	}

	return stdout, stderr, nil
}
