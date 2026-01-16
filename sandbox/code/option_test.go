package code

import (
	"testing"
)

func TestWithDataPlaneDomain_AppliesOption(t *testing.T) {
	customDomain := "internal.tencentags.com"

	// Test WithDataPlaneDomain applies to CreateOption
	createConfig := evaluateCreateOpts([]CreateOption{WithDataPlaneDomain(customDomain)})
	if len(createConfig.dataPlaneOptions) != 1 {
		t.Errorf("expected 1 core contact option, got %d", len(createConfig.dataPlaneOptions))
	}

	// Test WithDataPlaneDomain applies to ConnectOption
	connectConfig := evaluateConnectOpts([]ConnectOption{WithDataPlaneDomain(customDomain)})
	if len(connectConfig.dataPlaneOptions) != 1 {
		t.Errorf("expected 1 core contact option, got %d", len(connectConfig.dataPlaneOptions))
	}
}

func TestWithRegion_AppliesOption(t *testing.T) {
	region := "ap-shanghai"

	// Test WithRegion applies to CreateOption
	createConfig := evaluateCreateOpts([]CreateOption{WithRegion(region)})
	if len(createConfig.clientOptions) != 1 {
		t.Errorf("expected 1 client option, got %d", len(createConfig.clientOptions))
	}

	// Test WithRegion applies to ConnectOption
	connectConfig := evaluateConnectOpts([]ConnectOption{WithRegion(region)})
	if len(connectConfig.clientOptions) != 1 {
		t.Errorf("expected 1 client option, got %d", len(connectConfig.clientOptions))
	}

	// Test WithRegion applies to ListOption
	listConfig := evaluateListOpts([]ListOption{WithRegion(region)})
	if len(listConfig.clientOptions) != 1 {
		t.Errorf("expected 1 client option, got %d", len(listConfig.clientOptions))
	}

	// Test WithRegion applies to KillOption
	killConfig := evaluateKillOpts([]KillOption{WithRegion(region)})
	if len(killConfig.clientOptions) != 1 {
		t.Errorf("expected 1 client option, got %d", len(killConfig.clientOptions))
	}
}

func TestMultipleOptions_ApplyInOrder(t *testing.T) {
	// Test multiple options are applied
	createConfig := evaluateCreateOpts([]CreateOption{
		WithRegion("ap-shanghai"),
		WithDataPlaneDomain("internal.tencentags.com"),
	})
	if len(createConfig.clientOptions) != 1 {
		t.Errorf("expected 1 client option, got %d", len(createConfig.clientOptions))
	}
	if len(createConfig.dataPlaneOptions) != 1 {
		t.Errorf("expected 1 core contact option, got %d", len(createConfig.dataPlaneOptions))
	}
}

func TestEvaluateCreateOpts_InitializesCorrectly(t *testing.T) {
	config := evaluateCreateOpts(nil)

	if config.clientConfig == nil {
		t.Fatal("expected clientConfig to be initialized")
	}
	if config.clientOptions == nil {
		t.Fatal("expected clientOptions to be initialized")
	}
	if config.dataPlaneConfig == nil {
		t.Fatal("expected dataPlaneConfig to be initialized")
	}
	if config.dataPlaneOptions == nil {
		t.Fatal("expected dataPlaneOptions to be initialized")
	}
	if config.coreCreateOptions == nil {
		t.Fatal("expected coreCreateOptions to be initialized")
	}
}

func TestEvaluateConnectOpts_InitializesCorrectly(t *testing.T) {
	config := evaluateConnectOpts(nil)

	if config.clientConfig == nil {
		t.Fatal("expected clientConfig to be initialized")
	}
	if config.clientOptions == nil {
		t.Fatal("expected clientOptions to be initialized")
	}
	if config.dataPlaneConfig == nil {
		t.Fatal("expected dataPlaneConfig to be initialized")
	}
	if config.dataPlaneOptions == nil {
		t.Fatal("expected dataPlaneOptions to be initialized")
	}
	if config.coreConnectOptions == nil {
		t.Fatal("expected coreConnectOptions to be initialized")
	}
}

func TestEvaluateListOpts_InitializesCorrectly(t *testing.T) {
	config := evaluateListOpts(nil)

	if config.clientConfig == nil {
		t.Fatal("expected clientConfig to be initialized")
	}
	if config.clientOptions == nil {
		t.Fatal("expected clientOptions to be initialized")
	}
}

func TestEvaluateKillOpts_InitializesCorrectly(t *testing.T) {
	config := evaluateKillOpts(nil)

	if config.clientConfig == nil {
		t.Fatal("expected clientConfig to be initialized")
	}
	if config.clientOptions == nil {
		t.Fatal("expected clientOptions to be initialized")
	}
}

func TestDataPlaneOption_ImplementsInterfaces(_ *testing.T) {
	// Test that WithDataPlaneDomain returns a type that implements both CreateOption and ConnectOption
	option := WithDataPlaneDomain("test.domain.com")

	// Test it implements CreateOption
	var _ CreateOption = option
	// Test it implements ConnectOption
	var _ ConnectOption = option
	// Test it implements DataPlaneOption
	var _ DataPlaneOption = option
}

