//go:build internal
// +build internal

package test_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/code"
)

func TestCreate_WithInternalDomain_Integration(t *testing.T) {
	client := newAgsClient(t)
	if client == nil {
		t.Skip("missing cloud credentials, skip integration tests")
	}

	sb, err := code.Create(context.TODO(), "code-interpreter-v1",
		code.WithClient(client),
		code.WithSandboxTimeout(300*time.Second),
		code.WithDataPlaneDomain("internal.tencentags.com"),
	)
	if err != nil {
		t.Fatalf("create sandbox with internal domain: %v", err)
	}
	defer sb.Kill(context.TODO())

	// Verify the domain is set correctly
	if !strings.HasSuffix(sb.ConnectionConfig.Domain, "internal.tencentags.com") {
		t.Errorf("expected domain to end with 'internal.tencentags.com', got '%s'", sb.ConnectionConfig.Domain)
	}

	// Verify GetHost works correctly
	host := sb.GetHost(8080)
	if !strings.HasSuffix(host, "internal.tencentags.com") {
		t.Errorf("expected host to end with 'internal.tencentags.com', got '%s'", host)
	}

	// Test that the sandbox is functional with internal domain
	exec, err := sb.Code.RunCode(context.TODO(), "print('hello from internal domain')", nil, nil)
	if err != nil {
		t.Fatalf("RunCode with internal domain error: %v", err)
	}
	if exec.Error != nil {
		t.Fatalf("unexpected execution error: %+v", exec.Error)
	}
}

func TestConnect_WithInternalDomain_Integration(t *testing.T) {
	client := newAgsClient(t)
	if client == nil {
		t.Skip("missing cloud credentials, skip integration tests")
	}

	// First create a sandbox
	sb, err := code.Create(context.TODO(), "code-interpreter-v1",
		code.WithClient(client),
		code.WithSandboxTimeout(300*time.Second),
	)
	if err != nil {
		t.Fatalf("create sandbox: %v", err)
	}
	defer sb.Kill(context.TODO())

	sandboxId := sb.SandboxId

	// Connect with internal domain
	sbConnected, err := code.Connect(context.TODO(), sandboxId,
		code.WithClient(client),
		code.WithDataPlaneDomain("internal.tencentags.com"),
	)
	if err != nil {
		t.Fatalf("connect to sandbox with internal domain: %v", err)
	}

	// Verify the domain is set correctly
	if !strings.HasSuffix(sbConnected.ConnectionConfig.Domain, "internal.tencentags.com") {
		t.Errorf("expected domain to end with 'internal.tencentags.com', got '%s'", sbConnected.ConnectionConfig.Domain)
	}

	// Test that the connected sandbox is functional
	exec, err := sbConnected.Code.RunCode(context.TODO(), "print('hello from connected sandbox')", nil, nil)
	if err != nil {
		t.Fatalf("RunCode on connected sandbox error: %v", err)
	}
	if exec.Error != nil {
		t.Fatalf("unexpected execution error: %+v", exec.Error)
	}
}