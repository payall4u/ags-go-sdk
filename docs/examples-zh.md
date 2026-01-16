# 使用示例

本页提供从零开始的常用示例：如何创建沙箱、执行代码、读写文件、运行命令与管理进程。

## 先决条件

- 已具备腾讯云账号与 Agent Sandbox 访问权限
- Go 1.20+
- 已配置凭证（推荐使用 SDK Client 初始化）

## 导入路径约定

- 代码沙箱：`github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/code`
- 浏览器沙箱：`github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/browser`
- 核心包：`github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/core`
- 代码执行：`github.com/TencentCloudAgentRuntime/ags-go-sdk/tool/code`
- 命令/进程：`github.com/TencentCloudAgentRuntime/ags-go-sdk/tool/command`
- 文件系统：`github.com/TencentCloudAgentRuntime/ags-go-sdk/tool/filesystem`

# 代码沙箱 (sandbox/code)

代码沙箱提供高级 API，适用于大多数代码执行、文件操作和命令管理场景。

## 1. 创建代码沙箱并获取三大客户端

推荐通过腾讯云 SDK 的 AGS Client 注入，随后创建沙箱；返回的 Sandbox 聚合了 Files/Commands/Code 三个客户端。

```go
package main

import (
	"context"
	"log"
	"os"

	sandboxcode "github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/code"

	ags "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

func main() {
	ctx := context.Background()

	// 初始化 AGS Client（推荐）
	cred := &common.Credential{
		SecretId:  os.Getenv("TENCENTCLOUD_SECRET_ID"),
		SecretKey: os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	}
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ags.tencentcloudapi.com"
	client, err := ags.NewClient(cred, "ap-guangzhou", cpf)
	if err != nil {
		log.Fatal(err)
	}

	// 创建沙箱（tool 由服务端配置，比如 "code-interpreter-v1"）
	sb, err := sandboxcode.Create(ctx, "code-interpreter-v1", sandboxcode.WithClient(client))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = sb.Kill(ctx) }()

	log.Println("sandbox:", sb.SandboxId)
}
```

**要点**

- `sandbox/code.Create` 封装 `core.Create`，自动初始化 Files/Commands/Code 三个客户端，域名/鉴权统一透传。
- `sb.Kill(ctx)` 用于销毁释放；也可用 `sb.SetTimeoutSeconds(ctx, 600)` 设置过期时间。

## 2. 运行代码（Python 等）

`RunCode` 支持直接传入代码与可选语言；若不提供 ContextId 与 Language，默认 Language=python。

```go
exec, err := sb.Code.RunCode(ctx, "print('hello')", nil, nil)
if err != nil {
	log.Fatal(err)
}
log.Printf("results=%d, stdout=%v, error=%v", len(exec.Results), exec.Logs.Stdout, exec.Error)
```

最简实时输出（onOutput）示例：

```go
_, err := sb.Code.RunCode(ctx, "import time\nfor i in range(3):\n print(i); time.sleep(1)", &code.RunCodeConfig{Language: "python"}, &code.OnOutputConfig{
	OnStdout: func(s string) { log.Print("OUT:", s) },
	OnStderr: func(s string) { log.Print("ERR:", s) },
})
if err != nil {
	log.Fatal(err)
}
```

### 使用持久化代码上下文（可复用执行目录/会话）

```go
// 创建上下文（可指定语言与工作目录）
ctxResp, err := sb.Code.CreateCodeContext(ctx, &code.CreateCodeContextConfig{
	Cwd:      "/home/user/project",
	Language: "python",
})
if err != nil {
	log.Fatal(err)
}

// 绑定上下文执行
exec, err := sb.Code.RunCode(ctx, "print('stateful run')", &code.RunCodeConfig{
	ContextId: ctxResp.Id,
}, nil)
if err != nil {
	log.Fatal(err)
}
```

## 3. 文件系统操作（读/写/列/查/删/改名/建目录）

