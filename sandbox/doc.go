// Package sandbox provides a comprehensive framework for managing and interacting with sandboxes.
//
// This package contains three main sub-packages, each serving different use cases:
//
// # Code Sandbox (sandbox/code)
//
// Use the code package for creating and managing code execution sandboxes. This package provides
// high-level APIs for running code, managing files, and executing commands within isolated environments.
// This is the recommended package for most code execution use cases.
//
//	import "github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/code"
//
//	// Create a code sandbox with custom options
//	sandbox := code.NewSandbox("my-code-sandbox",
//		code.WithTimeout(30*time.Second),
//		code.WithRetryAttempts(3),
//	)
//
// # Browser Sandbox (sandbox/browser)
//
// NOTE: This package is not implemented yet.
//
// The browser package will provide APIs for web automation, browser control, and web-based testing environments.
// This will be the recommended package for browser automation and web testing use cases when implemented.
//
//	import "github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/browser"
//
//	// Create a browser sandbox (future implementation)
//	// sandbox := browser.NewSandbox("my-browser-sandbox",
//	//     browser.WithHeadless(true),
//	//     browser.WithTimeout(60*time.Second),
//	// )
//
// # Core Package (sandbox/core)
//
// The core package is a low-level package that defines the fundamental functionality for managing
// and communicating with sandboxes. Use this package only if you:
//
//   - Want to implement your own custom sandbox type
//   - Need direct control over sandbox management without the convenience APIs
//   - Don't care about the specific sandbox type and only need basic sandbox operations
//   - Require maximum flexibility and are willing to handle low-level details
//
// Most users should prefer the code package over the core package.
//
//	import "github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/core"
//
//	// Create a new sandbox (advanced use case)
//	coreInstance, err := core.Create(ctx, "my-tool",
//		core.WithCredential(credential),
//		core.WithRegion("ap-singapore"),
//	)
//
//	// Connect to existing sandbox
//	coreInstance, err := core.Connect(ctx, "sandbox-id",
//		core.WithCredential(credential),
//		core.WithRegion("ap-singapore"),
//	)
//
//	// List all sandboxes
//	instances, err := core.List(ctx,
//		core.WithCredential(credential),
//		core.WithRegion("ap-singapore"),
//	)
//
//	// Kill a sandbox
//	err := core.Kill(ctx, "sandbox-id",
//		core.WithCredential(credential),
//		core.WithRegion("ap-singapore"),
//	)
//
// # Package Selection Guide
//
//   - For code execution: Use sandbox/code
//   - For browser automation and web testing: Use sandbox/browser (not implemented yet)
//   - For custom sandbox implementations: Use sandbox/core
//   - For maximum control and flexibility: Use sandbox/core
//
// Each package follows the functional option pattern for configuration, providing
// flexible and extensible APIs with sensible defaults.
package sandbox
