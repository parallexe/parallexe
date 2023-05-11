package parallexe

type KeyValueVariable map[string]interface{}

type ExecVariables struct {
	// Variables will be injected for all hosts
	Variables KeyValueVariable
	// GroupVariables will be injected for corresponding group and will override Variables
	GroupVariables map[string]KeyValueVariable
	// HostVariables will be injected in template and will override GroupVariables
	HostVariables map[string]KeyValueVariable
}

// buildVariables will create a map[string]interface{} with all variables defined in execVariables
// This function will override variables:
// execVariables.HostVariables will override execVariables.GroupVariables
// execVariables.GroupVariables will override execVariables.Variables
func buildVariables(hostConfig HostConfig, execVariables *ExecVariables) map[string]interface{} {
	if execVariables == nil {
		return make(map[string]interface{})
	}

	variables := make(KeyValueVariable)

	// Build group variables
	groupVariables := make(KeyValueVariable)
	for _, group := range hostConfig.Groups {
		for key, value := range execVariables.GroupVariables[group] {
			groupVariables[key] = value
		}
	}

	// Build host variables
	hostVariables := execVariables.HostVariables[hostConfig.Host]

	// Merge variables
	mergeVariables(variables, execVariables.Variables)
	mergeVariables(variables, groupVariables)
	mergeVariables(variables, hostVariables)

	return variables
}

// mergeVariables merges source map into destination map
func mergeVariables(destination, source KeyValueVariable) {
	for key, value := range source {
		destination[key] = value
	}
}
