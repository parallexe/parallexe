package parallexe

import "strings"

type CommandResponse struct {
	Stdout string
	// Contains the error returns by the command
	Stderr string
	// Contains network error (ssh connection, ...)
	Error   error
	Code    int
	Success bool
}

type CommandStatus string

const (
	CommandStatusDone CommandStatus = "done"
	CommandStatusSkip CommandStatus = "skip"
)

type MultiCommandResponses struct {
	Status        CommandStatus
	Command       string
	HostResponses map[string]*CommandResponse
}

type CommandResponses struct {
	HostResponses map[string]*CommandResponse
}

func (r *CommandResponses) GetStdoutLines() map[string][]string {
	lineHosts := make(map[string][]string, 0)

	for host, commandResponse := range r.HostResponses {
		lineHosts[host] = splitLines(commandResponse.Stdout)
	}

	return lineHosts
}

func (r *CommandResponses) GetStderrLines() map[string][]string {
	lineHosts := make(map[string][]string, 0)

	for host, commandResponse := range r.HostResponses {
		lineHosts[host] = splitLines(commandResponse.Stderr)
	}

	return lineHosts
}

func splitLines(s string) []string {
	lines := make([]string, 0)

	for _, line := range strings.Split(s, "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines
}