func TestWithDataPlaneDomain_FunctionAdapterPattern(t *testing.T) {
	// Test that the function adapter pattern works correctly
	customDomain := "adapter.test.com"

	createConfig := evaluateCreateOpts([]CreateOption{WithDataPlaneDomain(customDomain)})
	if createConfig.dataPlaneConfig == nil {
		t.Fatal("expected dataPlaneConfig to be initialized")
	}
	if len(createConfig.dataPlaneOptions) != 1 {
		t.Errorf("expected 1 core contact option, got %d", len(createConfig.dataPlaneOptions))
	}

	connectConfig := evaluateConnectOpts([]ConnectOption{WithDataPlaneDomain(customDomain)})
	if connectConfig.dataPlaneConfig == nil {
		t.Fatal("expected dataPlaneConfig to be initialized")
	}
	if len(connectConfig.dataPlaneOptions) != 1 {
		t.Errorf("expected 1 core contact option, got %d", len(connectConfig.dataPlaneOptions))
	}
}

func TestWithDataPlaneDomain_MultipleApplications(t *testing.T) {
	// Test that applying WithDataPlaneDomain multiple times appends multiple options
	firstDomain := "first.domain.com"
	secondDomain := "second.domain.com"

	config := evaluateCreateOpts([]CreateOption{
		WithDataPlaneDomain(firstDomain),
		WithDataPlaneDomain(secondDomain),
	})

	// Both options should be appended
	if len(config.dataPlaneOptions) != 2 {
		t.Errorf("expected 2 core contact options, got %d", len(config.dataPlaneOptions))
	}
}

func TestWithDataPlaneDomain_WithOtherOptions(t *testing.T) {
	// Test that WithDataPlaneDomain works correctly when combined with other options
	customDomain := "combined.test.com"
	customRegion := "ap-shanghai"

	config := evaluateCreateOpts([]CreateOption{
		WithRegion(customRegion),
		WithDataPlaneDomain(customDomain),
	})

	if len(config.clientOptions) != 1 {
		t.Errorf("expected 1 client option, got %d", len(config.clientOptions))
	}
	if len(config.dataPlaneOptions) != 1 {
		t.Errorf("expected 1 core contact option, got %d", len(config.dataPlaneOptions))
	}
}

func TestDataPlaneConfig_Initialization(t *testing.T) {
	// Test that dataPlaneConfig is properly initialized in createConfig
	createConfig := evaluateCreateOpts(nil)
	if createConfig.dataPlaneConfig == nil {
		t.Fatal("expected dataPlaneConfig to be initialized in createConfig")
	}
	if createConfig.dataPlaneOptions == nil {
		t.Fatal("expected dataPlaneOptions to be initialized")
	}
	if len(createConfig.dataPlaneOptions) != 0 {
		t.Errorf("expected 0 core contact options initially, got %d", len(createConfig.dataPlaneOptions))
	}

	// Test that dataPlaneConfig is properly initialized in connectConfig
	connectConfig := evaluateConnectOpts(nil)
	if connectConfig.dataPlaneConfig == nil {
		t.Fatal("expected dataPlaneConfig to be initialized in connectConfig")
	}
	if connectConfig.dataPlaneOptions == nil {
		t.Fatal("expected dataPlaneOptions to be initialized")
	}
	if len(connectConfig.dataPlaneOptions) != 0 {
		t.Errorf("expected 0 core contact options initially, got %d", len(connectConfig.dataPlaneOptions))
	}
}

func TestWithDataPlaneDomain_AppendsToCorrectSlice(t *testing.T) {
	// Test that WithDataPlaneDomain appends to dataPlaneOptions, not coreCreateOptions or coreConnectOptions
	customDomain := "slice.test.com"

	createConfig := evaluateCreateOpts([]CreateOption{WithDataPlaneDomain(customDomain)})
	if len(createConfig.dataPlaneOptions) != 1 {
		t.Errorf("expected 1 core contact option, got %d", len(createConfig.dataPlaneOptions))
	}
	if len(createConfig.coreCreateOptions) != 0 {
		t.Errorf("expected 0 core create options, got %d", len(createConfig.coreCreateOptions))
	}

	connectConfig := evaluateConnectOpts([]ConnectOption{WithDataPlaneDomain(customDomain)})
	if len(connectConfig.dataPlaneOptions) != 1 {
		t.Errorf("expected 1 core contact option, got %d", len(connectConfig.dataPlaneOptions))
	}
	if len(connectConfig.coreConnectOptions) != 0 {
		t.Errorf("expected 0 core connect options, got %d", len(connectConfig.coreConnectOptions))
	}
}

func TestWithDataPlaneDomain_DoesNotApplyToListOrKill(_ *testing.T) {
	// WithDataPlaneDomain should only implement CreateOption and ConnectOption, not ListOption or KillOption
	// This is a compile-time check, but we can verify the behavior

	// These should compile
	var _ CreateOption = WithDataPlaneDomain("test.com")
	var _ ConnectOption = WithDataPlaneDomain("test.com")

	// These should NOT compile (commented out to avoid compilation errors)
	// var _ ListOption = WithDataPlaneDomain("test.com")  // Should not compile
	// var _ KillOption = WithDataPlaneDomain("test.com")  // Should not compile
}

func TestWithDataPlaneDomain_CombinedWithSandboxTimeout(t *testing.T) {
	// Test that WithDataPlaneDomain and WithSandboxTimeout work together
	customDomain := "timeout.test.com"

	config := evaluateCreateOpts([]CreateOption{
		WithDataPlaneDomain(customDomain),
		WithSandboxTimeout(600),
	})

	if len(config.dataPlaneOptions) != 1 {
		t.Errorf("expected 1 core contact option, got %d", len(config.dataPlaneOptions))
	}
	if len(config.coreCreateOptions) != 1 {
		t.Errorf("expected 1 core create option, got %d", len(config.coreCreateOptions))
	}
}