```go
// 写文件
_, err := sb.Files.Write(ctx, "/home/user/demo.txt", bytes.NewBufferString("hello"), &filesystem.WriteConfig{User: "user"})
if err != nil {
	log.Fatal(err)
}

// 读文件
r, err := sb.Files.Read(ctx, "/home/user/demo.txt", &filesystem.ReadConfig{User: "user"})
if err != nil {
	log.Fatal(err)
}
data, _ := io.ReadAll(r)
log.Println("content:", string(data))

// 列目录
entries, err := sb.Files.List(ctx, "/home/user", &filesystem.ListConfig{Depth: 1, User: "user"})
if err != nil {
	log.Fatal(err)
}
for _, e := range entries {
	log.Printf("[%s] %s size=%d", func() string {
		if e.Type != nil {
			return string(*e.Type)
		}
		return "unknown"
	}(), e.Path, e.Size)
}

// 获取信息
info, err := sb.Files.GetInfo(ctx, "/home/user/demo.txt", &filesystem.GetInfoConfig{User: "user"})
if err != nil {
	log.Fatal(err)
}
log.Printf("owner=%s perms=%s", info.Owner, info.Permissions)

// 判断是否存在
ok, err := sb.Files.Exists(ctx, "/home/user/demo.txt", &filesystem.ExistsConfig{User: "user"})
if err != nil {
	log.Fatal(err)
}
log.Println("exists:", ok)

// 改名
err = sb.Files.Rename(ctx, "/home/user/demo.txt", "/home/user/demo2.txt", &filesystem.RenameConfig{User: "user"})
if err != nil {
	log.Fatal(err)
}

// 创建目录
created, err := sb.Files.MakeDir(ctx, "/home/user/newdir", &filesystem.MakeDirConfig{User: "user"})
if err != nil {
	log.Fatal(err)
}
log.Println("dir created:", created)

// 删除
err = sb.Files.Remove(ctx, "/home/user/demo2.txt", &filesystem.RemoveConfig{User: "user"})
if err != nil {
	log.Fatal(err)
}
```

**注意**

- 所有 `*Config` 的 `User` 字段仅允许 "user" 或 "root"；为空时默认 "user"。
- 访问路径需使用沙箱内文件系统路径。

## 4. 命令/进程管理（前台/后台、输入/信号、进程列表）

### 便捷前台运行（自动聚合 stdout/stderr）

```go
res, err := sb.Commands.Run(ctx, "echo hello && uname -a", &command.ProcessConfig{
	User: "user", // 进程内用户，默认为 "user"
}, &command.OnOutputConfig{
	OnStdout: func(b []byte) { log.Printf("STDOUT: %s", string(b)) },
	OnStderr: func(b []byte) { log.Printf("STDERR: %s", string(b)) },
})
if err != nil {
	log.Fatal(err)
}
log.Printf("exit=%d, stdout=%q, stderr=%q, err=%v", res.ExitCode, string(res.Stdout), string(res.Stderr), res.Error)
```

### 后台启动 + 等待

```go
h, err := sb.Commands.Start(ctx, "sleep 5; echo done", &command.ProcessConfig{
	User: "user",
}, &command.OnOutputConfig{
	OnStdout: func(b []byte) { log.Print("OUT:", string(b)) },
	OnStderr: func(b []byte) { log.Print("ERR:", string(b)) },
})
if err != nil {
	log.Fatal(err)
}

ret, err := h.Wait(ctx)
if err != nil {
	log.Fatal(err)
}
log.Printf("exit=%d, error=%v", ret.ExitCode, ret.Error)
```

### 发送输入与信号

```go
// 创建进程
h, err := sb.Commands.Start(ctx, "cat", &command.ProcessConfig{
	User: "user",
}, &command.OnOutputConfig{
	OnStdout: func(b []byte) { log.Print("OUT:", string(b)) },
	OnStderr: func(b []byte) { log.Print("ERR:", string(b)) },
})
if err != nil {
	log.Fatal(err)
}

// 一次性发送输入到 stdin
_ = h.SendInput(ctx, h.Pid, []byte("hello\n"))

// 发送信号（SIGTERM=15 SIGKILL=9 目前只接受这两种信号）
_ = h.SendSignal(ctx, h.Pid, 15)

// 或直接发送 SIGKILL
_ = h.Kill(ctx)
```

