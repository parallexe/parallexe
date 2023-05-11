package parallexe

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/exp/slices"
	"os/exec"
	"strings"
	"sync"
)

type ExecConfig struct {
	Hosts  []string
	Groups []string
}

// Exec executes a command on a list of hosts
func (p *Parallexe) Exec(command string, execConfig *ExecConfig) (*CommandResponses, error) {
	// White list HostSession to execute only on desired hosts
	filteredHosts := getFilteredHosts(p.HostConnections, execConfig)

	var wg sync.WaitGroup
	wg.Add(len(filteredHosts))

	commandResponses := make(map[string]*CommandResponse, 0)
	errorHosts := make([]string, 0)

	for _, host := range filteredHosts {
		loopHost := host
		go func() {
			commandResponse := executeCommandOnHost(loopHost, command, &wg)

			commandResponses[loopHost.HostConfig.Host] = commandResponse
			if commandResponse.Error != nil || commandResponse.Stderr != "" {
				errorHosts = append(errorHosts, loopHost.HostConfig.Host)
			}
		}()
	}

	wg.Wait()

	var commandError error
	if len(errorHosts) > 0 {
		commandError = fmt.Errorf("error on hosts: %v", errorHosts)
	}

	return &CommandResponses{HostResponses: commandResponses}, commandError
}

// MultiExec executes a list of commands on a list of hosts
// It returns a list of MultiCommandResponses, each one containing the command and the responses for each host
// If a command fails on a host, the next commands will not be executed on any host
// Commands not executed on any host will have a status CommandStatusSkip
func (p *Parallexe) MultiExec(commands []string, execConfig *ExecConfig) ([]*MultiCommandResponses, error) {
	// White list HostSession to execute only on desired hosts
	filteredHosts := getFilteredHosts(p.HostConnections, execConfig)

	multiCommandResponses := make([]*MultiCommandResponses, 0)

	for _, command := range commands {
		multiCommandResponses = append(multiCommandResponses, &MultiCommandResponses{
			Command:       command,
			Status:        CommandStatusSkip,
			HostResponses: make(map[string]*CommandResponse),
		})
	}

	var m sync.Mutex
	errorHosts := make([]string, 0)

	for index, command := range commands {
		loopCommandIndex := index
		loopCommand := command

		var wg sync.WaitGroup
		wg.Add(len(filteredHosts))

		for _, host := range filteredHosts {
			loopHost := host
			go func() {

				commandResponse := executeCommandOnHost(loopHost, loopCommand, &wg)

				m.Lock()
				multiCommandResponses[loopCommandIndex].HostResponses[loopHost.HostConfig.Host] = commandResponse
				multiCommandResponses[loopCommandIndex].Status = CommandStatusDone
				m.Unlock()

				if commandResponse.Error != nil || commandResponse.Stderr != "" {
					errorHosts = append(errorHosts, loopHost.HostConfig.Host)
				}
			}()

			if len(errorHosts) > 0 {
				break
			}
		}

		if len(errorHosts) > 0 {
			break
		}

		wg.Wait()
	}

	var commandError error
	if len(errorHosts) > 0 {
		commandError = fmt.Errorf("error on hosts: %v", errorHosts)
	}

	return multiCommandResponses, commandError
}

// getFilteredHosts returns a list of HostSession filtered by ExecConfig
func getFilteredHosts(hostConnections []HostConnection, execConfig *ExecConfig) []HostConnection {
	filteredHosts := make([]HostConnection, 0)

	if execConfig == nil || (len(execConfig.Hosts) == 0 && len(execConfig.Groups) == 0) {
		for _, host := range hostConnections {
			filteredHosts = append(filteredHosts, host)
		}
	} else {
		if len(execConfig.Hosts) > 0 {
			for _, hostConnection := range hostConnections {
				if slices.Contains(execConfig.Hosts, hostConnection.HostConfig.Host) {
					filteredHosts = append(filteredHosts, hostConnection)
				}
			}
		}

		if len(execConfig.Groups) > 0 {
			for _, hostConnection := range hostConnections {
				for _, group := range hostConnection.HostConfig.Groups {
					if slices.Contains(execConfig.Groups, group) {
						filteredHosts = append(filteredHosts, hostConnection)
						break
					}
				}
			}
		}
	}

	return filteredHosts
}

// remoteExecute executes a command on a remote host
// If hostSession.Client is nil, run command locally
func executeCommandOnHost(hostSession HostConnection, cmd string, wg *sync.WaitGroup) *CommandResponse {
	defer wg.Done()

	if hostSession.Client == nil {
		return localExecute(cmd)
	}

	// Execute command on remote host
	return remoteExecute(hostSession, cmd)
}

// remoteExecute executes a command on a remote host
func remoteExecute(hostSession HostConnection, cmd string) *CommandResponse {
	var stdout strings.Builder
	var stderr strings.Builder

	session, err := hostSession.Client.NewSession()
	if err != nil {
		return &CommandResponse{
			Stdout:  "",
			Stderr:  "",
			Error:   fmt.Errorf("can't open SSH connection: %v", err),
			Code:    -1,
			Success: false,
		}
	}
	defer session.Close()

	session.Stdout = &stdout
	session.Stderr = &stderr

	var code int
	err = session.Run(cmd)

	if exitErr, ok := err.(*ssh.ExitError); ok {
		code = exitErr.ExitStatus()
	} else if err != nil {
		return &CommandResponse{
			Stdout:  "",
			Stderr:  "",
			Error:   err,
			Code:    -1,
			Success: false,
		}
	}

	return &CommandResponse{
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
		Error:   nil,
		Code:    code,
		Success: err == nil && stderr.String() == "",
	}
}

// localExecute executes a command locally
func localExecute(cmd string) *CommandResponse {
	var stdout strings.Builder
	var stderr strings.Builder

	command := exec.Command("sh", "-c", cmd)
	command.Stdout = &stdout
	command.Stderr = &stderr

	var code int
	err := command.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		code = exitErr.ExitCode()
	} else if err != nil {
		return &CommandResponse{
			Stdout:  "",
			Stderr:  "",
			Error:   err,
			Code:    -1,
			Success: false,
		}
	}

	return &CommandResponse{
		Stdout:  stdout.String(),
		Stderr:  stderr.String(),
		Error:   nil,
		Code:    code,
		Success: err == nil && stderr.String() == "",
	}
}
