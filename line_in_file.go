package parallexe

import "fmt"

type LineInFileConfig struct {
	// ExecConfig allows to filter hosts and groups
	ExecConfig *ExecConfig
	Absent     bool
}

// LineInFile checks if a line is present in a file.
// If absent is true, the line will be removed from the file if it exists or nothing will be done.
// If absent is false, if file does not exist, it will be created with the line.
func (p *Parallexe) LineInFile(path string, line string, config *LineInFileConfig) (*CommandResponses, error) {
	commands := ""

	if config.Absent {
		// If file exists, remove line from file
		// Otherwise, do nothing
		commands = fmt.Sprintf("sed '/%s/d' %s > %s.tmp 2>/dev/null && mv %s.tmp %s ", line, path, path, path, path)
	} else {
		// Append the line to the file or create the file if it does not exist
		commands = fmt.Sprintf("[ -f %s ] || printf '%s' > %s; grep -qF \"%s\" %s || printf '\n%s' >> %s", path, line, path, line, path, line, path)
	}

	return p.Exec(commands, config.ExecConfig)
}
