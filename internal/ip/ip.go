package ip

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type IPClient struct {
	Provider string
}

func NewIPClient(provider string) *IPClient {
	return &IPClient{
		Provider: provider,
	}
}

func (ic *IPClient) Get() (net.IP, error) {
	httpClient := http.Client{
		Timeout: 5,
	}

	url, err := url.Parse(ic.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	resp, err := httpClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    url,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get IP from provider: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	ip := net.ParseIP(strings.TrimSpace(string(body)))
	if ip == nil {
		return nil, fmt.Errorf("failed to parse IP from response: %s", string(body))
	}

	return ip, nil
}
