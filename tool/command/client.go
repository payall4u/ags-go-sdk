package command

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/TencentCloudAgentRuntime/ags-go-sdk/connection"
	"github.com/TencentCloudAgentRuntime/ags-go-sdk/pb/process"
	"github.com/TencentCloudAgentRuntime/ags-go-sdk/pb/process/processconnect"
)

const (
	// keepalivePingIntervalSeconds is the interval in seconds for sending keepalive pings
	// to maintain the connection with the remote process server.
	keepalivePingIntervalSeconds = "30"
)

// Client is a lightweight wrapper for process RPC operations
type Client struct {
	config     *connection.Config
	httpClient *http.Client
	rpcClient  processconnect.ProcessClient
}

// New creates a new command client with the given configuration
func New(config *connection.Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	return &Client{
		config:     config,
		httpClient: newHttpClient(config),
		rpcClient:  newRpcClient(config),
	}, nil
}

// Run executes a command in foreground and aggregates results (non-PTY)
func (client *Client) Run(ctx context.Context, cmd string, config *ProcessConfig, onOutput *OnOutputConfig) (*Result, error) {
	if client == nil {
		return nil, fmt.Errorf("command client not initialized")
	}
	config = config.loadDefault()
	onOutput = onOutput.loadDefault()
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	handle, err := client.Start(ctx, cmd, config, &OnOutputConfig{
		OnStdout: func(bytes []byte) {
			stdout.Write(bytes)
			onOutput.OnStdout(bytes)
		},
		OnStderr: func(bytes []byte) {
			stderr.Write(bytes)
			onOutput.OnStderr(bytes)
		},
	})
	if err != nil {
		return nil, err
	}
	result, err := handle.Wait(ctx)
	if err != nil {
		return nil, err
	}
	return &Result{
		ExitCode: result.ExitCode,
		Stdout:   stdout.Bytes(),
		Stderr:   stderr.Bytes(),
		Error:    result.Error,
	}, nil
}

// Start starts a new process and returns a handle to it
func (client *Client) Start(ctx context.Context, cmd string, config *ProcessConfig, onOutput *OnOutputConfig) (*Handle, error) {
	if client == nil {
		return nil, fmt.Errorf("command client not initialized")
	}
	config = config.loadDefault()
	onOutput = onOutput.loadDefault()
	processConfig := buildProcessConfig(cmd, config)
	req := newRequestWithHeaders(&process.StartRequest{Process: processConfig}, client.config)
	req.Header().Set(
		"Authorization",
		fmt.Sprintf(
			"Basic %v",
			base64.StdEncoding.EncodeToString(
				[]byte(config.User+":"),
			),
		),
	)
	req.Header().Set("Keepalive-Ping-Interval", keepalivePingIntervalSeconds)
	cancelCtx, cancel := context.WithCancel(ctx)
	stream, err := client.rpcClient.Start(cancelCtx, req)
	if err != nil {
		cancel()
		return nil, err
	}
	handle := &Handle{
		cancel:   cancel,
		client:   client,
		stream:   &startStream{stream: stream},
		onStdout: onOutput.OnStdout,
		onStderr: onOutput.OnStderr,
	}
	err = handle.processEvent()
	if err != nil {
		return nil, err
	}
	return handle, nil
}

// Connect connects to an existing process by PID
func (client *Client) Connect(ctx context.Context, pid uint32, onOutput *OnOutputConfig) (*Handle, error) {
	if client == nil {
		return nil, fmt.Errorf("command client not initialized")
	}
	onOutput = onOutput.loadDefault()
	req := newRequestWithHeaders(&process.ConnectRequest{
		Process: &process.ProcessSelector{
			Selector: &process.ProcessSelector_Pid{Pid: pid},
		},
	}, client.config)
	cancelCtx, cancel := context.WithCancel(ctx)
	stream, err := client.rpcClient.Connect(cancelCtx, req)
	if err != nil {
		cancel()
		return nil, err
	}
	handle := &Handle{
		cancel:   cancel,
		client:   client,
		stream:   &connectStream{stream: stream},
		onStdout: onOutput.OnStdout,
		onStderr: onOutput.OnStderr,
	}
	err = handle.processEvent()
	if err != nil {
		return nil, err
	}
	return handle, nil
}

// List returns information about currently running processes
func (client *Client) List(ctx context.Context) ([]ProcessInfo, error) {
	if client == nil || client.config == nil {
		return nil, fmt.Errorf("command client not initialized")
	}
	req := newRequestWithHeaders(&process.ListRequest{}, client.config)
	resp, err := client.rpcClient.List(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.Msg == nil {
		return nil, fmt.Errorf("empty response from List")
	}
	processes := resp.Msg.GetProcesses()
	out := make([]ProcessInfo, 0, len(processes))
	for _, p := range processes {
		if m := mapProcessInfo(p); m != nil {
			out = append(out, *m)
		}
	}
	return out, nil
}
