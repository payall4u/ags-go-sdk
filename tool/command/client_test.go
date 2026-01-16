package command

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"connectrpc.com/connect"
	"github.com/TencentCloudAgentRuntime/ags-go-sdk/connection"
	"github.com/TencentCloudAgentRuntime/ags-go-sdk/pb/process"
	"github.com/TencentCloudAgentRuntime/ags-go-sdk/pb/process/processconnect"
)

type testProcessHandler struct {
	processconnect.UnimplementedProcessHandler

	listResponse  *process.ListResponse
	startEvents   []*process.ProcessEvent
	connectEvents []*process.ProcessEvent
}

func (h *testProcessHandler) List(_ context.Context, _ *connect.Request[process.ListRequest]) (*connect.Response[process.ListResponse], error) {
	if h.listResponse == nil {
		return connect.NewResponse(&process.ListResponse{}), nil
	}
	return connect.NewResponse(h.listResponse), nil
}

func (h *testProcessHandler) Start(_ context.Context, _ *connect.Request[process.StartRequest], stream *connect.ServerStream[process.StartResponse]) error {
	for _, ev := range h.startEvents {
		if err := stream.Send(&process.StartResponse{Event: ev}); err != nil {
			return err
		}
	}
	return nil
}

func (h *testProcessHandler) Connect(_ context.Context, _ *connect.Request[process.ConnectRequest], stream *connect.ServerStream[process.ConnectResponse]) error {
	for _, ev := range h.connectEvents {
		if err := stream.Send(&process.ConnectResponse{Event: ev}); err != nil {
			return err
		}
	}
	return nil
}

func newTestClient(t *testing.T, handler processconnect.ProcessHandler) *Client {
	t.Helper()

	path, httpHandler := processconnect.NewProcessHandler(handler)
	mux := http.NewServeMux()
	mux.Handle(path, httpHandler)

	server := httptest.NewTLSServer(mux)
	t.Cleanup(server.Close)

	domain := strings.TrimPrefix(server.URL, "https://")
	cfg := &connection.Config{
		Domain:      domain,
		AccessToken: "test-token",
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	client.httpClient = server.Client()
	client.rpcClient = processconnect.NewProcessClient(server.Client(), server.URL, connect.WithProtoJSON())

	return client
}

func buildTestEvents() ([]*process.ProcessEvent, int32, string) {
	pid := uint32(123)
	stdout := []byte("hello")
	stderr := []byte("world")
	exitCode := int32(42)
	exitErr := "exit error"

	startEv := &process.ProcessEvent{
		Event: &process.ProcessEvent_Start{
			Start: &process.ProcessEvent_StartEvent{Pid: pid},
		},
	}

	stdoutEv := &process.ProcessEvent{
		Event: &process.ProcessEvent_Data{
			Data: &process.ProcessEvent_DataEvent{
				Output: &process.ProcessEvent_DataEvent_Stdout{Stdout: stdout},
			},
		},
	}

	stderrEv := &process.ProcessEvent{
		Event: &process.ProcessEvent_Data{
			Data: &process.ProcessEvent_DataEvent{
				Output: &process.ProcessEvent_DataEvent_Stderr{Stderr: stderr},
			},
		},
	}

	endEv := &process.ProcessEvent{
		Event: &process.ProcessEvent_End{
			End: &process.ProcessEvent_EndEvent{
				ExitCode: exitCode,
				Error:    &exitErr,
				Exited:   true,
				Status:   "exited",
			},
		},
	}

	return []*process.ProcessEvent{startEv, stdoutEv, stderrEv, endEv}, exitCode, exitErr
}

func TestNew_NilConfig(t *testing.T) {
	client, err := New(nil)
	if err == nil {
		t.Fatalf("expected error when config is nil, got client = %#v", client)
	}
}

func TestNew_Success(t *testing.T) {
	cfg := &connection.Config{Domain: "example.com"}
	client, err := New(cfg)
	if err != nil {
		t.Fatalf("unexpected error from New: %v", err)
	}
	if client.config != cfg {
		t.Fatalf("client.config not set correctly")
	}
	if client.httpClient == nil {
		t.Errorf("httpClient should not be nil")
	}
	if client.rpcClient == nil {
		t.Errorf("rpcClient should not be nil")
	}
}

func TestClient_Run_AggregatesOutputAndExitCode(t *testing.T) {
	events, exitCode, exitErr := buildTestEvents()
	handler := &testProcessHandler{startEvents: events}
	client := newTestClient(t, handler)

	var stdoutCalls [][]byte
	var stderrCalls [][]byte
	onOutput := &OnOutputConfig{
		OnStdout: func(p []byte) {
			cp := make([]byte, len(p))
			copy(cp, p)
			stdoutCalls = append(stdoutCalls, cp)
		},
		OnStderr: func(p []byte) {
			cp := make([]byte, len(p))
			copy(cp, p)
			stderrCalls = append(stderrCalls, cp)
		},
	}

	res, err := client.Run(context.Background(), "echo test", &ProcessConfig{User: "user"}, onOutput)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if res.ExitCode != exitCode {
		t.Errorf("unexpected exit code: got %d, want %d", res.ExitCode, exitCode)
	}
	if string(res.Stdout) != "hello" {
		t.Errorf("unexpected stdout: got %q, want %q", string(res.Stdout), "hello")
	}
	if string(res.Stderr) != "world" {
		t.Errorf("unexpected stderr: got %q, want %q", string(res.Stderr), "world")
	}
	if res.Error == nil || *res.Error != exitErr {
		if res.Error == nil {
			t.Errorf("expected error %q, got nil", exitErr)
		} else {
			t.Errorf("unexpected error message: got %q, want %q", *res.Error, exitErr)
		}
	}

	if len(stdoutCalls) != 1 || string(stdoutCalls[0]) != "hello" {
		t.Errorf("OnStdout not invoked as expected: %#v", stdoutCalls)
	}
	if len(stderrCalls) != 1 || string(stderrCalls[0]) != "world" {
		t.Errorf("OnStderr not invoked as expected: %#v", stderrCalls)
	}
}

// Test combinations of ProcessConfig (nil/non-nil) and OnOutputConfig (nil/non-nil) for Run.
func TestClient_Run_ConfigAndOnOutputVariants(t *testing.T) {
	events, _, _ := buildTestEvents()
	handler := &testProcessHandler{startEvents: events}
	client := newTestClient(t, handler)

	tests := []struct {
		name        string
		cfg         *ProcessConfig
		onOutputCfg *OnOutputConfig
		expectCalls bool
	}{
		{
			name:        "ConfigNil_OnOutputNil",
			cfg:         nil,
			onOutputCfg: nil,
			expectCalls: false,
		},
		{
			name:        "ConfigNonNil_OnOutputNil",
			cfg:         &ProcessConfig{User: "user"},
			onOutputCfg: nil,
			expectCalls: false,
		},
		{
			name: "ConfigNil_OnOutputNonNil",
			cfg:  nil,
			onOutputCfg: &OnOutputConfig{
				OnStdout: func([]byte) {},
				OnStderr: func([]byte) {},
			},
			expectCalls: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdoutCalls [][]byte
			var stderrCalls [][]byte

			onOutput := tt.onOutputCfg
			if onOutput != nil {
				onOutput = &OnOutputConfig{
					OnStdout: func(p []byte) {
						cp := make([]byte, len(p))
						copy(cp, p)
						stdoutCalls = append(stdoutCalls, cp)
					},
					OnStderr: func(p []byte) {
						cp := make([]byte, len(p))
						copy(cp, p)
						stderrCalls = append(stderrCalls, cp)
					},
				}
			}

			res, err := client.Run(context.Background(), "echo test", tt.cfg, onOutput)
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			if res.ExitCode == 0 && len(res.Stdout) == 0 && len(res.Stderr) == 0 {
				t.Fatalf("expected non-empty result from Run")
			}

			if tt.expectCalls {
				if len(stdoutCalls) == 0 {
					t.Errorf("expected OnStdout to be called")
				}
			}
		})
	}
}

