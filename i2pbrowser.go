package goi2pbrowser

import (
	"net"

	"github.com/go-i2p/go-i2ptunnel-config/i2pconv"
	embedding "github.com/go-i2p/go-i2ptunnel/lib/embedding"
	httpclient "github.com/go-i2p/go-i2ptunnel/lib/http/client"
)

type I2PBrowser struct {
	ProfileDir string
	*embedding.Tunnel
}

func (i *I2PBrowser) Start() error {
	// only start the tunnel if it is not already running
	if i.Tunnel != nil && i.Tunnel.IsRunning() {
		return nil
	}
	// also check for a listener on the tunnels configured port before starting the tunnel
	t := i.Tunnel.Tunnel()
	addr, err := t.LocalAddress()
	if err != nil {
		return err
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		// proceed anyway in this case, it means that I2P is already running and the tunnel is already running, so we can just return nil
		return nil
	}
	ln.Close()
	return i.Tunnel.Start()
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
