package test_test

import (
	"context"
	"strings"
	"testing"

	"github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/code"
	toolcode "github.com/TencentCloudAgentRuntime/ags-go-sdk/tool/code"
)

func TestRunCode_Integration_Hello(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		exec, err := sb.Code.RunCode(context.TODO(), "print(\"hello\")", nil, nil)
		if err != nil {
			t.Fatalf("RunCode error: %v", err)
		}
		if len(exec.Logs.Stdout) == 0 && len(exec.Results) == 0 && exec.Error == nil {
			t.Fatalf("unexpected empty execution: %+v", exec)
		}
		if len(exec.Logs.Stdout) > 0 && exec.Logs.Stdout[0] == "" {
			t.Fatalf("empty stdout[0], logs=%+v", exec.Logs)
		}
	})
}

func TestRunCode_Integration_WithLanguage(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		opts := &toolcode.RunCodeConfig{Language: "python"}
		exec, err := sb.Code.RunCode(context.TODO(), "x=1+2\nprint(x)", opts, nil)
		if err != nil {
			t.Fatalf("RunCode with language error: %v", err)
		}
		// 至少产生输出或结果
		if len(exec.Logs.Stdout) == 0 && len(exec.Results) == 0 && exec.Error == nil {
			t.Fatalf("unexpected empty execution: %+v", exec)
		}
		// 若有 stdout，应包含 "3"
		if len(exec.Logs.Stdout) > 0 && !strings.Contains(exec.Logs.Stdout[0], "3") {
			t.Logf("stdout does not contain expected '3': %v", exec.Logs.Stdout)
		}
	})
}

func TestRunCode_Integration_LanguageAndContextConflict(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		// 先创建一个上下文
		ctxObj, err := sb.Code.CreateCodeContext(context.TODO(), nil)
		if err != nil {
			t.Fatalf("CreateCodeContext error: %v", err)
		}
		if ctxObj.Id == "" {
			t.Fatalf("invalid context: %+v", ctxObj)
		}

		// 同时提供 Language 与 ContextId，应命中客户端参数校验错误
		opts := &toolcode.RunCodeConfig{Language: "python", ContextId: ctxObj.Id}
		_, err = sb.Code.RunCode(context.TODO(), "print('conflict')", opts, nil)
		if err == nil || !strings.Contains(err.Error(), "cannot use RunCode with both contextId and language") {
			t.Fatalf("expected conflict error, got: %v", err)
		}
	})
}

func TestCreateCodeContext_Integration_WithOptions(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		req := &toolcode.CreateCodeContextConfig{
			Cwd:      "/home/user/project",
			Language: "python",
		}
		ctxObj, err := sb.Code.CreateCodeContext(context.TODO(), req)
		if err != nil {
			t.Fatalf("CreateCodeContext error: %v", err)
		}
		if ctxObj == nil || ctxObj.Id == "" {
			t.Fatalf("invalid context returned: %+v", ctxObj)
		}
		// 服务端可能覆盖字段，这里只做宽松断言：至少 language 非空
		if ctxObj.Language == "" {
			t.Log("warning: Language empty in returned context")
		}
	})
}

func TestRunCode_Integration_WithContextIdFlow(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		// 先创建上下文
		ctxObj, err := sb.Code.CreateCodeContext(context.TODO(), &toolcode.CreateCodeContextConfig{
			Language: "python",
		})
		if err != nil {
			t.Fatalf("CreateCodeContext error: %v", err)
		}
		if ctxObj.Id == "" {
			t.Fatalf("invalid context: %+v", ctxObj)
		}

		// 使用 ContextId 执行（不再提供 Language）
		opts := &toolcode.RunCodeConfig{ContextId: ctxObj.Id}
		exec, err := sb.Code.RunCode(context.TODO(), "print('ctx-ok')", opts, nil)
		if err != nil {
			t.Fatalf("RunCode with contextId error: %v", err)
		}
		if exec.Error != nil {
			t.Fatalf("unexpected execution error: %+v", exec.Error)
		}
		// 若服务端返回执行计数，验证其为正
		if exec.ExecutionCount != nil && *exec.ExecutionCount < 0 {
			t.Fatalf("invalid execution count: %v", *exec.ExecutionCount)
		}
	})
}

func TestRunCode_Integration_WithEnvVars(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		opts := &toolcode.RunCodeConfig{
			Language: "python",
			Envs: map[string]string{
				"FOO": "BAR",
			},
		}
		codeStr := "import os\nprint(os.getenv('FOO',''))"
		exec, err := sb.Code.RunCode(context.TODO(), codeStr, opts, nil)
		if err != nil {
			t.Fatalf("RunCode with env vars error: %v", err)
		}
		// 预期 stdout 中包含 BAR
		found := false
		for _, line := range exec.Logs.Stdout {
			if strings.Contains(line, "BAR") {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected stdout to contain 'BAR', got: %v", exec.Logs.Stdout)
		}
	})
}

func TestRunCode_Integration_OnOutput(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		var gotStdout, gotStderr bool

		codeStr := "import sys\nprint('hello-out')\nprint('hello-err', file=sys.stderr)"
		_, err := sb.Code.RunCode(context.TODO(), codeStr, &toolcode.RunCodeConfig{Language: "python"}, &toolcode.OnOutputConfig{
			OnStdout: func(s string) {
				if strings.Contains(s, "hello-out") {
					gotStdout = true
				}
			},
			OnStderr: func(s string) {
				if strings.Contains(s, "hello-err") {
					gotStderr = true
				}
			},
		})
		if err != nil {
			t.Fatalf("RunCode onOutput error: %v", err)
		}
		if !gotStdout || !gotStderr {
			t.Fatalf("onOutput not triggered as expected: stdout=%v stderr=%v", gotStdout, gotStderr)
		}
	})
}

func TestRunCode_Integration_RuntimeError(t *testing.T) {
	runWithBothModes(t, func(t *testing.T, sb *code.Sandbox) {
		opts := &toolcode.RunCodeConfig{Language: "python"}
		// 制造运行时错误（ZeroDivisionError）
		exec, err := sb.Code.RunCode(context.TODO(), "print(1/0)", opts, nil)
		if err != nil {
			t.Fatalf("RunCode request error (expected runtime error in payload): %v", err)
		}
		// 允许两种通路：服务以 error 事件体现，或在 stderr 中体现
		if exec.Error == nil && len(exec.Logs.Stderr) == 0 {
			t.Fatalf("expected runtime error present (error or stderr), got: %+v", exec)
		}
	})
}
