package parallexe

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type SendConfig struct {
	// ExecConfig allows to filter hosts and groups
	ExecConfig *ExecConfig
	// CompileTemplate If sourcePath is a go template, Parallexe will compile this template for each host before sending it
	CompileTemplate bool
	// ExecVariables contains the runtime variables for compiling the templates when sending them.
	// These variables allow you to customize the rendering of templates for each host.
	// Variables can be overridden by host group specific variables and host specific variables.
	ExecVariables *ExecVariables
	// Owner is the owner of the destination file
	Owner string
	//	Group is the group of the destination file
	Mode string
	// IgnoreIfExists indicates whether the upload should be ignored if the destination file already exists on the remote host.
	// If this value is set to true, the upload will not be performed and no error will be returned if the file already exists.
	// If this value is set to false, the upload will be performed even if the file already exists, causing it to be overwritten.
	// The default value is false.
	IgnoreIfExists bool
}

// Send sends a source file to a destination on remote hosts.
// The source file can be a template that will be rendered before sending.
func (p *Parallexe) Send(sourcePath string, destPath string, config *SendConfig) (*CommandResponses, error) {
	file, err := os.Open(sourcePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	txtContent := string(content)
	hostTxtContent := make(map[string]string)

	if config.CompileTemplate {
		// Parse the template
		tmpl, err := template.New(filepath.Base(sourcePath)).Parse(txtContent)
		if err != nil {
			return nil, err
		}

		// Render the template with the provided data per host
		for _, hostConnection := range p.HostConnections {
			variables := buildVariables(hostConnection.HostConfig, config.ExecVariables)

			// Build variables for this host
			var rendered bytes.Buffer
			err = tmpl.Execute(&rendered, variables)
			if err != nil {
				return nil, err
			}

			hostTxtContent[hostConnection.HostConfig.Host] = rendered.String()
		}
	}

	response, err := p.execSend(destPath, hostTxtContent, txtContent, config.IgnoreIfExists, config.ExecConfig)
	if err != nil {
		return response, err
	}

	if config.Owner != "" {
		response, err = p.Exec(fmt.Sprintf("chown %s %s", config.Owner, destPath), config.ExecConfig)
		if err != nil {
			return response, err
		}
	}

	if config.Mode != "" {
		_, err = p.Exec(fmt.Sprintf("chmod %s %s", config.Mode, destPath), config.ExecConfig)
		if err != nil {
			return response, err
		}
	}

	return response, nil
}

// execSend executes the actual send command to the destination path on the remote hosts.
func (p *Parallexe) execSend(destPath string, hostContent map[string]string, content string, ignoreIfExists bool, config *ExecConfig) (*CommandResponses, error) {
	// Check if file already exist and add a condition if we must override it or not
	preCommand := ""
	if ignoreIfExists {
		preCommand = fmt.Sprintf("[ -f '%s' ] || ", destPath)
	}

	// Build the final command if the content is the same for all hosts
	if len(hostContent) == 0 {
		command := fmt.Sprintf("%sprintf '%s' > %s", preCommand, content, destPath)
		return p.Exec(command, config)
	}

	// Build specific content for each host

	// White list HostSession to execute only on desired hosts
	filteredHosts := getFilteredHosts(p.HostConnections, config)

	var wg sync.WaitGroup
	wg.Add(len(filteredHosts))

	commandResponses := make(map[string]*CommandResponse, 0)
	errorHosts := make([]string, 0)

	for _, host := range filteredHosts {
		loopHost := host
		go func() {
			command := fmt.Sprintf("%sprintf '%s' > %s", preCommand, hostContent[loopHost.HostConfig.Host], destPath)
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
