package code

import (
	"time"
)

// ExecutionError represents an error that occurred during code execution
type ExecutionError struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Traceback string `json:"traceback"`
}

// Logs contains stdout and stderr output
type Logs struct {
	Stdout []string `json:"stdout"`
	Stderr []string `json:"stderr"`
}

// Result contains multi-format execution result data
type Result struct {
	Text         *string        `json:"text,omitempty"`
	Html         *string        `json:"html,omitempty"`
	Markdown     *string        `json:"markdown,omitempty"`
	Svg          *string        `json:"svg,omitempty"`
	Png          *string        `json:"png,omitempty"`
	Jpeg         *string        `json:"jpeg,omitempty"`
	Pdf          *string        `json:"pdf,omitempty"`
	Latex        *string        `json:"latex,omitempty"`
	Json         map[string]any `json:"json,omitempty"`
	Javascript   *string        `json:"javascript,omitempty"`
	Data         map[string]any `json:"data,omitempty"`
	Chart        map[string]any `json:"chart,omitempty"`
	IsMainResult bool           `json:"is_main_result"`
	Extra        map[string]any `json:"extra,omitempty"`
}

// Execution aggregates results, logs, errors, and execution count
type Execution struct {
	Results        []Result        `json:"results"`
	Logs           Logs            `json:"logs"`
	Error          *ExecutionError `json:"error,omitempty"`
	ExecutionCount *int            `json:"execution_count,omitempty"`
}

// CodeContext represents a code execution context
type CodeContext struct {
	Id       string `json:"id"`
	Language string `json:"language"`
	Cwd      string `json:"cwd"`
}

// createContextRequest is the request body for creating a context
type createContextRequest struct {
	Cwd      *string `json:"cwd,omitempty"`
	Language *string `json:"language,omitempty"`
}

// executeRequest is the request body for code execution
type executeRequest struct {
	Code      string            `json:"code"`
	ContextId *string           `json:"context_id,omitempty"`
	Language  *string           `json:"language,omitempty"`
	EnvVars   map[string]string `json:"env_vars,omitempty"`
}

// RequestTimeouts holds timeout configuration for requests
type RequestTimeouts struct {
	Timeout time.Duration
}

// OnOutputConfig holds callbacks for execution output
// Allows real-time streaming callbacks while preserving aggregated logs
// in the returned Execution object.
type OnOutputConfig struct {
	OnStdout func(string)
	OnStderr func(string)
}

func (config *OnOutputConfig) loadDefault() *OnOutputConfig {
	if config == nil {
		config = &OnOutputConfig{
			OnStdout: func(string) {},
			OnStderr: func(string) {},
		}
	}
	if config.OnStderr == nil {
		config.OnStderr = func(string) {}
	}
	if config.OnStdout == nil {
		config.OnStdout = func(string) {}
	}
	return config
}
