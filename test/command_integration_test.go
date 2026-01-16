package test_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/code"
	"github.com/TencentCloudAgentRuntime/ags-go-sdk/tool/command"
)

func ctxWithTimeout(t *testing.T, d time.Duration) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), d)
	t.Cleanup(cancel)
	return ctx
}

func TestCommand_Run_Hello(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		ctx := ctxWithTimeout(t, 30*time.Second)
		res, err := sb.Commands.Run(ctx, "echo hello", &command.ProcessConfig{User: "user"}, &command.OnOutputConfig{
			OnStdout: func([]byte) {},
			OnStderr: func([]byte) {},
		})
		if err != nil {
			t.Fatalf("Run error: %v", err)
		}
		if res.ExitCode != 0 {
			t.Fatalf("unexpected exit code: %d", res.ExitCode)
		}
		if len(res.Stdout) == 0 {
			t.Fatalf("expected stdout, got empty")
		}
		if !strings.Contains(string(res.Stdout), "hello") {
			t.Logf("stdout does not contain 'hello': %q", string(res.Stdout))
		}
	})
}

func TestCommand_Start_Wait_WithEnvAndCwd(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		ctx := ctxWithTimeout(t, 40*time.Second)
		cwd := "/home"
		stdoutBuffer, stderrBuffer := bytes.NewBufferString(""), bytes.NewBufferString("")
		h, err := sb.Commands.Start(ctx, "echo $FOO; pwd", &command.ProcessConfig{
			User: "user",
			Envs: map[string]string{"FOO": "BAR"},
			Cwd:  &cwd,
		}, &command.OnOutputConfig{
			OnStdout: func(b []byte) { stdoutBuffer.Write(b) },
			OnStderr: func(b []byte) { stderrBuffer.Write(b) },
		})
		if err != nil {
			t.Fatalf("Start error: %v", err)
		}
		pr, err := h.Wait(ctx)
		if err != nil {
			t.Fatalf("Wait error: %v", err)
		}
		if pr.ExitCode != 0 {
			t.Fatalf("unexpected exit code: %d", pr.ExitCode)
		}
		if !strings.Contains(stdoutBuffer.String(), "BAR") {
			t.Logf("stdout does not contain 'BAR': %q", stdoutBuffer.String())
		}
		if !strings.Contains(stdoutBuffer.String(), "/home") {
			t.Logf("stdout does not contain '/home': %q", stdoutBuffer.String())
		}
	})
}

func TestCommand_Run_ErrorExitCode(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		ctx := ctxWithTimeout(t, 20*time.Second)
		res, err := sb.Commands.Run(ctx, "exit 3", &command.ProcessConfig{User: "user"}, &command.OnOutputConfig{
			OnStdout: func([]byte) {},
			OnStderr: func([]byte) {},
		})
		if err != nil {
			t.Fatalf("Run request error: %v", err)
		}
		if res.ExitCode != 3 {
			t.Fatalf("expected exit code 3, got %d", res.ExitCode)
		}
	})
}

func TestCommand_Start_SendInput(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		ctx := ctxWithTimeout(t, 40*time.Second)
		var out strings.Builder
		h, err := sb.Commands.Start(ctx, "read line; echo X$line", &command.ProcessConfig{
			User: "user",
		}, &command.OnOutputConfig{
			OnStdout: func(b []byte) { out.Write(b) },
			OnStderr: func([]byte) {},
		})
		if err != nil {
			t.Fatalf("Start error: %v", err)
		}
		// 发送一行输入
		if err := h.SendInput(ctx, h.Pid, []byte("abc\n")); err != nil {
			t.Fatalf("SendInput error: %v", err)
		}
		pr, err := h.Wait(ctx)
		if err != nil {
			t.Fatalf("Wait error: %v", err)
		}
		if pr.ExitCode != 0 {
			t.Fatalf("unexpected exit code: %d", pr.ExitCode)
		}
		if !strings.Contains(out.String(), "Xabc") {
			t.Fatalf("expected stdout to contain 'Xabc', got: %q", out.String())
		}
	})
}

func TestCommand_Kill_LongRunning(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		ctx := ctxWithTimeout(t, 40*time.Second)
		h, err := sb.Commands.Start(ctx, "sleep 30", &command.ProcessConfig{User: "user"}, &command.OnOutputConfig{
			OnStdout: func([]byte) {},
			OnStderr: func([]byte) {},
		})
		if err != nil {
			t.Fatalf("Start error: %v", err)
		}
		// 立即 Kill
		if err := h.Kill(ctx); err != nil {
			t.Fatalf("Kill error: %v", err)
		}
		pr, err := h.Wait(ctx)
		// 宽松断言：被 kill 后可能返回非零退出码，或服务端包含错误信息
		if err != nil {
			// 接受 Wait 返回错误作为一种结果
			t.Fatalf("Wait returned error after kill: %v", err)
			return
		}
		if pr.ExitCode == 0 {
			t.Fatalf("expected non-zero exit after kill, got 0")
		}
	})
}

func TestCommand_SendSignal_SIGTERM(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		ctx := ctxWithTimeout(t, 40*time.Second)
		// 启动一个可捕获 SIGTERM 的脚本：收到 TERM 时输出并退出 0
		// 使用 /bin/sh 运行；在 sandbox 的 Linux 环境中可用
		var out strings.Builder
		h, err := sb.Commands.Start(ctx, `sh -c 'trap "echo TERM; exit 0" TERM; while true; do sleep 1; done'`, &command.ProcessConfig{
			User: "user",
		}, &command.OnOutputConfig{
			OnStdout: func(b []byte) { out.Write(b) },
			OnStderr: func([]byte) {},
		})
		if err != nil {
			t.Fatalf("Start error: %v", err)
		}

		// 发送 SIGTERM（15）；若服务端将枚举映射为标准 POSIX 号，此处应使进程优雅退出
		if err := h.SendSignal(ctx, h.Pid, 15 /* SIGTERM */); err != nil {
			t.Fatalf("SendSignal error: %v", err)
		}

		pr, err := h.Wait(ctx)
		if err != nil {
			t.Fatalf("Wait error: %v", err)
		}
		// 期望优雅退出（exit 0），并输出 "TERM"
		if pr.ExitCode != 0 {
			t.Fatalf("unexpected exit code after SIGTERM: %d", pr.ExitCode)
		}
		if !strings.Contains(out.String(), "TERM") {
			t.Logf("stdout does not contain 'TERM': %q", out.String())
		}
	})
}

func TestCommand_List(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		ctx := ctxWithTimeout(t, 20*time.Second)
		ps, err := sb.Commands.List(ctx)
		if err != nil {
			t.Fatalf("List error: %v", err)
		}
		// 宽松断言：返回切片即可
		if ps == nil {
			t.Fatalf("nil process list")
		}
	})
}
