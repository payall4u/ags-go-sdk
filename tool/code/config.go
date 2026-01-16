package code

// RunCodeConfig holds configuration for code execution
type RunCodeConfig struct {
	Language  string
	ContextId string
	Envs      map[string]string
	// Maximum buffer size for scanning response lines, defaults to 1MB
	// Note: MaxBufferSize smaller than 64KB will be ignored
	MaxBufferSize int
}

func (config *RunCodeConfig) loadDefault() *RunCodeConfig {
	if config == nil {
		config = &RunCodeConfig{}
	}
	if config.Language == "" && config.ContextId == "" {
		// Default to Python when no specific context is provided
		config.Language = "python"
	}
	if config.Envs == nil {
		config.Envs = make(map[string]string)
	}
	if config.MaxBufferSize <= 0 {
		config.MaxBufferSize = 1 << 20 // 1MB
	}
	return config
}

// CreateCodeContextConfig holds configuration for creating a code context
type CreateCodeContextConfig struct {
	Cwd      string
	Language string
}

func (config *CreateCodeContextConfig) loadDefault() *CreateCodeContextConfig {
	if config == nil {
		config = &CreateCodeContextConfig{}
	}
	if config.Cwd == "" {
		config.Cwd = "/home/user"
	}
	if config.Language == "" {
		config.Language = "python"
	}
	return config
}
