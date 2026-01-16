# 贡献指南

感谢您对 Agent Sandbox Go SDK 的关注！我们欢迎任何形式的贡献，包括但不限于：

- 报告 Bug
- 提交功能建议
- 提交代码修复或新功能
- 改进文档

## 开发环境

### 先决条件

- Go 1.22+
- [buf](https://buf.build/) (用于 proto 代码生成)
- 腾讯云账号与 Agent Sandbox 访问权限

### 安装开发工具

```bash
# 安装 buf
make tools
```

### 克隆仓库

```bash
git clone github.com/TencentCloudAgentRuntime/ags-go-sdk.git
cd ags-go-sdk
```

## 开发流程

### 1. 创建分支

请基于 `main` 分支创建您的功能分支：

```bash
git checkout -b feature/your-feature-name
# 或
git checkout -b fix/your-bug-fix
```

分支命名规范：
- `feature/xxx` - 新功能
- `fix/xxx` - Bug 修复
- `docs/xxx` - 文档更新
- `refactor/xxx` - 代码重构
- `chore/xxx` - 构建/工具相关

### 2. 编写代码

请遵循以下规范：

- 遵循 Go 官方代码风格，使用 `gofmt` 格式化代码
- 为公开的函数、类型添加文档注释
- 确保代码通过 `go vet` 和 `golint` 检查
- 编写单元测试覆盖新增代码

### 3. Proto 文件修改

如果您修改了 `proto/` 目录下的 proto 文件，请运行：

```bash
make gen
```

重新生成 Go 代码。

### 4. 运行测试

```bash
cd test
go test -v ./...
```

### 5. 提交代码

提交信息请遵循以下格式：

```
<type>(<scope>): <subject>

<body>

<footer>
```

类型 (type)：
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式调整（不影响代码逻辑）
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建/工具相关

示例：
```
feat(sandbox): 添加浏览器沙箱支持

- 实现 BrowserSandbox 类型
- 添加浏览器操作相关 API
- 新增单元测试

Closes #123
```

### 6. 提交 Merge Request

- 确保所有测试通过
- 确保代码已格式化
- 填写清晰的 MR 描述
- 关联相关的 Issue 或 TAPD 单

## 代码规范

### Go 代码风格

- 使用 `gofmt` 格式化代码
- 遵循 [Effective Go](https://golang.org/doc/effective_go.html)
- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### 命名规范

- 包名：小写，简短，无下划线
- 导出标识符：使用驼峰命名法（CamelCase）
- 非导出标识符：使用小驼峰命名法（camelCase）
- 常量：使用驼峰命名法

### 错误处理

- 优先使用 `errors.New()` 或 `fmt.Errorf()` 创建错误
- 错误信息以小写开头，不以标点结尾
- 使用 `%w` 包装错误以保留错误链

```go
if err != nil {
    return fmt.Errorf("failed to create sandbox: %w", err)
}
```

### 注释规范

- 所有导出的类型、函数、常量必须有文档注释
- 注释以被描述对象的名称开头
- 使用完整的句子

```go
// Create creates a new code sandbox with the given template.
// It returns a Sandbox instance and any error encountered.
func Create(ctx context.Context, template string, opts ...Option) (*Sandbox, error) {
    // ...
}
```

## 目录结构

```
ags-go-sdk/
├── connection/     # 连接管理
├── constant/       # 常量定义
├── docs/           # 文档
├── example/        # 使用示例
├── pb/             # 生成的 protobuf 代码
├── proto/          # proto 定义文件
├── sandbox/        # 沙箱核心实现
│   ├── code/       # 代码沙箱
│   ├── core/       # 核心功能
│   └── browser/    # 浏览器沙箱（待实现）
├── test/           # 测试代码
└── tool/           # 工具客户端
    ├── code/       # 代码执行
    ├── command/    # 命令执行
    └── filesystem/ # 文件系统操作
```

## 报告 Bug

如果您发现了 Bug，请通过以下方式报告：

1. 在项目中创建 Issue
2. 包含以下信息：
   - SDK 版本
   - Go 版本
   - 操作系统
   - 复现步骤
   - 期望行为
   - 实际行为
   - 相关日志或错误信息

## 功能建议

如果您有功能建议，请：

1. 在项目中创建 Issue
2. 描述您的使用场景
3. 说明期望的功能行为

## 联系方式

如有任何问题，请联系项目维护者。

## 许可证

通过提交代码，您同意您的贡献将按照项目许可证进行授权。
