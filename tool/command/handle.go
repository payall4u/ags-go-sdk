package command

import (
	"context"
	"fmt"

	"github.com/TencentCloudAgentRuntime/ags-go-sdk/pb/process"
)

type exitResult struct {
	result *ProcessResult
	err    error
}

// Handle provides control over a running process
type Handle struct {
	Pid      uint32
	exitChan chan exitResult
	doneChan chan struct{} // signals when the goroutine has completely finished
	cancel   context.CancelFunc
	client   *Client
	stream   eventStream
	onStdout func([]byte)
	onStderr func([]byte)
}

// Wait waits for the process to complete and returns the result
func (handle *Handle) Wait(ctx context.Context) (*ProcessResult, error) {
	if handle == nil {
		return nil, fmt.Errorf("handle is nil")
	}
	if handle.exitChan == nil {
		return nil, fmt.Errorf("handle not initialized")
	}
	if handle.doneChan == nil {
		return nil, fmt.Errorf("handle not initialized")
	}
	// Make sure every path waits for goroutine to stop
	defer func() {
		handle.cancel()
		<-handle.doneChan
	}()
	select {
	case ret := <-handle.exitChan:
		return ret.result, ret.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Kill terminates the process with SIGKILL
func (handle *Handle) Kill(ctx context.Context) error {
	if handle == nil {
		return fmt.Errorf("handle is nil")
	}
	if handle.client == nil {
		return fmt.Errorf("client is nil")
	}
	cli := handle.client.rpcClient
	req := newRequestWithHeaders(&process.SendSignalRequest{
		Process: &process.ProcessSelector{
			Selector: &process.ProcessSelector_Pid{Pid: handle.Pid},
		},
		Signal: process.Signal_SIGNAL_SIGKILL,
	}, handle.client.config)
	_, err := cli.SendSignal(ctx, req)
	return err
}

// Disconnect disconnects from the process while keeping it running
func (handle *Handle) Disconnect(_ context.Context) error {
	if handle == nil {
		return fmt.Errorf("handle is nil")
	}
	if handle.cancel != nil {
		handle.cancel()
	}
	return nil
}

// SendInput sends input to the process (stdin or PTY)
func (handle *Handle) SendInput(ctx context.Context, pid uint32, stdin []byte) error {
	if handle == nil {
		return fmt.Errorf("handle is nil")
	}
	if handle.client == nil {
		return fmt.Errorf("client is nil")
	}
	cli := handle.client.rpcClient
	req := newRequestWithHeaders(&process.SendInputRequest{
		Process: &process.ProcessSelector{
			Selector: &process.ProcessSelector_Pid{Pid: pid},
		},
		Input: &process.ProcessInput{
			Input: &process.ProcessInput_Stdin{Stdin: stdin},
		},
	}, handle.client.config)
	_, err := cli.SendInput(ctx, req)
	return err
}

// SendSignal sends a signal to the process (e.g., SIGTERM, SIGKILL)
func (handle *Handle) SendSignal(ctx context.Context, pid uint32, sig process.Signal) error {
	if handle == nil {
		return fmt.Errorf("handle is nil")
	}
	if handle.client == nil {
		return fmt.Errorf("client is nil")
	}
	if sig != 9 && sig != 15 {
		return fmt.Errorf("invalid signal: %d, only SIGTERM(15) and SIGKILL(9) supported", sig)
	}
	cli := handle.client.rpcClient
	req := newRequestWithHeaders(&process.SendSignalRequest{
		Process: &process.ProcessSelector{
			Selector: &process.ProcessSelector_Pid{Pid: pid},
		},
		Signal: sig,
	}, handle.client.config)
	_, err := cli.SendSignal(ctx, req)
	return err
}

// processEvent starts a background loop to consume events, trigger callbacks, and handle process completion
func (handle *Handle) processEvent() error {
	if handle == nil {
		return fmt.Errorf("handle is nil")
	}
	if handle.exitChan == nil {
		handle.exitChan = make(chan exitResult, 1)
	}
	if handle.doneChan == nil {
		handle.doneChan = make(chan struct{})
	}
	ev, err := handle.stream.next()
	if err != nil {
		handle.cancel()
		return err
	}
	if se := ev.GetStart(); se != nil {
		handle.Pid = se.GetPid()
	}
	go func() {
		defer close(handle.doneChan) // Signal completion when goroutine exits
		defer handle.cancel()        // Ensure cancel is always called
		for {
			// If handle.cancel is called, the context will be canceled,
			// which will cause the stream to be closed.
			// handle.stream.next() will return error when the stream is closed.
			ev, err := handle.stream.next()
			if err != nil {
				handle.exitChan <- exitResult{err: err}
				return
			}
			switch e := ev.Event.(type) {
			case *process.ProcessEvent_Data:
				if d := e.Data; d != nil {
					if out := d.GetStdout(); len(out) > 0 {
						if handle.onStdout != nil {
							handle.onStdout(out)
						}
					}
					if stderr := d.GetStderr(); len(stderr) > 0 {
						if handle.onStderr != nil {
							handle.onStderr(stderr)
						}
					}
				}
			case *process.ProcessEvent_End:
				handle.exitChan <- exitResult{result: &ProcessResult{
					ExitCode: e.End.GetExitCode(),
					Error:    e.End.Error,
				}}
				return
			default:
			}
		}
	}()
	return nil
}
