package goi2pbrowser

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/go-i2p/go-i2ptunnel-config/i2pconv"
	embedding "github.com/go-i2p/go-i2ptunnel/lib/embedding"
	httpclient "github.com/go-i2p/go-i2ptunnel/lib/http/client"
)

// BrowserOption configures NewI2PBrowser.
type BrowserOption func(*browserConfig)

type browserConfig struct {
	samAddr     string
	metricsAddr string
}

// WithSAMAddr overrides the SAM bridge address (default "127.0.0.1:7656").
func WithSAMAddr(addr string) BrowserOption {
	return func(c *browserConfig) { c.samAddr = addr }
}

// WithMetricsAddr overrides the metrics listener address (default "127.0.0.1:9090").
// Pass an empty string to disable the metrics server.
func WithMetricsAddr(addr string) BrowserOption {
	return func(c *browserConfig) { c.metricsAddr = addr }
}

type I2PBrowser struct {
	ProfileDir string
	*embedding.Tunnel
}

func (i *I2PBrowser) Start() error {
	// guard against a nil tunnel; returning an error is safer than panicking
	if i.Tunnel == nil {
		return fmt.Errorf("i2pbrowser: tunnel is not initialized")
	}
	// only start the tunnel if it is not already running
	if i.Tunnel.IsRunning() {
		return nil
	}
	// also check for a listener on the tunnels configured port before starting the tunnel
	addr, err := i.Tunnel.Tunnel().LocalAddress()
	if err != nil {
		return err
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		// Port is occupied — assume I2P is already running, but warn that the
		// assumption may be wrong (any process holding port 4444 triggers this path).
		log.Printf("WARNING: port %s is already in use; assuming I2P proxy is running. "+
			"If I2P is not running, browser traffic may not be routed through I2P.", addr)
		return nil
	}
	if err := ln.Close(); err != nil {
		log.Printf("WARNING: failed to close probe listener on %s: %v", addr, err)
	}
	// TOCTOU: if the port is grabbed between Close and Start, treat that the
	// same as the probe-failure case above (another I2P proxy is now listening).
	startErr := i.Tunnel.Start()
	if startErr != nil {
		var opErr *net.OpError
		if errors.As(startErr, &opErr) && isAddrInUse(opErr) {
			log.Printf("WARNING: port %s was taken after probe; assuming I2P proxy is running. "+
				"If I2P is not running, browser traffic may not be routed through I2P.", addr)
			return nil
		}
		return startErr
	}
	return nil
}

// isAddrInUse reports whether a *net.OpError wraps an "address already in use" error.
func isAddrInUse(opErr *net.OpError) bool {
	return opErr != nil && opErr.Op == "listen" &&
		opErr.Err != nil && (containsSuffix(opErr.Err.Error(), "address already in use") ||
		containsSuffix(opErr.Err.Error(), "only one usage of each socket"))
}

// containsSuffix reports whether s contains substr.
func containsSuffix(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func (i *I2PBrowser) Stop() error {
	// only stop the tunnel if the tunnel is running
	if i.Tunnel != nil && i.Tunnel.IsRunning() {
		return i.Tunnel.Stop()
	}
	return nil
}

// NewI2PBrowser creates and starts an I2PBrowser with the given profile directory.
// Optional BrowserOption values override the default SAM address and metrics address.
func NewI2PBrowser(profileDir string, opts ...BrowserOption) (*I2PBrowser, error) {
	cfg := browserConfig{
		samAddr:     "127.0.0.1:7656",
		metricsAddr: "127.0.0.1:9090",
	}
	for _, o := range opts {
		o(&cfg)
	}

	tunnelCfg := i2pconv.TunnelConfig{
		Name:      "http-proxy",
		Type:      "httpclient",
		Interface: "127.0.0.1",
		Port:      4444,
	}

	rawTunnel, err := httpclient.NewHTTPClient(tunnelCfg, cfg.samAddr)
	if err != nil {
		return nil, err
	}

	wrapOpts := []embedding.Option{embedding.WithSAMAddr(cfg.samAddr)}
	if cfg.metricsAddr != "" {
		wrapOpts = append(wrapOpts, embedding.WithMetricsAddr(cfg.metricsAddr))
	}

	t, err := embedding.Wrap(rawTunnel, wrapOpts...)
	if err != nil {
		return nil, err
	}
	ibb := &I2PBrowser{ProfileDir: profileDir, Tunnel: t}
	if err := ibb.Start(); err != nil {
		return nil, err
	}
	return ibb, nil
}