### 列出运行中进程

```go
ps, err := sb.Commands.List(ctx)
if err != nil {
	log.Fatal(err)
}
for _, p := range ps {
	log.Printf("pid=%d cmd=%s args=%v cwd=%v", p.Pid, p.Cmd, p.Args, p.Cwd)
}
```

## 5. 代码沙箱列表和管理

```go
// 列表所有沙箱
instances, err := sandboxcode.List(ctx, sandboxcode.WithClient(client))
if err != nil {
	log.Fatal(err)
}
for _, ins := range instances {
	log.Println(*ins.InstanceId, *ins.Status)
}

// 连接到现有代码沙箱
sb, err := sandboxcode.Connect(ctx, "SBOX-XXXX", sandboxcode.WithClient(client))
if err != nil {
	log.Fatal(err)
}
defer sb.Kill(ctx)
```

**提示**

- `sandbox/code.Create` 已经为你构造好连接域名与 AccessToken 并分配到各工具客户端，无需手动管理。
- 如需自定义代理、额外请求头，可在连接后修改 `sb.Core.ConnectionConfig`，再按需重新创建客户端。

# 浏览器沙箱 (sandbox/browser)

**注意：浏览器沙箱包尚未实现。**

浏览器沙箱将提供 Web 自动化、浏览器控制和基于 Web 的测试环境 API。

# 核心包 (sandbox/core)

核心包提供低级 API，适用于自定义沙箱实现或需要最大控制权的场景。

## 1. 直接使用核心 API 创建沙箱

```go
import "github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/core"

// 创建沙箱
coreInstance, err := core.Create(ctx, "code-interpreter-v1",
	core.WithCredential(cred),
	core.WithRegion("ap-guangzhou"),
)
if err != nil {
	log.Fatal(err)
}
defer coreInstance.Kill(ctx)

log.Println("沙箱ID:", coreInstance.SandboxId)
log.Println("连接配置:", coreInstance.ConnectionConfig)
```

## 2. 连接到现有沙箱

```go
// 连接到指定沙箱（仅获取 token，不初始化工具客户端）
coreInstance, err := core.Connect(ctx, "SBOX-XXXX", 
	core.WithCredential(cred),
	core.WithRegion("ap-guangzhou"),
)
if err != nil {
	log.Fatal(err)
}
defer coreInstance.Kill(ctx)
```

## 3. 列出所有沙箱

```go
// 列出所有沙箱实例
instances, err := core.List(ctx,
	core.WithCredential(cred),
	core.WithRegion("ap-guangzhou"),
)
if err != nil {
	log.Fatal(err)
}

for _, ins := range instances {
	log.Printf("沙箱: %s, 状态: %s", *ins.InstanceId, *ins.Status)
}
```

## 4. 销毁沙箱

```go
// 直接销毁指定沙箱
err := core.Kill(ctx, "SBOX-XXXX",
	core.WithCredential(cred),
	core.WithRegion("ap-guangzhou"),
)
if err != nil {
	log.Fatal(err)
}
```

## 5. 核心实例方法

```go
// 获取沙箱主机地址
host := coreInstance.GetHost(8080)
log.Println("沙箱主机:", host)

// 设置超时时间
err := coreInstance.SetTimeoutSeconds(ctx, 600)
if err != nil {
	log.Fatal(err)
}

// 获取沙箱详细信息
info, err := coreInstance.GetInfo(ctx)
if err != nil {
	log.Fatal(err)
}
log.Printf("沙箱信息: %+v", info)
```

## 使用场景

**使用代码沙箱 (sandbox/code) 当你需要：**

- 代码执行和文件操作
- 命令和进程管理
- 开箱即用的高级 API

**使用浏览器沙箱 (sandbox/browser) 当你需要：**

- Web 自动化和浏览器控制（尚未实现）
- 基于 Web 的测试环境

**使用核心包 (sandbox/core) 当你需要：**

- 实现自定义沙箱类型
- 直接控制沙箱管理
- 最大的灵活性和控制权
- 不关心具体沙箱类型，只需要基本沙箱操作
