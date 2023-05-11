# Parallexe


Parallexe executes scripts in parallel on remote servers. It simplifies management by allowing users to execute commands described in a markdown file or directly from code.

## Installation

```bash
go get github.com/parallexe/parallexe
```

## Usage

### Code

```go
package main

import (
    "log"
    "github.com/parallexe/parallexe"
)

func main() {
	hostConfigs := []parallexe.HostConfig{
		{
			Host:   "53.0.0.1",
			Groups: []string{"prod"},
			SshConfig: &parallexe.SshConfig{
				User:           "root",
				PrivateKeyPath: "/home/user/.ssh/id_rsa",
			},
		},
		{
			Host:   "53.0.0.2",
			Groups: []string{"uat"},
			SshConfig: &parallexe.SshConfig{
				User:           "root",
				PrivateKeyPath: "/home/user/.ssh/id_rsa",
			},
		},
	}

	pexe, err := parallexe.New(hostConfigs)
	if err != nil {
		panic(err)
	}
	defer p.Close()

	responses, err := pexe.Exec("ls -l", nil)
	if err != nil {
		panic(err)
	}
	log.Print(responses.GetStdoutLines())
	
	// Send file.tpl to /tmp/file.txt on hosts and compile go template
	_, err := pexe.Send("./file.tpl", "/tmp/file.txt", &parallexe.SendConfig{
		CompileTemplate: true,
		ExecVariables:   &parallexe.ExecVariables{
			Variables: parallexe.KeyValueVariable{
				"var1": "value1",
			},
		},
	})
}
```

### Command line

COMING SOON