package test_test

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/code"
	ags "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

// newAgsClient creates an AGS client from environment credentials.
// Returns nil if credentials are missing.
func newAgsClient(t *testing.T) *ags.Client {
	t.Helper()

	secretID := os.Getenv("TENCENTCLOUD_SECRET_ID")
	secretKey := os.Getenv("TENCENTCLOUD_SECRET_KEY")
	if secretID == "" || secretKey == "" {
		return nil
	}

	cred := &common.Credential{
		SecretId:  secretID,
		SecretKey: secretKey,
	}
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "ags.tencentcloudapi.com"

	client, err := ags.NewClient(cred, "ap-guangzhou", clientProfile)
	if err != nil {
		t.Fatalf("create ags client: %v", err)
	}
	return client
}

// newSandbox creates a sandbox using Create, skips if credentials are missing.
func newSandbox(t *testing.T) *code.Sandbox {
	t.Helper()

	client := newAgsClient(t)
	if client == nil {
		t.Skip("missing cloud credentials, skip integration tests")
	}

	sb, err := code.Create(context.TODO(), "code-interpreter-v1", code.WithClient(client), code.WithSandboxTimeout(300*time.Second))
	if err != nil {
		t.Fatalf("create sandbox: %v", err)
	}

	t.Cleanup(func() {
		_ = sb.Kill(context.TODO())
	})
	return sb
}

// newSandboxWithConnect creates a sandbox using Create, then reconnects using Connect.
// This tests the Connect flow with a valid sandbox ID.
func newSandboxWithConnect(t *testing.T) *code.Sandbox {
	t.Helper()

	client := newAgsClient(t)
	if client == nil {
		t.Skip("missing cloud credentials, skip integration tests")
	}

	// First create a sandbox to get a valid sandbox ID
	sb, err := code.Create(context.TODO(), "code-interpreter-v1", code.WithClient(client), code.WithSandboxTimeout(300*time.Second))
	if err != nil {
		t.Fatalf("create sandbox: %v", err)
	}

	sandboxId := sb.SandboxId

	// Now connect to the same sandbox using Connect
	sbConnected, err := code.Connect(context.TODO(), sandboxId, code.WithClient(client))
	if err != nil {
		// Cleanup the created sandbox on error
		_ = sb.Kill(context.TODO())
		t.Fatalf("connect to sandbox: %v", err)
	}

	t.Cleanup(func() {
		_ = sbConnected.Kill(context.TODO())
	})
	return sbConnected
}

// fmtInt64 formats int64 to string without importing fmt.
func fmtInt64(v int64) string {
	return strconv.FormatInt(v, 10)
}

// runWithBothModes runs the test function with both Create and Connect modes.
// Each mode runs as a subtest.
func runWithBothModes(t *testing.T, testFn func(t *testing.T, sb *code.Sandbox)) {
	t.Helper()

	t.Run("Create", func(t *testing.T) {
		sb := newSandbox(t)
		testFn(t, sb)
	})

	t.Run("Connect", func(t *testing.T) {
		sb := newSandboxWithConnect(t)
		testFn(t, sb)
	})
}
