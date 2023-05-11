package parallexe

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestLineInFileAddLine(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Error during file test creation: %v", err)
	}
	defer os.Remove(file.Name())

	file.WriteString("toto\n")

	pexe, err := New([]HostConfig{{Host: "localhost"}})
	if err != nil {
		t.Fatalf("Error during Parallexe creation: %v", err)
	}
	defer pexe.Close()

	response, err := pexe.LineInFile(file.Name(), "tata", &LineInFileConfig{
		ExecConfig: nil,
		Absent:     false,
	})
	if err != nil {
		t.Fatalf("Error during LineInFile: %v", err)
	}

	if response == nil {
		t.Fatalf("Response is nil")
	}

	if response.HostResponses == nil {
		t.Fatalf("HostResponses is nil")
	}

	if len(response.HostResponses) != 1 {
		t.Fatalf("HostResponses length is not 1")
	}

	content, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("Error during file test reading: %v", err)
	}

	// Check if temp file contains "tata"
	if !strings.Contains(string(content), "tata") {
		t.Fatalf("File content is not correct")
	}

	// Check if temp file still contains "toto"
	if !strings.Contains(string(content), "toto") {
		t.Fatalf("File content is not correct")
	}

	// Check if temp file contains "tata" twice
	response, err = pexe.LineInFile(file.Name(), "tata", &LineInFileConfig{
		ExecConfig: nil,
		Absent:     false,
	})
	if err != nil {
		t.Fatalf("Error during LineInFile : %v", err)
	}

	content, err = os.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("Error during file test reading : %v", err)
	}
	if strings.Contains(string(content), "tata\ntata") {
		t.Fatalf("File content is not correct")
	}

	// Test against a file that does not exist (in tmp directory)
	fakeFile := fmt.Sprintf("%s/%s", os.TempDir(), "doesnotexist")
	defer os.Remove(fakeFile)

	response, err = pexe.LineInFile(fakeFile, "tata", &LineInFileConfig{
		ExecConfig: nil,
		Absent:     false,
	})
	if err != nil {
		t.Fatalf("Error during LineInFile : %v", err)
	}

	if response == nil {
		t.Fatalf("Response is nil")
	}

	// Check if file still does not exist
	if _, err := os.Stat(fakeFile); os.IsNotExist(err) {
		t.Fatalf("File should exist")
	}
}

func TestLineInFileRemoveLine(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Error during file test creation : %v", err)
	}
	defer os.Remove(file.Name())

	file.WriteString("toto\n")

	pexe, err := New([]HostConfig{{Host: "localhost"}})
	if err != nil {
		t.Fatalf("Error during Parallexe creation : %v", err)
	}
	defer pexe.Close()

	response, err := pexe.LineInFile(file.Name(), "toto", &LineInFileConfig{
		ExecConfig: nil,
		Absent:     true,
	})
	if err != nil {
		t.Fatalf("Error during LineInFile : %v", err)
	}

	if response == nil {
		t.Fatalf("Response is nil")
	}

	// Check if temp file contains "toto"
	content, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("Error during file test reading : %v", err)
	}
	if strings.Contains(string(content), "toto") {
		t.Fatalf("File content is not correct")
	}

	// Execute lineInFile on a file that does not exist
	fakeFile := fmt.Sprintf("%s/%s", os.TempDir(), "doesnotexist")
	response, err = pexe.LineInFile(fakeFile, "toto", &LineInFileConfig{
		ExecConfig: nil,
		Absent:     true,
	})
	if err != nil {
		t.Fatalf("Error during LineInFile : %v", err)
	}
	if response == nil {
		t.Fatalf("Response is nil")
	}

	// Check if file still does not exist
	if _, err := os.Stat(fakeFile); !os.IsNotExist(err) {
		t.Fatalf("File should not exist")
	}
}
