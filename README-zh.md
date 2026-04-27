# ags-go-sdk（payall4u fork）

[`TencentCloudAgentRuntime/ags-go-sdk`](https://github.com/TencentCloudAgentRuntime/ags-go-sdk) 的 fork 版本，去掉了执行命令时的 login shell 包装。

## 解决了什么问题

上游 SDK 执行所有命令时会包装成：

```
/bin/bash -l -c <cmd>
```

`-l` 表示 **login shell**，每次执行命令都会加载 `/etc/profile`、`~/.bash_profile`、`~/.profile` 等登录配置。这会带来明显的启动延迟，并且实际运行环境依赖镜像里的登录脚本——对短小的脚本化命令来说既慢又不可预测。

本 fork 改用普通的 POSIX shell：

```
/bin/sh -c <cmd>
```

不再加载任何登录配置，命令启动更快，运行环境更确定。

其他逻辑与上游完全一致。

## 使用方式

在你的 `go.mod` 里加一条 `replace`，将上游模块路径指向本 fork：

```
require github.com/TencentCloudAgentRuntime/ags-go-sdk v0.1.6

replace github.com/TencentCloudAgentRuntime/ags-go-sdk => github.com/payall4u/ags-go-sdk v0.1.6
```

然后：

```bash
go mod tidy
```

代码里的 import 路径无需改动，继续使用 `github.com/TencentCloudAgentRuntime/ags-go-sdk/...` 即可。

## 上游文档

完整文档、示例和 SDK 参考请见上游仓库：
<https://github.com/TencentCloudAgentRuntime/ags-go-sdk>
