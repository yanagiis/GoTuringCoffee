package mdns

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/grandcat/zeroconf"
)

type Request struct {
	serviceName string
	reply       chan Response
	timeout     time.Duration
}

type Response struct {
	IPv4 []net.IP
	Port int
	Err  error
}

type Config struct {
	Instance string `mapstructure:"instance"`
	Domain   string `mapstructure:"domain"`
}

type MDNS struct {
	Conf     Config
	services map[string]*zeroconf.Server
	resolver *zeroconf.Resolver
}

func NewMDNS(conf Config) *MDNS {
	return &MDNS{
		Conf:     conf,
		services: make(map[string]*zeroconf.Server),
	}
}

func (m *MDNS) Register(service string, port int) (err error) {
	var server *zeroconf.Server
	if server, err = zeroconf.Register(m.Conf.Instance, service, m.Conf.Domain, port, nil, nil); err != nil {
		return
	}
	m.services[service] = server
	return
}

func (m *MDNS) Stop() {
	for name, service := range m.services {
		service.Shutdown()
		delete(m.services, name)
	}
}

func (m *MDNS) Lookup(serviceName string, timeout time.Duration) ([]net.IP, int, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	entries := make(chan *zeroconf.ServiceEntry)
	if err = resolver.Browse(ctx, serviceName, "local.", entries); err != nil {
		return nil, 0, err
	}
	<-ctx.Done()

	for entry := range entries {
		return entry.AddrIPv4, entry.Port, nil
	}

	return nil, 0, fmt.Errorf("Cannot find %q address", serviceName)
}
