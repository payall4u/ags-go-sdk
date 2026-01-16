package code

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TencentCloudAgentRuntime/ags-go-sdk/connection"
)

func TestRunCode_MaxBufferSize(t *testing.T) {
	// Generate a long string (70KB, larger than default 64KB initial buffer)
	longData := strings.Repeat("a", 70*1024)
	respData := map[string]interface{}{
		"type": "stdout",
		"text": longData,
	}
	jsonBytes, err := json.Marshal(respData)
	if err != nil {
		panic(err)
	}
	jsonLine := string(jsonBytes) + "\n"

	// Create a TLS server to match the client's HTTPS requirement
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, jsonLine)
	}))
	defer server.Close()

	// Extract host from server URL (remove https:// prefix)
	serverURL := server.URL
	domain := strings.TrimPrefix(serverURL, "https://")

	cfg := &connection.Config{
		Domain: domain,
	}
	client := New(cfg)
	// Inject the server's client which trusts the test server's certificate
	client.httpClient = server.Client()

	t.Run("SufficientBuffer", func(t *testing.T) {
		// Set buffer to 100KB (larger than the ~70KB line)
		runConfig := &RunCodeConfig{
			MaxBufferSize: 100 * 1024,
		}
		exec, err := client.RunCode(context.Background(), "print('hello')", runConfig, nil)
		if err != nil {
			t.Fatalf("Expected success, got error: %v", err)
		}
		if len(exec.Logs.Stdout) != 1 || exec.Logs.Stdout[0] != longData {
			t.Errorf("Unexpected stdout content")
		}
	})

	t.Run("InsufficientBuffer", func(t *testing.T) {
		// Set buffer to 65KB (larger than initial 64KB, but smaller than 70KB data)
		// Note: Since initial buffer is 64KB, we must test with data > 64KB to verify limit
		runConfig := &RunCodeConfig{
			MaxBufferSize: 65 * 1024,
		}
		_, err := client.RunCode(context.Background(), "print('hello')", runConfig, nil)
		if err == nil {
			t.Fatal("Expected error due to small buffer, got nil")
		}
		// Expect bufio.Scanner: token too long
		if !strings.Contains(err.Error(), "token too long") {
			t.Errorf("Expected 'token too long' error, got: %v", err)
		}
	})

	t.Run("DefaultBuffer", func(t *testing.T) {
		// Default buffer is 1GB, should handle 70KB easily
		runConfig := &RunCodeConfig{}
		exec, err := client.RunCode(context.Background(), "print('hello')", runConfig, nil)
		if err != nil {
			t.Fatalf("Expected success with default buffer, got error: %v", err)
		}
		if len(exec.Logs.Stdout) != 1 || exec.Logs.Stdout[0] != longData {
			t.Errorf("Unexpected stdout content")
		}
	})
}
