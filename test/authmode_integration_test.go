package test_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/TencentCloudAgentRuntime/ags-go-sdk/constant"
	"github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/code"
	"github.com/TencentCloudAgentRuntime/ags-go-sdk/tool/command"
)

// newNoAuthSandbox creates a sandbox with AuthMode=NONE.
func newNoAuthSandbox(t *testing.T) *code.Sandbox {
	t.Helper()

	client := newAgsClient(t)
	if client == nil {
		t.Skip("missing cloud credentials, skip integration tests")
	}

	authMode := "NONE"
	sb, err := code.Create(context.TODO(), "code-interpreter-v1",
		code.WithClient(client),
		code.WithSandboxTimeout(300*time.Second),
		code.WithSandboxConfig(&code.SandboxConfig{
			AuthMode: &authMode,
		}),
	)
	if err != nil {
		t.Fatalf("create sandbox with AuthMode=NONE: %v", err)
	}

	t.Cleanup(func() {
		_ = sb.Kill(context.TODO())
	})
	return sb
}

// TestAuthModeNone_SDK_RunCode verifies that the sandbox created with AuthMode=NONE
// works normally through the SDK (code execution).
func TestAuthModeNone_SDK_RunCode(t *testing.T) {
	sb := newNoAuthSandbox(t)

	exec, err := sb.Code.RunCode(context.TODO(), "print('auth-none-ok')", nil, nil)
	if err != nil {
		t.Fatalf("RunCode error: %v", err)
	}
	if exec.Error != nil {
		t.Fatalf("unexpected execution error: %+v", exec.Error)
	}
	found := false
	for _, line := range exec.Logs.Stdout {
		if strings.Contains(line, "auth-none-ok") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected stdout to contain 'auth-none-ok', got: %v", exec.Logs.Stdout)
	}
}

// TestAuthModeNone_SDK_RunCommand verifies that the sandbox created with AuthMode=NONE
// works normally through the SDK (command execution).
func TestAuthModeNone_SDK_RunCommand(t *testing.T) {
	sb := newNoAuthSandbox(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := sb.Commands.Run(ctx, "echo auth-none-cmd", &command.ProcessConfig{User: "user"}, nil)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if res.ExitCode != 0 {
		t.Fatalf("unexpected exit code: %d", res.ExitCode)
	}
	if !strings.Contains(string(res.Stdout), "auth-none-cmd") {
		t.Fatalf("expected stdout to contain 'auth-none-cmd', got: %q", string(res.Stdout))
	}
}

// TestAuthModeNone_RawHTTP_NoToken verifies that the data plane of a sandbox created
// with AuthMode=NONE can be accessed via raw HTTP requests WITHOUT the X-Access-Token header.
// This proves that the AuthMode=NONE truly disables token authentication on the data plane.
func TestAuthModeNone_RawHTTP_NoToken(t *testing.T) {
	sb := newNoAuthSandbox(t)

	// Build the code execution URL: https://{CodePort}-{sandboxId}.{region}.tencentags.com/execute
	codeHost := sb.GetHost(constant.CodePort)
	executeURL := fmt.Sprintf("https://%s/execute", codeHost)

	// Prepare the request body (same format as the SDK uses)
	lang := "python"
	body := map[string]any{
		"code":     "print('raw-no-token')",
		"language": &lang,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}

	// Make a raw HTTP request WITHOUT X-Access-Token header
	req, err := http.NewRequest(http.MethodPost, executeURL, bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// Intentionally NOT setting X-Access-Token header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("expected 2xx status, got %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse NDJSON response lines — look for stdout containing our marker
	found := false
	for _, line := range bytes.Split(respBody, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		var raw map[string]any
		if err := json.Unmarshal(line, &raw); err != nil {
			continue
		}
		if raw["type"] == "stdout" {
			if txt, ok := raw["text"].(string); ok && strings.Contains(txt, "raw-no-token") {
				found = true
				break
			}
		}
	}
	if !found {
		t.Fatalf("expected response to contain 'raw-no-token', got: %s", string(respBody))
	}
}
