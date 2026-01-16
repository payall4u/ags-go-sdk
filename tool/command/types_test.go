package command

import (
	"testing"
)

func TestProcessConfig_loadDefault_NilConfig(t *testing.T) {
	var config *ProcessConfig
	result := config.loadDefault()

	if result == nil {
		t.Fatal("loadDefault() should not return nil")
	}
	if result.User != "user" {
		t.Errorf("expected User to be 'user', got %q", result.User)
	}
}

func TestProcessConfig_loadDefault_NonNilConfig(t *testing.T) {
	config := &ProcessConfig{User: "custom"}
	result := config.loadDefault()

	if result != config {
		t.Fatal("loadDefault() should return the same pointer for non-nil config")
	}
	if result.User != "custom" {
		t.Errorf("expected User to be 'custom', got %q", result.User)
	}
}

func TestProcessConfig_loadDefault_NonNilConfig_EmptyUser(t *testing.T) {
	config := &ProcessConfig{}
	result := config.loadDefault()

	if result != config {
		t.Fatal("loadDefault() should return the same pointer for non-nil config")
	}
	if result.User != "user" {
		t.Errorf("expected User to be 'user', got %q", result.User)
	}
}

func TestOnOutputConfig_loadDefault_NilConfig(t *testing.T) {
	var config *OnOutputConfig
	result := config.loadDefault()

	if result == nil {
		t.Fatal("loadDefault() should not return nil")
	}
	if result.OnStdout == nil {
		t.Error("OnStdout should not be nil")
	}
	if result.OnStderr == nil {
		t.Error("OnStderr should not be nil")
	}

	// Verify callbacks don't panic
	result.OnStdout([]byte("test"))
	result.OnStderr([]byte("test"))
}

func TestOnOutputConfig_loadDefault_NonNilConfig_WithCallbacks(t *testing.T) {
	var stdoutCalled, stderrCalled bool
	config := &OnOutputConfig{
		OnStdout: func([]byte) { stdoutCalled = true },
		OnStderr: func([]byte) { stderrCalled = true },
	}
	result := config.loadDefault()

	if result != config {
		t.Fatal("loadDefault() should return the same pointer for non-nil config")
	}

	result.OnStdout([]byte("test"))
	result.OnStderr([]byte("test"))

	if !stdoutCalled {
		t.Error("OnStdout callback should have been called")
	}
	if !stderrCalled {
		t.Error("OnStderr callback should have been called")
	}
}

func TestOnOutputConfig_loadDefault_NonNilConfig_NilOnStdout(t *testing.T) {
	var stderrCalled bool
	config := &OnOutputConfig{
		OnStdout: nil,
		OnStderr: func([]byte) { stderrCalled = true },
	}
	result := config.loadDefault()

	if result != config {
		t.Fatal("loadDefault() should return the same pointer")
	}
	if result.OnStdout == nil {
		t.Error("OnStdout should be set to default function when nil")
	}

	// Verify default OnStdout doesn't panic
	result.OnStdout([]byte("test"))

	result.OnStderr([]byte("test"))
	if !stderrCalled {
		t.Error("OnStderr callback should have been called")
	}
}

func TestOnOutputConfig_loadDefault_NonNilConfig_NilOnStderr(t *testing.T) {
	var stdoutCalled bool
	config := &OnOutputConfig{
		OnStdout: func([]byte) { stdoutCalled = true },
		OnStderr: nil,
	}
	result := config.loadDefault()

	if result != config {
		t.Fatal("loadDefault() should return the same pointer")
	}
	if result.OnStderr == nil {
		t.Error("OnStderr should be set to default function when nil")
	}

	result.OnStdout([]byte("test"))
	if !stdoutCalled {
		t.Error("OnStdout callback should have been called")
	}

	// Verify default OnStderr doesn't panic
	result.OnStderr([]byte("test"))
}

func TestOnOutputConfig_loadDefault_NonNilConfig_BothCallbacksNil(t *testing.T) {
	config := &OnOutputConfig{
		OnStdout: nil,
		OnStderr: nil,
	}
	result := config.loadDefault()

	if result != config {
		t.Fatal("loadDefault() should return the same pointer")
	}
	if result.OnStdout == nil {
		t.Error("OnStdout should be set to default function when nil")
	}
	if result.OnStderr == nil {
		t.Error("OnStderr should be set to default function when nil")
	}

	// Verify default callbacks don't panic
	result.OnStdout([]byte("test"))
	result.OnStderr([]byte("test"))
}
