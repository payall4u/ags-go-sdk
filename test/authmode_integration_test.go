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

// newPublicAuthSandbox creates a sandbox with AuthMode=PUBLIC.
// In PUBLIC mode only the ENVD management port (49983) requires a TOKEN;
// all other ports can be accessed without a TOKEN.
func newPublicAuthSandbox(t *testing.T) *code.Sandbox {
	t.Helper()

	client := newAgsClient(t)
	if client == nil {
		t.Skip("missing cloud credentials, skip integration tests")
	}

	authMode := "PUBLIC"
	sb, err := code.Create(context.TODO(), "code-interpreter-v1",
		code.WithClient(client),
		code.WithSandboxTimeout(300*time.Second),
		code.WithSandboxConfig(&code.SandboxConfig{
			AuthMode: &authMode,
		}),
	)
	if err != nil {
		t.Fatalf("create sandbox with AuthMode=PUBLIC: %v", err)
	}

	t.Cleanup(func() {
		_ = sb.Kill(context.TODO())
	})
	return sb
}

// TestAuthModePublic_SDK_RunCode verifies that a sandbox created with AuthMode=PUBLIC
// still works through the SDK (code execution on the Code port, no token required).
func TestAuthModePublic_SDK_RunCode(t *testing.T) {
	sb := newPublicAuthSandbox(t)

	exec, err := sb.Code.RunCode(context.TODO(), "print('auth-public-ok')", nil, nil)
	if err != nil {
		t.Fatalf("RunCode error: %v", err)
	}
	if exec.Error != nil {
		t.Fatalf("unexpected execution error: %+v", exec.Error)
	}
	found := false
	for _, line := range exec.Logs.Stdout {
		if strings.Contains(line, "auth-public-ok") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected stdout to contain 'auth-public-ok', got: %v", exec.Logs.Stdout)
	}
}

// TestAuthModePublic_SDK_RunCommand verifies that a sandbox created with AuthMode=PUBLIC
// still works through the SDK for command execution (ENVD port, SDK supplies the token internally).
func TestAuthModePublic_SDK_RunCommand(t *testing.T) {
	sb := newPublicAuthSandbox(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := sb.Commands.Run(ctx, "echo auth-public-cmd", &command.ProcessConfig{User: "user"}, nil)
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if res.ExitCode != 0 {
		t.Fatalf("unexpected exit code: %d", res.ExitCode)
	}
	if !strings.Contains(string(res.Stdout), "auth-public-cmd") {
		t.Fatalf("expected stdout to contain 'auth-public-cmd', got: %q", string(res.Stdout))
	}
}

// TestAuthModePublic_RawHTTP_CodePort_NoToken verifies that in PUBLIC mode the
// non-ENVD ports (e.g. Code port) can be accessed via raw HTTP WITHOUT X-Access-Token.
func TestAuthModePublic_RawHTTP_CodePort_NoToken(t *testing.T) {
	sb := newPublicAuthSandbox(t)

	codeHost := sb.GetHost(constant.CodePort)
	executeURL := fmt.Sprintf("https://%s/execute", codeHost)

	lang := "python"
	body := map[string]any{
		"code":     "print('public-code-no-token')",
		"language": &lang,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, executeURL, bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// Intentionally NOT setting X-Access-Token header.

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("expected 2xx on Code port in PUBLIC mode without token, got %d: %s",
			resp.StatusCode, string(respBody))
	}

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
			if txt, ok := raw["text"].(string); ok && strings.Contains(txt, "public-code-no-token") {
				found = true
				break
			}
		}
	}
	if !found {
		t.Fatalf("expected response to contain 'public-code-no-token', got: %s", string(respBody))
	}
}

// TestAuthModePublic_RawHTTP_EnvdPort_NoTokenRejected verifies that in PUBLIC mode
// the ENVD management port (49983) STILL requires a TOKEN: a raw HTTP request
// without X-Access-Token must be rejected (non-2xx status).
func TestAuthModePublic_RawHTTP_EnvdPort_NoTokenRejected(t *testing.T) {
	sb := newPublicAuthSandbox(t)

	envdHost := sb.GetHost(constant.EnvdPort)
	// Any ENVD endpoint will do; use a lightweight one. We only care about the
	// authentication result, not the concrete endpoint semantics.
	probeURL := fmt.Sprintf("https://%s/", envdHost)

	req, err := http.NewRequest(http.MethodGet, probeURL, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	// Intentionally NOT setting X-Access-Token header.

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected ENVD port (49983) to reject requests without token in PUBLIC mode, got %d: %s",
			resp.StatusCode, string(respBody))
	}
}
