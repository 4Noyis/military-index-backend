package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// ServiceProxy represents a reverse proxy to a microservice
type ServiceProxy struct {
	ServiceName string
	TargetURL   *url.URL
	Proxy       *httputil.ReverseProxy
}

// NewServiceProxy creates a new reverse proxy for a microservice
func NewServiceProxy(serviceName, targetURL string) (*ServiceProxy, error) {
	url, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse target URL: %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	
	// Customize the proxy director
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = url.Host
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Origin-Host", url.Host)
	}

	// Add error handling
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, fmt.Sprintf("Service unavailable: %v", err), http.StatusServiceUnavailable)
	}

	// Set timeout
	proxy.Transport = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}

	return &ServiceProxy{
		ServiceName: serviceName,
		TargetURL:   url,
		Proxy:       proxy,
	}, nil
}

// ServeHTTP implements the http.Handler interface
func (sp *ServiceProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sp.Proxy.ServeHTTP(w, r)
}
