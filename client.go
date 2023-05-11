package parallexe

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
)

// createClient creates a new SSH client
// If sshConfig.PrivateKeyPath and sshConfig.Password are empty, it will try to connect to local SSH agent
func createClient(addr string, sshConfig *SshConfig) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            sshConfig.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	authMethods, err := getAuthMethod(sshConfig.Password, sshConfig.PrivateKeyPath, sshConfig.PrivateKey)
	if err != nil {
		return nil, err
	}

	config.Auth = authMethods

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", addr), config)
	if err != nil {
		return nil, fmt.Errorf("can't connect to %s: %v", addr, err)
	}

	return conn, nil
}

// getAuthMethod returns the SSH auth method to use
// If password is not empty, it will use it
// If privateKey is not empty, it will use it
// If privateKeyPath is not empty, it will use it
// Otherwise, it will try to connect to local SSH agent
func getAuthMethod(password string, privateKeyPath string, privateKey []byte) ([]ssh.AuthMethod, error) {
	if len(password) > 0 {
		return []ssh.AuthMethod{
			ssh.Password(password),
		}, nil
	}

	if len(privateKey) > 0 {
		signer, err := ssh.ParsePrivateKey(privateKey)
		if err != nil {
			return []ssh.AuthMethod{}, fmt.Errorf("can't parse private key: %s", err)
		}

		return []ssh.AuthMethod{ssh.PublicKeys(signer)}, nil
	}

	if len(privateKeyPath) > 0 {
		key, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return []ssh.AuthMethod{}, fmt.Errorf("can't read private key: %s", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return []ssh.AuthMethod{}, fmt.Errorf("can't parse private key: %s", err)
		}

		return []ssh.AuthMethod{ssh.PublicKeys(signer)}, nil
	}

	// Create connection with local SSH agent
	conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return []ssh.AuthMethod{}, fmt.Errorf("can't connect to local SSH agent: %s", err)
	}

	// Create client agent
	agentClient := agent.NewClient(conn)

	return []ssh.AuthMethod{ssh.PublicKeysCallback(agentClient.Signers)}, nil
}
