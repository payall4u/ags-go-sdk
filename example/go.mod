module test

go 1.22

replace github.com/TencentCloudAgentRuntime/ags-go-sdk v0.0.0 => ../

require (
	github.com/TencentCloudAgentRuntime/ags-go-sdk v0.0.0
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags v1.3.16
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.3.16
)

require (
	connectrpc.com/connect v1.18.1 // indirect
	google.golang.org/protobuf v1.36.7 // indirect
)