// Test Start with config/onOutput being nil and non-nil.
func TestClient_Start_ConfigAndOnOutputVariants(t *testing.T) {
	events, exitCode, _ := buildTestEvents()
	handler := &testProcessHandler{startEvents: events}
	client := newTestClient(t, handler)

	tests := []struct {
		name        string
		cfg         *ProcessConfig
		onOutputCfg *OnOutputConfig
		expectCalls bool
	}{
		{
			name:        "ConfigNil_OnOutputNil",
			cfg:         nil,
			onOutputCfg: nil,
			expectCalls: false,
		},
		{
			name:        "ConfigNonNil_OnOutputNil",
			cfg:         &ProcessConfig{User: "user"},
			onOutputCfg: nil,
			expectCalls: false,
		},
		{
			name: "ConfigNonNil_OnOutputNonNil",
			cfg:  &ProcessConfig{User: "user"},
			onOutputCfg: &OnOutputConfig{
				OnStdout: func([]byte) {},
				OnStderr: func([]byte) {},
			},
			expectCalls: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdoutCalls [][]byte
			var stderrCalls [][]byte

			onOutput := tt.onOutputCfg
			if onOutput != nil {
				onOutput = &OnOutputConfig{
					OnStdout: func(p []byte) {
						cp := make([]byte, len(p))
						copy(cp, p)
						stdoutCalls = append(stdoutCalls, cp)
					},
					OnStderr: func(p []byte) {
						cp := make([]byte, len(p))
						copy(cp, p)
						stderrCalls = append(stderrCalls, cp)
					},
				}
			}

			handle, err := client.Start(context.Background(), "echo test", tt.cfg, onOutput)
			if err != nil {
				t.Fatalf("Start() error = %v", err)
			}

			res, err := handle.Wait(context.Background())
			if err != nil {
				t.Fatalf("Wait() error = %v", err)
			}
			if res.ExitCode != exitCode {
				t.Errorf("unexpected exit code: got %d, want %d", res.ExitCode, exitCode)
			}

			if tt.expectCalls && len(stdoutCalls) == 0 {
				t.Errorf("expected OnStdout to be called")
			}
		})
	}
}

