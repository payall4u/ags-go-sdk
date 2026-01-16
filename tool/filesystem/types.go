package filesystem

import (
	"time"
)

// FileType represents the type of a filesystem entry
type FileType string

const (
	File FileType = "file"
	Dir  FileType = "dir"
)

// WriteInfo contains basic information about a written file
type WriteInfo struct {
	Name string    `json:"name"`
	Type *FileType `json:"type,omitempty"`
	Path string    `json:"path"`
}

// EntryInfo contains detailed information about a filesystem entry
type EntryInfo struct {
	WriteInfo
	Size          int64     `json:"size"`
	Mode          int       `json:"mode"`
	Permissions   string    `json:"permissions"`
	Owner         string    `json:"owner"`
	Group         string    `json:"group"`
	ModifiedTime  time.Time `json:"modified_time"`
	SymlinkTarget *string   `json:"symlink_target,omitempty"`
}

// FilesystemEventType represents the type of filesystem event
type FilesystemEventType string

const (
	EventCreate FilesystemEventType = "create"
	EventWrite  FilesystemEventType = "write"
	EventRemove FilesystemEventType = "remove"
	EventRename FilesystemEventType = "rename"
	EventChmod  FilesystemEventType = "chmod"
)

// FilesystemEvent represents a filesystem change event
type FilesystemEvent struct {
	Name string              `json:"name"`
	Type FilesystemEventType `json:"type"`
}
