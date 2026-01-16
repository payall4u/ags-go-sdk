package core

import (
	"testing"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
)

func TestEvaluateCreateOpts_DefaultDomain(t *testing.T) {
	config := evaluateCreateOpts(nil)

	if config.domain != "tencentags.com" {
		t.Errorf("expected default domain 'tencentags.com', got '%s'", config.domain)
	}
	if config.clientConfig.region != regions.Guangzhou {
		t.Errorf("expected default region '%s', got '%s'", regions.Guangzhou, config.clientConfig.region)
	}
}

func TestEvaluateCreateOpts_WithDataPlaneDomain(t *testing.T) {
	customDomain := "internal.tencentags.com"
	config := evaluateCreateOpts([]CreateOption{WithDataPlaneDomain(customDomain)})

	if config.domain != customDomain {
		t.Errorf("expected domain '%s', got '%s'", customDomain, config.domain)
	}
}

func TestEvaluateConnectOpts_DefaultDomain(t *testing.T) {
	config := evaluateConnectOpts(nil)

	if config.domain != "tencentags.com" {
		t.Errorf("expected default domain 'tencentags.com', got '%s'", config.domain)
	}
	if config.clientConfig.region != regions.Guangzhou {
		t.Errorf("expected default region '%s', got '%s'", regions.Guangzhou, config.clientConfig.region)
	}
}

func TestEvaluateConnectOpts_WithDataPlaneDomain(t *testing.T) {
	customDomain := "internal.tencentags.com"
	config := evaluateConnectOpts([]ConnectOption{WithDataPlaneDomain(customDomain)})

	if config.domain != customDomain {
		t.Errorf("expected domain '%s', got '%s'", customDomain, config.domain)
	}
}

func TestEvaluateListOpts_DefaultRegion(t *testing.T) {
	config := evaluateListOpts(nil)

	if config.clientConfig.region != regions.Guangzhou {
		t.Errorf("expected default region '%s', got '%s'", regions.Guangzhou, config.clientConfig.region)
	}
}

func TestEvaluateKillOpts_DefaultRegion(t *testing.T) {
	config := evaluateKillOpts(nil)

	if config.clientConfig.region != regions.Guangzhou {
		t.Errorf("expected default region '%s', got '%s'", regions.Guangzhou, config.clientConfig.region)
	}
}

func TestWithDataPlaneDomain_OverridesDefault(t *testing.T) {
	// Test that WithDataPlaneDomain can override the default domain
	tests := []struct {
		name   string
		domain string
	}{
		{"internal domain", "internal.tencentags.com"},
		{"custom domain", "sandbox.example.com"},
		{"empty domain", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := evaluateCreateOpts([]CreateOption{WithDataPlaneDomain(tt.domain)})
			if config.domain != tt.domain {
				t.Errorf("expected domain '%s', got '%s'", tt.domain, config.domain)
			}
		})
	}
}

func TestWithDataPlaneDomain_AppliesToCreateAndConnect(t *testing.T) {
	customDomain := "custom.tencentags.com"

	// Test Create
	createConfig := evaluateCreateOpts([]CreateOption{WithDataPlaneDomain(customDomain)})
	if createConfig.domain != customDomain {
		t.Errorf("Create: expected domain '%s', got '%s'", customDomain, createConfig.domain)
	}

	// Test Connect
	connectConfig := evaluateConnectOpts([]ConnectOption{WithDataPlaneDomain(customDomain)})
	if connectConfig.domain != customDomain {
		t.Errorf("Connect: expected domain '%s', got '%s'", customDomain, connectConfig.domain)
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
	if createConfig.domain != customDomain {
		t.Errorf("expected domain '%s', got '%s'", customDomain, createConfig.domain)
	}

	connectConfig := evaluateConnectOpts([]ConnectOption{WithDataPlaneDomain(customDomain)})
	if connectConfig.dataPlaneConfig == nil {
		t.Fatal("expected dataPlaneConfig to be initialized")
	}
	if connectConfig.domain != customDomain {
		t.Errorf("expected domain '%s', got '%s'", customDomain, connectConfig.domain)
	}
}

func TestWithDataPlaneDomain_MultipleApplications(t *testing.T) {
	// Test that applying WithDataPlaneDomain multiple times uses the last value
	firstDomain := "first.domain.com"
	secondDomain := "second.domain.com"

	config := evaluateCreateOpts([]CreateOption{
		WithDataPlaneDomain(firstDomain),
		WithDataPlaneDomain(secondDomain),
	})

	if config.domain != secondDomain {
		t.Errorf("expected domain '%s' (last applied), got '%s'", secondDomain, config.domain)
	}
}

func TestWithDataPlaneDomain_WithOtherOptions(t *testing.T) {
	// Test that WithDataPlaneDomain works correctly when combined with other options
	customDomain := "combined.test.com"
	customRegion := regions.Shanghai

	config := evaluateCreateOpts([]CreateOption{
		WithRegion(customRegion),
		WithDataPlaneDomain(customDomain),
	})

	if config.domain != customDomain {
		t.Errorf("expected domain '%s', got '%s'", customDomain, config.domain)
	}
	if config.region != customRegion {
		t.Errorf("expected region '%s', got '%s'", customRegion, config.region)
	}
}

func TestDefaultDataPlaneDomain_Constant(t *testing.T) {
	// Test that the constant is set correctly
	if DefaultDataPlaneDomain != "tencentags.com" {
		t.Errorf("expected DefaultDataPlaneDomain to be 'tencentags.com', got '%s'", DefaultDataPlaneDomain)
	}

	// Test that configs use the constant
	createConfig := evaluateCreateOpts(nil)
	if createConfig.domain != DefaultDataPlaneDomain {
		t.Errorf("expected default domain to use DefaultDataPlaneDomain constant")
	}

	connectConfig := evaluateConnectOpts(nil)
	if connectConfig.domain != DefaultDataPlaneDomain {
		t.Errorf("expected default domain to use DefaultDataPlaneDomain constant")
	}
}

func TestDataPlaneConfig_Initialization(t *testing.T) {
	// Test that dataPlaneConfig is properly initialized in createConfig
	createConfig := evaluateCreateOpts(nil)
	if createConfig.dataPlaneConfig == nil {
		t.Fatal("expected dataPlaneConfig to be initialized in createConfig")
	}
	if createConfig.dataPlaneConfig.domain != DefaultDataPlaneDomain {
		t.Errorf("expected default domain '%s', got '%s'", DefaultDataPlaneDomain, createConfig.dataPlaneConfig.domain)
	}

	// Test that dataPlaneConfig is properly initialized in connectConfig
	connectConfig := evaluateConnectOpts(nil)
	if connectConfig.dataPlaneConfig == nil {
		t.Fatal("expected dataPlaneConfig to be initialized in connectConfig")
	}
	if connectConfig.dataPlaneConfig.domain != DefaultDataPlaneDomain {
		t.Errorf("expected default domain '%s', got '%s'", DefaultDataPlaneDomain, connectConfig.dataPlaneConfig.domain)
	}
}