// Test Connect with onOutput nil and non-nil.
func TestClient_Connect_OnOutputVariants(t *testing.T) {
	events, exitCode, _ := buildTestEvents()
	handler := &testProcessHandler{connectEvents: events}
	client := newTestClient(t, handler)

	t.Run("OnOutputNil", func(t *testing.T) {
		handle, err := client.Connect(context.Background(), 999, nil)
		if err != nil {
			t.Fatalf("Connect() error = %v", err)
		}
		res, err := handle.Wait(context.Background())
		if err != nil {
			t.Fatalf("Wait() error = %v", err)
		}
		if res.ExitCode != exitCode {
			t.Errorf("unexpected exit code: got %d, want %d", res.ExitCode, exitCode)
		}
	})

	t.Run("OnOutputNonNil", func(t *testing.T) {
		var stdoutCalls [][]byte
		onOutput := &OnOutputConfig{
			OnStdout: func(p []byte) {
				cp := make([]byte, len(p))
				copy(cp, p)
				stdoutCalls = append(stdoutCalls, cp)
			},
		}

		handle, err := client.Connect(context.Background(), 999, onOutput)
		if err != nil {
			t.Fatalf("Connect() error = %v", err)
		}
		res, err := handle.Wait(context.Background())
		if err != nil {
			t.Fatalf("Wait() error = %v", err)
		}
		if res.ExitCode != exitCode {
			t.Errorf("unexpected exit code: got %d, want %d", res.ExitCode, exitCode)
		}
		if len(stdoutCalls) == 0 {
			t.Errorf("expected OnStdout to be called")
		}
	})
}

func TestClient_Connect_StartsAndCollectsResult(t *testing.T) {
	events, exitCode, _ := buildTestEvents()
	handler := &testProcessHandler{connectEvents: events}
	client := newTestClient(t, handler)

	var stdoutCalls [][]byte
	onOutput := &OnOutputConfig{
		OnStdout: func(p []byte) {
			cp := make([]byte, len(p))
			copy(cp, p)
			stdoutCalls = append(stdoutCalls, cp)
		},
	}

	handle, err := client.Connect(context.Background(), 999, onOutput)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	res, err := handle.Wait(context.Background())
	if err != nil {
		t.Fatalf("Wait() error = %v", err)
	}
	if res.ExitCode != exitCode {
		t.Errorf("unexpected exit code from handle: got %d, want %d", res.ExitCode, exitCode)
	}
	if len(stdoutCalls) == 0 || string(stdoutCalls[0]) != "hello" {
		t.Errorf("expected OnStdout to receive 'hello', got %#v", stdoutCalls)
	}
}

func TestClient_List_MapsProcessInfo(t *testing.T) {
	tag := "tag1"
	cwd := "/tmp"
	handler := &testProcessHandler{
		listResponse: &process.ListResponse{
			Processes: []*process.ProcessInfo{
				{
					Config: &process.ProcessConfig{
						Cmd:  "bash",
						Args: []string{"-c", "echo"},
						Envs: map[string]string{"KEY": "VALUE"},
						Cwd:  &cwd,
					},
					Pid: 123,
					Tag: &tag,
				},
			},
		},
	}

	client := newTestClient(t, handler)
	got, err := client.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 process, got %d", len(got))
	}

	p := got[0]
	if p.Pid != 123 {
		t.Errorf("unexpected pid: got %d, want %d", p.Pid, 123)
	}
	if p.Tag == nil || *p.Tag != tag {
		t.Errorf("unexpected tag: got %#v, want %q", p.Tag, tag)
	}
	if p.Cmd != "bash" {
		t.Errorf("unexpected cmd: got %q, want %q", p.Cmd, "bash")
	}
	if len(p.Args) != 2 || p.Args[0] != "-c" || p.Args[1] != "echo" {
		t.Errorf("unexpected args: %#v", p.Args)
	}
	if p.Envs["KEY"] != "VALUE" {
		t.Errorf("unexpected envs: %#v", p.Envs)
	}
	if p.Cwd == nil || *p.Cwd != cwd {
		t.Errorf("unexpected cwd: got %#v, want %q", p.Cwd, cwd)
	}
}

func TestClient_ErrorCases(t *testing.T) {
	var nilClient *Client

	if _, err := nilClient.Run(context.Background(), "cmd", nil, nil); err == nil {
		t.Errorf("expected error from Run on nil client, got nil")
	}

	if _, err := nilClient.Start(context.Background(), "cmd", nil, nil); err == nil {
		t.Errorf("expected error from Start on nil client, got nil")
	}

	if _, err := nilClient.Connect(context.Background(), 1, nil); err == nil {
		t.Errorf("expected error from Connect on nil client, got nil")
	}

	if _, err := nilClient.List(context.Background()); err == nil {
		t.Errorf("expected error from List on nil client, got nil")
	}

	client := &Client{}
	if _, err := client.List(context.Background()); err == nil {
		t.Errorf("expected error from List with nil config, got nil")
	}
}
