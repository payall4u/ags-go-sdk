# SDK 参考

按包列出导出的类型、函数、方法与关键说明。为简明起见，仅列出对外可用成员与必要注释。

**包导入前缀：** `github.com/TencentCloudAgentRuntime/ags-go-sdk`

## 目录

- [connection](#connection)
- [constant](#constant)
- [sandbox/core](#sandboxcore)
- [sandbox/code](#sandboxcode)
- [tool/code](#toolcode远程代码执行)
- [tool/command](#toolcommand命令进程)
- [tool/filesystem](#toolfilesystem文件系统)

---

## connection

**导入：** `github.com/TencentCloudAgentRuntime/ags-go-sdk/connection`

### 类型

```go
type Config struct {
	Domain      string
	AccessToken string
	Headers     http.Header
	Proxy       *url.URL
}
```

### 构造

```go
func NewConfig() *Config
```

### 说明

- 用于各工具客户端的连接信息。
- `Domain` 形如 `"<sandboxId>.<region>.tencentags.com"`（由 `sandbox/core` 生成）。
- `Domain` 在`tool/code`等Client中包含对应的端口，形如 `"{port}-{sandboxId}.{region}.tencentags.com"`
- `AccessToken` 由 `core.Create/Connect` 内部下发。

---

## constant

**导入：** `github.com/TencentCloudAgentRuntime/ags-go-sdk/constant`

### 常量

```go
const (
	AgentSandboxInternalEndpoint = "ags.internal.tencentcloudapi.com"
	EnvdPort                     = 49983
	CodePort                     = 49999
)
```

---

## sandbox/core

**导入：** `github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/core`

### 类型

```go
type Core struct {
	SandboxId        string
	ConnectionConfig *connection.Config
	// 内部字段：client
}
```

### 方法

```go
func NewCore(client *ags.Client, sandboxId string, cfg *connection.Config) *Core
```

创建一个新的 Core 实例。

```go
func (c *Core) GetHost(port int) string
```

生成形如 `"{port}-{sandboxId}.{region}.tencentags.com"` 的域名。

```go
func (c *Core) Kill(ctx context.Context) error
```

终止沙箱实例。

```go
func (c *Core) SetTimeoutSeconds(ctx context.Context, seconds int) error
```

设置沙箱超时时间（秒）。

```go
func (c *Core) GetInfo(ctx context.Context) (*ags.SandboxInstance, error)
```

获取沙箱详细信息。

### 顶层函数（操作）

```go
func Create(ctx context.Context, toolName string, opts ...CreateOption) (*Core, error)
```

创建新沙箱并返回 Core 实例。

```go
func Connect(ctx context.Context, sandboxId string, opts ...ConnectOption) (*Core, error)
```

连接到现有沙箱并返回 Core 实例。

```go
func List(ctx context.Context, opts ...ListOption) ([]*ags.SandboxInstance, error)
```

列出所有沙箱实例。

```go
func Kill(ctx context.Context, sandboxId string, opts ...KillOption) error
```

销毁指定沙箱。

### 可选项（选项接口）

```go
type ClientOption interface {
	CreateOption
	ConnectOption
	ListOption
	KillOption
}
```

```go
func WithClient(client *ags.Client) ClientOption
```

设置 AGS 客户端实例。此选项优先级高于 `WithCredential` 和 `WithRegion`。

```go
func WithCredential(cred common.CredentialIface) ClientOption
```

设置腾讯云凭证。需要与 `WithRegion` 一起使用来创建新的 AGS 客户端。

```go
func WithRegion(region string) ClientOption
```

设置腾讯云区域。需要与 `WithCredential` 一起使用来创建新的 AGS 客户端。默认区域为 `ap-guangzhou`。

### 说明

- 未指定 `WithClient` 时，需同时提供 `WithCredential` 与 `WithRegion` 才可创建 AGS Client。
- `Create` 返回的 Core 内置 `ConnectionConfig`，后续可用于初始化各工具客户端。

---

## sandbox/code

**导入：** `github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/code`

### 类型

```go
type Sandbox struct {
	*core.Core
	Files    *filesystem.Client
	Commands *command.Client
	Code     *code.Client
}
```

代码沙箱实例，嵌入 `core.Core` 并提供三个工具客户端。

### 操作

```go
func Create(ctx context.Context, toolName string, opts ...CreateOption) (*Sandbox, error)
```

创建沙箱并初始化 Files/Commands/Code 三个客户端。

```go
func Connect(ctx context.Context, sandboxId string, opts ...ConnectOption) (*Sandbox, error)
```

连接到现有沙箱并初始化三个客户端。

```go
func List(ctx context.Context, opts ...ListOption) ([]*ags.SandboxInstance, error)
```

列出所有沙箱实例。

```go
func Kill(ctx context.Context, sandboxId string, opts ...KillOption) error
```

销毁指定沙箱。

### 选项

```go
type ClientOption interface {
	CreateOption
	ConnectOption
	ListOption
	KillOption
}
```

```go
func WithClient(client *ags.Client) ClientOption
func WithCredential(cred common.CredentialIface) ClientOption
func WithRegion(region string) ClientOption
```

选项与 `sandbox/core` 中的选项类似，透传到底层 core 操作。

---

## tool/code（远程代码执行）

**导入：** `github.com/TencentCloudAgentRuntime/ags-go-sdk/tool/code`

### 客户端

```go
type Client struct {
	// 内部字段：config, httpClient
}
```

```go
func New(cfg *connection.Config) *Client
```

创建新的代码执行客户端。

### 方法

```go
func (c *Client) RunCode(ctx context.Context, code string, config *RunCodeConfig, onOutput *OnOutputConfig) (*Execution, error)
```

执行代码并返回执行结果。

- 如果 `config` 同时设置 `Language` 与 `ContextId` 会报错。
- 默认 `Language=python`（当未指定 `Language` 与 `ContextId` 时）。
- `onOutput` 可选：实时回调输出（stdout/stderr），不影响最终返回的聚合结果。

```go
func (c *Client) CreateCodeContext(ctx context.Context, config *CreateCodeContextConfig) (*CodeContext, error)
```

创建代码执行上下文。

### 配置类型

```go
type RunCodeConfig struct {
	Language  string
	ContextId string
	Envs      map[string]string
}
```

代码执行配置。

```go
type OnOutputConfig struct {
	OnStdout func(string)
	OnStderr func(string)
}
```

输出回调配置：
- `OnStdout`：接收实时标准输出文本
- `OnStderr`：接收实时标准错误文本
- 若传入 `nil` 或其中任一回调为 `nil`，对应回调不会被调用；最终 `Execution.Logs` 仍保留聚合日志

```go
type CreateCodeContextConfig struct {
	Cwd      string
	Language string
}
```

创建上下文配置。默认 `Cwd="/home/user"`，`Language="python"`。

### 结果类型

```go
type Execution struct {
	Results        []Result
	Logs           Logs
	Error          *ExecutionError
	ExecutionCount *int
}
```

代码执行结果汇总。

```go
type Result struct {
	Text         *string
	Html         *string
	Markdown     *string
	Svg          *string
	Png          *string
	Jpeg         *string
	Pdf          *string
	Latex        *string
	Javascript   *string
	Json         map[string]any
	Data         map[string]any
	Chart        map[string]any
	IsMainResult bool
	Extra        map[string]any
}
```

单个执行结果，支持多种格式。

```go
type Logs struct {
	Stdout []string
	Stderr []string
}
```

标准输出和标准错误日志。

```go
type ExecutionError struct {
	Name      string
	Value     string
	Traceback string
}
```

执行错误信息。

```go
type CodeContext struct {
	Id       string
	Language string
	Cwd      string
}
```

代码上下文信息。

---

## tool/command（命令/进程）

**导入：** `github.com/TencentCloudAgentRuntime/ags-go-sdk/tool/command`

### 客户端

```go
type Client struct {
	// 内部字段：config, httpClient, rpcClient
}
```

```go
func New(cfg *connection.Config) (*Client, error)
```

创建新的命令执行客户端。

### 方法

```go
func (c *Client) Run(ctx context.Context, cmd string, config *ProcessConfig, onOutput *OnOutputConfig) (*Result, error)
```

便捷前台运行命令，聚合 stdout/stderr 与退出码。

```go
func (c *Client) Start(ctx context.Context, cmd string, config *ProcessConfig, onOutput *OnOutputConfig) (*Handle, error)
```

后台运行命令，返回 Handle。

```go
func (c *Client) Connect(ctx context.Context, pid uint32, onOutput *OnOutputConfig) (*Handle, error)
```

连接到现有进程。

```go
func (c *Client) List(ctx context.Context) ([]ProcessInfo, error)
```

列出当前运行中的进程。

### 配置类型

```go
type ProcessConfig struct {
	User string
	Args []string
	Envs map[string]string
	Cwd  *string
}
```

进程配置。`User` 默认为 "user"，仅允许 "user" 或 "root"。

```go
type OnOutputConfig struct {
	OnStdout func([]byte)
	OnStderr func([]byte)
}
```

输出回调配置。

### 结果类型

```go
type ProcessInfo struct {
	Pid  uint32
	Tag  *string
	Cmd  string
	Args []string
	Envs map[string]string
	Cwd  *string
}
```

进程信息。

```go
type ProcessResult struct {
	ExitCode int32
	Error    *string
}
```

进程执行结果。

```go
type Result struct {
	ExitCode int32
	Stdout   []byte
	Stderr   []byte
	Error    *string
}
```

聚合的前台执行结果。

### 句柄

```go
type Handle struct {
	Pid uint32
	// 内部字段：cancel, stream, client, onStdout, onStderr
}
```

进程句柄，用于控制和监控进程。

#### 方法

```go
func (h *Handle) Wait(ctx context.Context) (*ProcessResult, error)
```

等待进程结束并返回结果。

```go
func (h *Handle) Kill(ctx context.Context) error
```

发送 SIGKILL 信号终止进程。

```go
func (h *Handle) Disconnect(ctx context.Context) error
```

断开连接但不终止进程。

```go
func (h *Handle) SendInput(ctx context.Context, pid uint32, stdin []byte) error
```

向进程发送输入。

```go
func (h *Handle) SendSignal(ctx context.Context, pid uint32, sig process.Signal) error
```

向进程发送信号。信号值例如：`15` (SIGTERM)、`9` (SIGKILL)。

---

## tool/filesystem（文件系统）

**导入：** `github.com/TencentCloudAgentRuntime/ags-go-sdk/tool/filesystem`

### 客户端

```go
type Client struct {
	// 内部字段：config, httpClient, rpcClient
}
```

```go
func New(cfg *connection.Config) (*Client, error)
```

创建新的文件系统客户端。

### 方法

```go
func (c *Client) Read(ctx context.Context, path string, config *ReadConfig) (io.Reader, error)
```

读取文件内容。

```go
func (c *Client) Write(ctx context.Context, path string, data io.Reader, config *WriteConfig) (*WriteInfo, error)
```

写入文件。

```go
func (c *Client) List(ctx context.Context, path string, config *ListConfig) ([]EntryInfo, error)
```

列出目录内容。

```go
func (c *Client) Exists(ctx context.Context, path string, config *ExistsConfig) (bool, error)
```

检查文件或目录是否存在。

```go
func (c *Client) GetInfo(ctx context.Context, path string, config *GetInfoConfig) (*EntryInfo, error)
```

获取文件或目录详细信息。

```go
func (c *Client) Remove(ctx context.Context, path string, config *RemoveConfig) error
```

删除文件或目录。

```go
func (c *Client) Rename(ctx context.Context, oldPath, newPath string, config *RenameConfig) error
```

重命名或移动文件/目录。

```go
func (c *Client) MakeDir(ctx context.Context, path string, config *MakeDirConfig) (bool, error)
```

创建目录。返回 `true` 表示创建成功。

### 类型

```go
type FileType string

const (
	File FileType = "file"
	Dir  FileType = "dir"
)
```

文件类型常量。

```go
type WriteInfo struct {
	Name string
	Type *FileType
	Path string
}
```

写入操作返回的基本信息。

```go
type EntryInfo struct {
	WriteInfo
	Size          int64
	Mode          int
	Permissions   string
	Owner         string
	Group         string
	ModifiedTime  time.Time
	SymlinkTarget *string
}
```

文件或目录的详细信息。

### 事件类型（预留）

```go
type FilesystemEventType string

const (
	EventCreate FilesystemEventType = "create"
	EventWrite  FilesystemEventType = "write"
	EventRemove FilesystemEventType = "remove"
	EventRename FilesystemEventType = "rename"
	EventChmod  FilesystemEventType = "chmod"
)
```

文件系统事件类型（当前为预留功能）。

### 配置类型

所有文件系统操作的配置类型：

```go
type ReadConfig struct{ User string }
type WriteConfig struct{ User string }
type ListConfig struct{ Depth int; User string }
type ExistsConfig struct{ User string }
type GetInfoConfig struct{ User string }
type RemoveConfig struct{ User string }
type RenameConfig struct{ User string }
type MakeDirConfig struct{ User string }
```

### 说明

- 所有 `*Config.User` 仅允许 "user" 或 "root"；留空时默认 "user"。
- `Read/Write` 通过 HTTP `/files` 端点；其余通过 RPC 接口完成。
- `ListConfig.Depth` 指定递归深度，默认为 1。
