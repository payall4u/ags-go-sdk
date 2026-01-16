package command

// ProcessConfig holds configuration for process execution
type ProcessConfig struct {
	User string
	Args []string
	Envs map[string]string
	Cwd  *string
}

func (config *ProcessConfig) loadDefault() *ProcessConfig {
	if config == nil {
		config = &ProcessConfig{}
	}
	if config.User == "" {
		config.User = "user"
	}
	return config
}

// ProcessInfo contains information about a running process
type ProcessInfo struct {
	Pid  uint32
	Tag  *string
	Cmd  string
	Args []string
	Envs map[string]string
	Cwd  *string
}

// ProcessResult contains the result of process execution
type ProcessResult struct {
	ExitCode int32
	Error    *string
}

// Result contains aggregated output for foreground execution
type Result struct {
	ExitCode int32
	Stdout   []byte
	Stderr   []byte
	Error    *string
}

// OnOutputConfig holds callbacks for process output
type OnOutputConfig struct {
	OnStdout func([]byte)
	OnStderr func([]byte)
}

func (config *OnOutputConfig) loadDefault() *OnOutputConfig {
	if config == nil {
		config = &OnOutputConfig{
			OnStdout: func([]byte) {},
			OnStderr: func([]byte) {},
		}
	}
	if config.OnStdout == nil {
		config.OnStdout = func([]byte) {}
	}
	if config.OnStderr == nil {
		config.OnStderr = func([]byte) {}
	}
	return config
}
