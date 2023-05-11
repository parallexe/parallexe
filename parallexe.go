package parallexe

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"sync"
)

// localHostValues contains what is considered as localhost value given in HostConfig.Host
var localHostValues = []string{"localhost", "127.0.0.1"}

// SshConfig contains all needed configuration to SSH access to a specific host
type SshConfig struct {
	User           string
	Password       string
	PrivateKeyPath string
	PrivateKey     []byte
}

// HostConnection contains the SSH Client and the HostConfig.
// It is used internally to filter connections against what HostConfig contains
type HostConnection struct {
	HostConfig HostConfig
	Client     *ssh.Client
}

type HostConfig struct {
	SshConfig *SshConfig
	Host      string
	Groups    []string
}

type Parallexe struct {
	HostConnections []HostConnection
}

// New creates a new Parallexe client with a list of HostConfig
// It returns an error if at least one host is not reachable
// If Host is localhost or 127.0.0.1, it will create a HostConnection with Client nil
func New(configs []HostConfig) (*Parallexe, error) {
	var wg sync.WaitGroup
	wg.Add(len(configs))

	var p Parallexe

	hostErrors := make([]error, 0)

	for _, config := range configs {
		loopConfig := config
		go func() {
			defer wg.Done()

			err := p.AddHost(loopConfig)
			if err != nil {
				hostErrors = append(hostErrors, err)
			}
		}()
	}

	wg.Wait()

	if len(hostErrors) > 0 {
		return nil, fmt.Errorf("error while creating Parallexe client: %v", hostErrors)
	}

	return &p, nil
}

func (p *Parallexe) AddHost(hostConfig HostConfig) error {
	// Skip createClient if host is localhost
	for _, localHostValue := range localHostValues {
		if hostConfig.Host == localHostValue {
			p.HostConnections = append(p.HostConnections, HostConnection{
				HostConfig: hostConfig,
				Client:     nil,
			})

			return nil
		}
	}

	newClient, err := createClient(hostConfig.Host, hostConfig.SshConfig)
	if err != nil {
		return err
	}

	p.HostConnections = append(p.HostConnections, HostConnection{
		HostConfig: hostConfig,
		Client:     newClient,
	})

	return nil
}

func (p *Parallexe) Close() error {
	for _, hostConnection := range p.HostConnections {
		if hostConnection.Client != nil {
			err := hostConnection.Client.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
