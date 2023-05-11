package parallexe

import "testing"

func TestBuildVariables(t *testing.T) {
	hostConfigs := HostConfig{
		SshConfig: nil,
		Host:      "100.0.0.1",
		Groups:    []string{"group1", "group2"},
	}

	t.Run("Build variables", func(t *testing.T) {
		result := buildVariables(hostConfigs, &ExecVariables{
			Variables: KeyValueVariable{
				"var1": "value1",
			},
			GroupVariables: map[string]KeyValueVariable{
				"group1": {
					"var2": "value2",
				},
			},
			HostVariables: map[string]KeyValueVariable{
				"100.0.0.1": {
					"var3": "value3",
				},
			},
		})

		if len(result) != 3 {
			t.Fatalf("Expected 3 variables, got %d", len(result))
		}

		if result["var1"] != "value1" {
			t.Errorf("Expected var1 to be 'value1', got '%s'", result["var1"])
		}

		if result["var2"] != "value2" {
			t.Errorf("Expected var2 to be 'value2', got '%s'", result["var2"])
		}

		if result["var3"] != "value3" {
			t.Errorf("Expected var3 to be 'value3', got '%s'", result["var3"])
		}
	})

	t.Run("Override host variables", func(t *testing.T) {
		result := buildVariables(hostConfigs, &ExecVariables{
			Variables: KeyValueVariable{
				"var1": "value1",
			},
			GroupVariables: map[string]KeyValueVariable{
				"group1": {
					"var2": "value2",
				},
			},
			HostVariables: map[string]KeyValueVariable{
				"100.0.0.1": {
					"var1": "value3",
				},
			},
		})

		if len(result) != 2 {
			t.Fatalf("Expected 2 variables, got %d", len(result))
		}

		if result["var1"] != "value3" {
			t.Errorf("Expected var1 to be 'value3', got '%s'", result["var1"])
		}
	})

	t.Run("Override group variables", func(t *testing.T) {
		result := buildVariables(hostConfigs, &ExecVariables{
			Variables: KeyValueVariable{
				"var1": "value1",
			},
			GroupVariables: map[string]KeyValueVariable{
				"group1": {
					"var1": "value2",
				},
			},
			HostVariables: map[string]KeyValueVariable{
				"100.0.0.1": {
					"var3": "value3",
				},
			},
		})

		if len(result) != 2 {
			t.Fatalf("Expected 2 variables, got %d", len(result))
		}

		if result["var1"] != "value2" {
			t.Errorf("Expected var1 to be 'value2', got '%s'", result["var1"])
		}
	})

	t.Run("Do not override variables", func(t *testing.T) {
		result := buildVariables(hostConfigs, &ExecVariables{
			Variables: KeyValueVariable{
				"var1": "value1",
			},
			GroupVariables: map[string]KeyValueVariable{
				"group3": {
					"var1": "value2",
				},
			},
			HostVariables: map[string]KeyValueVariable{
				"100.0.0.2": {
					"var3": "value3",
				},
			},
		})

		if len(result) != 1 {
			t.Fatalf("Expected 1 variable, got %d", len(result))
		}
		if result["var1"] != "value1" {
			t.Errorf("Expected var1 to be 'value1', got '%s'", result["var1"])
		}
	})
}
