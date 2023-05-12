package parallexe

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestSend(t *testing.T) {
	pexe, err := New([]HostConfig{{Host: "localhost"}})
	if err != nil {
		t.Fatalf("Error during Parallexe creation: %v", err)
	}
	defer pexe.Close()

	t.Run("Send file", func(t *testing.T) {
		file, err := os.CreateTemp("", "testfile")
		if err != nil {
			t.Fatalf("Error during file test creation: %v", err)
		}
		defer os.Remove(file.Name())
		file.WriteString("toto\n")

		// Check if file is correctly sent
		destCopyFile := fmt.Sprintf("%s/%s", os.TempDir(), "file")
		defer os.Remove(destCopyFile)

		response, err := pexe.Send(file.Name(), destCopyFile, &SendConfig{
			ExecConfig:      nil,
			CompileTemplate: false,
			IgnoreIfExists:  false,
		})
		if err != nil {
			t.Fatalf("Error during Send: %v", err)
		}

		if response == nil {
			t.Fatalf("Response is nil")
		}

		// Check if file is correctly sent
		content, err := os.ReadFile(destCopyFile)
		if err != nil {
			t.Fatalf("Error during file test reading: %v", err)
		}
		// Check if temp file contains "toto"
		if string(content) != "toto\n" {
			t.Fatalf("File content is not correct")
		}
	})

	t.Run("Ignore if exists", func(t *testing.T) {
		file, err := os.CreateTemp("", "testfile")
		if err != nil {
			t.Fatalf("Error during file test creation: %v", err)
		}
		defer os.Remove(file.Name())
		file.WriteString("toto\n")

		// Create a file with the same name
		destFile, err := os.CreateTemp("", "fileExists")
		if err != nil {
			t.Fatalf("Error during file test creation: %v", err)
		}
		defer os.Remove(destFile.Name())

		// Add some content
		destFile.WriteString("tata\n")

		response, err := pexe.Send(file.Name(), destFile.Name(), &SendConfig{
			ExecConfig:      nil,
			CompileTemplate: false,
			IgnoreIfExists:  true,
		})
		if err != nil {
			t.Fatalf("Error during Send: %v", err)
		}
		if response == nil {
			t.Fatalf("Response is nil")
		}

		// Read destFile and check if content is still the same
		content, err := os.ReadFile(destFile.Name())
		if err != nil {
			t.Fatalf("Error during file test reading: %v", err)
		}
		// Check if temp file contains "tata"
		if strings.Contains(string(content), "toto") {
			t.Fatalf("File content is not correct")
		}
	})

	t.Run("Compile template", func(t *testing.T) {
		file, err := os.CreateTemp("", "testfile")
		if err != nil {
			t.Fatalf("Error during file test creation: %v", err)
		}
		defer os.Remove(file.Name())

		// Add some content
		file.WriteString("{{ .Name }}\n")

		// Check if file is correctly sent
		destCopyFile := fmt.Sprintf("%s/%s", os.TempDir(), "file")
		defer os.Remove(destCopyFile)

		response, err := pexe.Send(file.Name(), destCopyFile, &SendConfig{
			ExecConfig:      nil,
			CompileTemplate: true,
			ExecVariables:   &ExecVariables{Variables: KeyValueVariable{"Name": "tutu"}},
			Owner:           "",
			Mode:            "",
			IgnoreIfExists:  false,
		})
		if err != nil {
			t.Fatalf("Error during Send: %v", err)
		}
		if response == nil {
			t.Fatalf("Response is nil")
		}

		content, err := os.ReadFile(destCopyFile)
		if err != nil {
			t.Fatalf("Error during file test reading: %v", err)
		}
		// Check if temp file contains "tutu"
		if string(content) != "tutu\n" {
			t.Fatalf("File content is not correct")
		}

	})
}

func ExampleParallexe_Send_copy() {
	pexe, err := New([]HostConfig{{Host: "localhost"}})
	if err != nil {
		panic(err)
	}
	defer pexe.Close()

	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())
	file.WriteString("toto\n")

	// Check if file is correctly sent
	destCopyFile := fmt.Sprintf("%s/%s", os.TempDir(), "destFile")
	defer os.Remove(destCopyFile)

	_, err = pexe.Send(file.Name(), destCopyFile, &SendConfig{
		ExecConfig:      nil,
		CompileTemplate: false,
		IgnoreIfExists:  false,
	})
	if err != nil {
		panic(err)
	}

	content, err := os.ReadFile(destCopyFile)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(content))
	// Output: toto
}

func ExampleParallexe_Send_template() {
	pexe, err := New([]HostConfig{{Host: "localhost"}})
	if err != nil {
		panic(err)
	}
	defer pexe.Close()

	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())
	file.WriteString("{{ .Name }}\n")

	// Check if file is correctly sent
	destCopyFile := fmt.Sprintf("%s/%s", os.TempDir(), "destFile")
	defer os.Remove(destCopyFile)

	_, err = pexe.Send(file.Name(), destCopyFile, &SendConfig{
		ExecConfig:      nil,
		CompileTemplate: true,
		ExecVariables:   &ExecVariables{Variables: KeyValueVariable{"Name": "tutu"}},
	})
	if err != nil {
		panic(err)
	}

	content, err := os.ReadFile(destCopyFile)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(content))
	// Output: tutu
}
