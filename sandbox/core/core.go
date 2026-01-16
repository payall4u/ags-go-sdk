package core

import (
	"context"
	"fmt"

	"github.com/TencentCloudAgentRuntime/ags-go-sdk/connection"
	ags "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
)

// Core contains the core sandbox info and methods.
// This struct should be embedded by sandbox structs of any types within the sandbox package.
// It has necessary info and client needed to contact the Agent Sandbox.
type Core struct {
	client           *ags.Client
	SandboxId        string
	ConnectionConfig *connection.Config
}

// NewCore creates a new Core struct.
func NewCore(client *ags.Client, sandboxId string, connectionConfig *connection.Config) *Core {
	return &Core{
		client:           client,
		SandboxId:        sandboxId,
		ConnectionConfig: connectionConfig,
	}
}

// GetHost returns the hostname of the sandbox.
func (core *Core) GetHost(port int) string {
	return fmt.Sprintf("%v-%v.%v", port, core.SandboxId, core.ConnectionConfig.Domain)
}

// Kill terminates the sandbox.
func (core *Core) Kill(ctx context.Context) error {
	_, err := core.client.StopSandboxInstanceWithContext(ctx, &ags.StopSandboxInstanceRequest{
		InstanceId: &core.SandboxId,
	})
	if err != nil {
		return err
	}
	return nil
}

// SetTimeoutSeconds sets the timeout for the sandbox.
func (core *Core) SetTimeoutSeconds(ctx context.Context, seconds int) error {
	timeoutString := fmt.Sprintf("%ds", seconds)
	_, err := core.client.UpdateSandboxInstanceWithContext(ctx, &ags.UpdateSandboxInstanceRequest{
		InstanceId: &core.SandboxId,
		Timeout:    &timeoutString,
	})
	if err != nil {
		return err
	}
	return nil
}

// GetInfo retrieves detailed information about the sandbox.
func (core *Core) GetInfo(ctx context.Context) (*ags.SandboxInstance, error) {
	info, err := core.client.DescribeSandboxInstanceListWithContext(ctx, &ags.DescribeSandboxInstanceListRequest{
		InstanceIds: []*string{&core.SandboxId},
	})
	if err != nil {
		return nil, err
	}
	if len(info.Response.InstanceSet) > 0 {
		return info.Response.InstanceSet[0], nil
	}
	return nil, fmt.Errorf("not found sandbox instance")
}
