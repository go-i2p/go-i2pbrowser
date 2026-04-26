package goi2pbrowser

import (
	"fmt"
	"log"
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
	return i.Tunnel.Start()
}

func (i *I2PBrowser) Stop() error {
	// only stop the tunnel if the tunnel is running
	if i.Tunnel != nil && i.Tunnel.IsRunning() {
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
		embedding.WithMetricsAddr("127.0.0.1:9090"),
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
