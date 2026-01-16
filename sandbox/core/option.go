package core

import (
	"time"

	ags "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
)

const (
	// DefaultDataPlaneDomain is the default domain for sandbox data plane connections
	DefaultDataPlaneDomain = "tencentags.com"
)

// #===================================================================================================================#
// #                                                 ClientOption                                                      #
// #===================================================================================================================#

// configs for Tencent Cloud SDK Agent Sandbox client.
type clientConfig struct {
	// Core functionality options
	// Either set tencentcloud-sdk-go Agent Sandbox client or set Tencent Cloud credentials and region.
	credential common.CredentialIface
	region     string

	client *ags.Client
}

// ClientOption defines a option for configuring Tencent Cloud SDK Agent Sandbox client.
type ClientOption interface {
	CreateOption
	ConnectOption
	ListOption
	KillOption
}

// clientOptionFunc is a function adapter for ClientOption
type clientOptionFunc func(*clientConfig)

func (f clientOptionFunc) applyCreate(config *createConfig) {
	f(config.clientConfig)
}

func (f clientOptionFunc) applyConnect(config *connectConfig) {
	f(config.clientConfig)
}

func (f clientOptionFunc) applyList(config *listConfig) {
	f(config.clientConfig)
}

func (f clientOptionFunc) applyKill(config *killConfig) {
	f(config.clientConfig)
}

// WithClient sets the AGS client instance. This option has higher priority than WithCredential and WithRegion.
// When this option is set, operations are performed using the given AGS client.
func WithClient(client *ags.Client) ClientOption {
	return clientOptionFunc(func(config *clientConfig) {
		config.client = client
	})
}

// WithCredential sets the credentials of the AGS client that will be created in order to call
// Tencent Cloud AgentSandbox APIs to perform operations.
// When WithClient option is not set, this option is used together with WithRegion to create a new AGS client which
// will be used to perform operations.
// However, if WithRegion option is not set, the default region [regions.Guangzhou] will be used.
func WithCredential(credential common.CredentialIface) ClientOption {
	return clientOptionFunc(func(config *clientConfig) {
		config.credential = credential
	})
}

// WithRegion sets the Tencent Cloud region of the AGS client that will be created in order to call Tencent Cloud
// AgentSandbox APIs to perform operations.
// When WithClient option is not set, this option is used together with WithCredential to create a new AGS client which
// will be used to perform operations.
// However, if WithCredential option is not set, an error will be returned by the operation function.
func WithRegion(region string) ClientOption {
	return clientOptionFunc(func(config *clientConfig) {
		config.region = region
	})
}

// #===================================================================================================================#
// #                                                 CreateOption                                                      #
// #===================================================================================================================#

// configs for Create operations.
type createConfig struct {
	*clientConfig
	*dataPlaneConfig
	// Sandbox instance timeout (default: 300s)
	sandboxTimeout time.Duration
}

// evaluateCreateOpts evaluates the provided options and returns the configuration
func evaluateCreateOpts(options []CreateOption) *createConfig {
	config := &createConfig{
		clientConfig: &clientConfig{
			region: regions.Guangzhou,
		},
		dataPlaneConfig: &dataPlaneConfig{
			domain: DefaultDataPlaneDomain,
		},
		sandboxTimeout: 300 * time.Second, // Default timeout of 300 seconds
	}

	for _, option := range options {
		option.applyCreate(config)
	}

	return config
}

// CreateOption defines options specific to Create operations
type CreateOption interface {
	applyCreate(*createConfig)
}

// createOptionFunc is a function adapter for CreateOption
type createOptionFunc func(*createConfig)

func (f createOptionFunc) applyCreate(config *createConfig) {
	f(config)
}

// WithSandboxTimeout sets the timeout for the sandbox instance lifecycle.
// The timeout parameter should be a time.Duration (e.g., 300*time.Second, 5*time.Minute, 1*time.Hour).
// This determines how long the sandbox instance will remain active before automatic termination.
// Default timeout is 300s if not specified.
func WithSandboxTimeout(timeout time.Duration) CreateOption {
	return createOptionFunc(func(config *createConfig) {
		config.sandboxTimeout = timeout
	})
}

// #===================================================================================================================#
// #                                                 ConnectOption                                                     #
// #===================================================================================================================#

// configs for Connect operations.
type connectConfig struct {
	*clientConfig
	*dataPlaneConfig
}

// evaluateConnectOpts evaluates the provided options and returns the configuration
func evaluateConnectOpts(options []ConnectOption) *connectConfig {
	config := &connectConfig{
		clientConfig: &clientConfig{
			region: regions.Guangzhou,
		},
		dataPlaneConfig: &dataPlaneConfig{
			domain: DefaultDataPlaneDomain,
		},
	}

	for _, option := range options {
		option.applyConnect(config)
	}

	return config
}

// ConnectOption defines options specific to Connect operations
type ConnectOption interface {
	applyConnect(*connectConfig)
}

// #===================================================================================================================#
// #                                            DataPlaneOption                                                   #
// #===================================================================================================================#

// dataPlaneConfig is the config for contacting sandboxes. This is the data plane config.
type dataPlaneConfig struct {
	domain string
}

// DataPlaneOption defines options for configuring domain, applicable to Create and Connect operations.
type DataPlaneOption interface {
	CreateOption
	ConnectOption
}

// dataPlaneOptionFunc is a function adapter for DataPlaneOption
type dataPlaneOptionFunc func(*dataPlaneConfig)

func (f dataPlaneOptionFunc) applyCreate(config *createConfig) {
	f(config.dataPlaneConfig)
}

func (f dataPlaneOptionFunc) applyConnect(config *connectConfig) {
	f(config.dataPlaneConfig)
}

// WithDataPlaneDomain sets a custom domain for contacting sandboxes. This sets the data plane domain.
// Default is DefaultDataPlaneDomain.
func WithDataPlaneDomain(domain string) DataPlaneOption {
	return dataPlaneOptionFunc(func(config *dataPlaneConfig) {
		config.domain = domain
	})
}

// #===================================================================================================================#
// #                                                 ListOption                                                        #
// #===================================================================================================================#

// configs for List operations.
type listConfig struct {
	*clientConfig
}

// evaluateListOpts evaluates the provided options and returns the configuration
func evaluateListOpts(options []ListOption) *listConfig {
	config := &listConfig{
		clientConfig: &clientConfig{
			region: regions.Guangzhou,
		},
	}

	for _, option := range options {
		option.applyList(config)
	}

	return config
}

// ListOption defines options specific to List operations
type ListOption interface {
	applyList(*listConfig)
}

// #===================================================================================================================#
// #                                                 KillOption                                                        #
// #===================================================================================================================#

// configs for Kill operations.
type killConfig struct {
	*clientConfig
}

// evaluateKillOpts evaluates the provided options and returns the configuration
func evaluateKillOpts(options []KillOption) *killConfig {
	config := &killConfig{
		clientConfig: &clientConfig{
			region: regions.Guangzhou,
		},
	}

	for _, option := range options {
		option.applyKill(config)
	}

	return config
}

// KillOption defines options specific to Kill operations
type KillOption interface {
	applyKill(*killConfig)
}
