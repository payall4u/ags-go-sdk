package code

import (
	"context"
	"io"
	"net/http"

	"github.com/TencentCloudAgentRuntime/ags-go-sdk/connection"
)

// newHttpClient creates an HTTP client with proxy support
func newHttpClient(config *connection.Config) *http.Client {
	httpClient := &http.Client{}
	if config.Proxy != nil {
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(config.Proxy),
		}
	}
	return httpClient
}

// newHttpRequestWithHeaders creates an HTTP request with configured headers
func newHttpRequestWithHeaders(
	ctx context.Context, method string, url string, body io.Reader, cfg *connection.Config,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	for k, vv := range cfg.Headers {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}
	if cfg.AccessToken != "" {
		req.Header.Set("X-Access-Token", cfg.AccessToken)
	}
	return req, err
}
