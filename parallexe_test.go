package parallexe

import (
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Test with empty host config", func(t *testing.T) {
		_, err := New([]HostConfig{})
		if err != nil {
			t.Fatalf("Error during Parallexe creation: %v", err)
		}
	})

	t.Run("Test with one host config", func(t *testing.T) {
		pexe, err := New([]HostConfig{{Host: "localhost"}})
		if err != nil {
			t.Fatalf("Error during Parallexe creation: %v", err)
		}

		if pexe == nil {
			t.Fatalf("Parallexe is nil")
		}

		if len(pexe.HostConnections) != 1 {
			t.Fatalf("Parallexe hosts length is not 1")
		}
	})
}
