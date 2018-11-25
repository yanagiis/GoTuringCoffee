package uartwrap

import (
	"context"
	"fmt"
	"net"
	"time"

	"GoTuringCoffee/internal/service/mdns"

	"github.com/rs/zerolog/log"
)

type TCPUARTClient struct {
	service string
	md      *mdns.MDNS
	conn    net.Conn
}

func NewTCPUARTClientMDNS(service string, md *mdns.MDNS) *TCPUARTClient {
	return &TCPUARTClient{
		service: service,
		md:      md,
	}
}

func (c *TCPUARTClient) Open(ctx context.Context) (err error) {
	var addrs []net.IP
	var port int
	if addrs, port, err = c.md.Lookup(ctx, c.service, time.Second); err != nil {
		err = fmt.Errorf("Cannot lookup %q service", c.service)
		return
	}
	if len(addrs) == 0 {
		err = fmt.Errorf("Cannot lookup %q ip and port", c.service)
		log.Error().Msg(err.Error())
		return
	}

	url := fmt.Sprintf("%s:%d", addrs[0], port)
	if c.conn, err = net.Dial("tcp", url); err != nil {
		log.Error().Msg(err.Error())
		return
	}
	return
}

func (c *TCPUARTClient) IsOpen() bool {
	return c.conn != nil
}

func (c *TCPUARTClient) Close() (err error) {
	if c.conn != nil {
		if err = c.conn.Close(); err != nil {
			return
		}
		c.conn = nil
	}
	return
}

func (c *TCPUARTClient) Read(p []byte) (int, error) {
	return c.conn.Read(p)
}

func (c *TCPUARTClient) Write(p []byte) (int, error) {
	return c.conn.Write(p)
}
