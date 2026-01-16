package test_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/code"
)

func TestCreate_WithDefaultDomain_Integration(t *testing.T) {
	client := newAgsClient(t)
	if client == nil {
		t.Skip("missing cloud credentials, skip integration tests")
	}

	sb, err := code.Create(context.TODO(), "code-interpreter-v1",
		code.WithClient(client),
		code.WithSandboxTimeout(300*time.Second),
	)
	if err != nil {
		t.Fatalf("create sandbox with default domain: %v", err)
	}
	defer sb.Kill(context.TODO())

	// Verify the domain is set correctly
	if sb.ConnectionConfig.Domain == "" {
		t.Fatal("expected non-empty domain")
	}
	if !strings.HasSuffix(sb.ConnectionConfig.Domain, "tencentags.com") {
		t.Errorf("expected domain to end with 'tencentags.com', got '%s'", sb.ConnectionConfig.Domain)
	}

	// Verify GetHost works correctly
	host := sb.GetHost(8080)
	if !strings.Contains(host, sb.SandboxId) {
		t.Errorf("expected host to contain sandbox ID '%s', got '%s'", sb.SandboxId, host)
	}
	if !strings.HasSuffix(host, "tencentags.com") {
		t.Errorf("expected host to end with 'tencentags.com', got '%s'", host)
	}
}
