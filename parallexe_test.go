package parallexe

import (
	"fmt"
	"testing"
)

func TestParallexe(t *testing.T) {
	pexe, err := New([]HostConfig{{
		Host: "localhost",
		SshConfig: &SshConfig{
			User: "root",
		}},
	})

	if err != nil {
		panic(err)
	}

	defer pexe.Close()

	response, err := pexe.LineInFile("./tutu.txt", "tata", &LineInFileConfig{
		ExecConfig: nil,
		Absent:     false,
	})
	fmt.Printf("%+v", response)
	fmt.Printf("%+v", err)
}
