package filesystem

import "fmt"

// User constants for filesystem operations
const (
	UserDefault = "user"
	UserRoot    = "root"
)

// ReadConfig holds configuration for reading files
type ReadConfig struct{ User string }

func (config *ReadConfig) valid() (*ReadConfig, error) {
	if config == nil {
		config = &ReadConfig{}
	}
	if config.User == "" {
		config.User = UserDefault
	}
	if config.User != UserDefault && config.User != UserRoot {
		return nil, fmt.Errorf("invalid user: %s, must be user or root", config.User)
	}
	return config, nil
}

// WriteConfig holds configuration for writing files
type WriteConfig struct{ User string }

func (config *WriteConfig) valid() (*WriteConfig, error) {
	if config == nil {
		config = &WriteConfig{}
	}
	// Set default before validation for consistency
	if config.User == "" {
		config.User = UserDefault
	}
	if config.User != UserDefault && config.User != UserRoot {
		return nil, fmt.Errorf("invalid user: %s, must be user or root", config.User)
	}
	return config, nil
}

// ListConfig holds configuration for listing directory entries
type ListConfig struct {
	Depth int
	User  string
}

func (config *ListConfig) valid() (*ListConfig, error) {
	if config == nil {
		config = &ListConfig{}
	}
	if config.User == "" {
		config.User = UserDefault
	}
	if config.User != UserDefault && config.User != UserRoot {
		return nil, fmt.Errorf("invalid user: %s, must be user or root", config.User)
	}
	if config.Depth < 0 {
		return nil, fmt.Errorf("invalid depth: %d, must be >= 0", config.Depth)
	}
	return config, nil
}

// ExistsConfig holds configuration for checking file existence
type ExistsConfig struct{ User string }

func (config *ExistsConfig) valid() (*ExistsConfig, error) {
	if config == nil {
		config = &ExistsConfig{}
	}
	if config.User == "" {
		config.User = UserDefault
	}
	if config.User != UserDefault && config.User != UserRoot {
		return nil, fmt.Errorf("invalid user: %s, must be user or root", config.User)
	}
	return config, nil
}

// GetInfoConfig holds configuration for getting file information
type GetInfoConfig struct{ User string }

func (config *GetInfoConfig) valid() (*GetInfoConfig, error) {
	if config == nil {
		config = &GetInfoConfig{}
	}
	if config.User == "" {
		config.User = UserDefault
	}
	if config.User != UserDefault && config.User != UserRoot {
		return nil, fmt.Errorf("invalid user: %s, must be user or root", config.User)
	}
	return config, nil
}

// RemoveConfig holds configuration for removing files
type RemoveConfig struct{ User string }

func (config *RemoveConfig) valid() (*RemoveConfig, error) {
	if config == nil {
		config = &RemoveConfig{}
	}
	if config.User == "" {
		config.User = UserDefault
	}
	if config.User != UserDefault && config.User != UserRoot {
		return nil, fmt.Errorf("invalid user: %s, must be user or root", config.User)
	}
	return config, nil
}

// RenameConfig holds configuration for renaming files
type RenameConfig struct{ User string }

func (config *RenameConfig) valid() (*RenameConfig, error) {
	if config == nil {
		config = &RenameConfig{}
	}
	if config.User == "" {
		config.User = UserDefault
	}
	if config.User != UserDefault && config.User != UserRoot {
		return nil, fmt.Errorf("invalid user: %s, must be user or root", config.User)
	}
	return config, nil
}

// MakeDirConfig holds configuration for creating directories
type MakeDirConfig struct{ User string }

func (config *MakeDirConfig) valid() (*MakeDirConfig, error) {
	if config == nil {
		config = &MakeDirConfig{}
	}
	if config.User == "" {
		config.User = UserDefault
	}
	if config.User != UserDefault && config.User != UserRoot {
		return nil, fmt.Errorf("invalid user: %s, must be user or root", config.User)
	}
	return config, nil
}
