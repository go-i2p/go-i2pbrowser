package goi2pbrowser

import (
	"github.com/go-i2p/go-i2ptunnel-config/i2pconv"
	embedding "github.com/go-i2p/go-i2ptunnel/lib/embedding"
	httpclient "github.com/go-i2p/go-i2ptunnel/lib/http/client"
)

type I2PBrowser struct {
	ProfileDir string
	*embedding.Tunnel
}

func NewI2PBrowser(profileDir string) (*I2PBrowser, error) {
	const samAddr = "127.0.0.1:7656"

	cfg := i2pconv.TunnelConfig{
		Name:      "http-proxy",
		Type:      "httpclient",
		Interface: "127.0.0.1",
		Port:      4444,
	}

	rawTunnel, err := httpclient.NewHTTPClient(cfg, samAddr)
	if err != nil {
		return nil, err
	}

	t, err := embedding.Wrap(rawTunnel,
		embedding.WithSAMAddr(samAddr),
		embedding.WithMetricsAddr(":9090"),
	)
	if err != nil {
		return nil, err
	}
	return &I2PBrowser{ProfileDir: profileDir, Tunnel: t}, nil
}
