package core

import (
	"testing"

	"github.com/TencentCloudAgentRuntime/ags-go-sdk/connection"
)

func TestGetHost_WithDefaultDomain(t *testing.T) {
	c := &Core{
		SandboxId: "test-sandbox-id",
		ConnectionConfig: &connection.Config{
			Domain: "ap-guangzhou.tencentags.com",
		},
	}

	host := c.GetHost(8080)
	expected := "8080-test-sandbox-id.ap-guangzhou.tencentags.com"

	if host != expected {
		t.Errorf("expected host '%s', got '%s'", expected, host)
	}
}

func TestGetHost_WithCustomDomain(t *testing.T) {
	c := &Core{
		SandboxId: "test-sandbox-id",
		ConnectionConfig: &connection.Config{
			Domain: "ap-guangzhou.internal.tencentags.com",
		},
	}

	host := c.GetHost(3000)
	expected := "3000-test-sandbox-id.ap-guangzhou.internal.tencentags.com"

	if host != expected {
		t.Errorf("expected host '%s', got '%s'", expected, host)
	}
}

func TestGetHost_DifferentPorts(t *testing.T) {
	c := &Core{
		SandboxId: "sandbox-123",
		ConnectionConfig: &connection.Config{
			Domain: "ap-shanghai.tencentags.com",
		},
	}

	tests := []struct {
		port     int
		expected string
	}{
		{80, "80-sandbox-123.ap-shanghai.tencentags.com"},
		{443, "443-sandbox-123.ap-shanghai.tencentags.com"},
		{8080, "8080-sandbox-123.ap-shanghai.tencentags.com"},
		{3000, "3000-sandbox-123.ap-shanghai.tencentags.com"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			host := c.GetHost(tt.port)
			if host != tt.expected {
				t.Errorf("expected host '%s', got '%s'", tt.expected, host)
			}
		})
	}
}
