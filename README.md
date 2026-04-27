# ags-go-sdk (payall4u fork)

A fork of [`TencentCloudAgentRuntime/ags-go-sdk`](https://github.com/TencentCloudAgentRuntime/ags-go-sdk) that drops the login-shell wrapper when running commands.

## Why this fork

Upstream wraps every command as:

```
/bin/bash -l -c <cmd>
```

The `-l` flag starts a **login shell**, which sources `/etc/profile`, `~/.bash_profile`, `~/.profile`, etc. on every command. This adds noticeable startup latency and makes the runtime environment depend on whatever the sandbox image happens to have in its login profile — surprising for short, scripted commands.

This fork uses a plain POSIX shell instead:

```
/bin/sh -c <cmd>
```

No login profile is sourced. Commands start faster and run in a predictable, minimal environment.

Everything else is identical to upstream.

## Usage

Add a `replace` directive to your `go.mod` so the upstream module path resolves to this fork:

```
require github.com/TencentCloudAgentRuntime/ags-go-sdk v0.1.6

replace github.com/TencentCloudAgentRuntime/ags-go-sdk => github.com/payall4u/ags-go-sdk v0.1.6
```

Then:

```bash
go mod tidy
```

Your import paths stay the same — keep using `github.com/TencentCloudAgentRuntime/ags-go-sdk/...`.

## Upstream

For documentation, examples, and the SDK reference, see the upstream repo:
<https://github.com/TencentCloudAgentRuntime/ags-go-sdk>
