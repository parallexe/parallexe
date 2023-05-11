package parallexe

import (
	"testing"
)

func TestGetFilteredHosts(t *testing.T) {
	hostConnections := []HostConnection{
		{
			HostConfig: HostConfig{
				SshConfig: nil,
				Host:      "100.0.0.1",
				Groups:    []string{"group1", "group2"},
			},
			Client: nil,
		},
		{
			HostConfig: HostConfig{
				SshConfig: nil,
				Host:      "100.0.0.2",
				Groups:    []string{"group1", "group3"},
			},
			Client: nil,
		},
		{
			HostConfig: HostConfig{
				SshConfig: nil,
				Host:      "100.0.0.3",
				Groups:    []string{"group3"},
			},
			Client: nil,
		},
		{
			HostConfig: HostConfig{
				SshConfig: nil,
				Host:      "100.0.0.4",
				Groups:    []string{"group1", "group2", "group3"},
			},
			Client: nil,
		},
		{
			HostConfig: HostConfig{
				SshConfig: nil,
				Host:      "100.0.0.5",
				Groups:    []string{"group1", "group2", "group3"},
			},
			Client: nil,
		},
	}

	t.Run("filter by host", func(t *testing.T) {
		// Should return 1 host
		result := getFilteredHosts(hostConnections, &ExecConfig{
			Hosts:  []string{"100.0.0.4"},
			Groups: nil,
		})
		if len(result) != 1 {
			t.Errorf("Expected 1 host, got %d", len(result))
		}
		if result[0].HostConfig.Host != "100.0.0.4" {
			t.Errorf("Wrong host returned, expected 100.0.0.4 got %s", result[0].HostConfig.Host)
		}

		// Should return 2 hosts
		result = getFilteredHosts(hostConnections, &ExecConfig{
			Hosts:  []string{"100.0.0.4", "100.0.0.5"},
			Groups: nil,
		})
		if len(result) != 2 {
			t.Errorf("Expected 2 hosts, got %d", len(result))
		}
		if result[0].HostConfig.Host != "100.0.0.4" {
			t.Errorf("Wrong host returned, expected 100.0.0.4 got %s", result[0].HostConfig.Host)
		}
		if result[1].HostConfig.Host != "100.0.0.5" {
			t.Errorf("Wrong host returned, expected 100.0.0.5 got %s", result[0].HostConfig.Host)
		}

		// Should return 1 host if second host is invalid
		result = getFilteredHosts(hostConnections, &ExecConfig{
			Hosts:  []string{"100.0.0.4", "100.0.1.0"},
			Groups: nil,
		})
		if len(result) != 1 {
			t.Errorf("Expected 1 host, got %d", len(result))
		}
		if result[0].HostConfig.Host != "100.0.0.4" {
			t.Errorf("Wrong host returned, expected 100.0.0.4 got %s", result[0].HostConfig.Host)
		}
	})

	t.Run("filter by group", func(t *testing.T) {
		// Should return 4 hosts
		result := getFilteredHosts(hostConnections, &ExecConfig{
			Hosts:  nil,
			Groups: []string{"group1"},
		})
		if len(result) != 4 {
			t.Errorf("Expected 4 hosts, got %d", len(result))
		}

		// Should return 3 hosts
		result = getFilteredHosts(hostConnections, &ExecConfig{
			Hosts:  nil,
			Groups: []string{"group2"},
		})
		if len(result) != 3 {
			t.Errorf("Expected 3 hosts, got %d", len(result))
		}

		// Should return 4 hosts
		result = getFilteredHosts(hostConnections, &ExecConfig{
			Hosts:  nil,
			Groups: []string{"group3"},
		})
		if len(result) != 4 {
			t.Errorf("Expected 4 hosts, got %d", len(result))
		}

		// Should return 5 hosts
		result = getFilteredHosts(hostConnections, &ExecConfig{
			Hosts:  nil,
			Groups: []string{"group1", "group3"},
		})
		if len(result) != 5 {
			t.Errorf("Expected 5 hosts, got %d", len(result))
		}
	})

	t.Run("filter by host and group", func(t *testing.T) {
		// Should return 4 hosts
		result := getFilteredHosts(hostConnections, &ExecConfig{
			Hosts:  []string{"100.0.0.2"},
			Groups: []string{"group2"},
		})
		if len(result) != 4 {
			t.Errorf("Expected 4 hosts, got %d", len(result))
		}
	})
}
