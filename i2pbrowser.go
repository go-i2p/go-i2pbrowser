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

func (i *I2PBrowser) Stop() error {
	// only stop the tunnel if the tunnel is running
	if i.Tunnel != nil {
		// check if the tunnel is running before stopping it
		if i.Tunnel.IsRunning() {
			return i.Tunnel.Stop()
		}
		return i.Tunnel.Stop()
	}
	return nil
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
	ibb := &I2PBrowser{ProfileDir: profileDir, Tunnel: t}
	err = ibb.Start()
	if err != nil {
		return nil, err
	}
	return ibb, nil
}
