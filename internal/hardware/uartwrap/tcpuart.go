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

func (c *TCPUARTClient) Open() (err error) {
	var addrs []net.IP
	var port int
	if addrs, port, err = c.md.Lookup(c.service, time.Second); err != nil {
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
	if err = c.conn.Close(); err != nil {
		return
	}
	c.conn = nil
	return
}

func (c *TCPUARTClient) Read(p []byte) (int, error) {
	return c.conn.Read(p)
}

func (c *TCPUARTClient) Write(p []byte) (int, error) {
	return c.conn.Write(p)
}

type TCPUARTServer struct {
	Service string
	Port    int
	md      *mdns.MDNS
	conn    net.Conn
	uart    UART
	ctx     context.Context
}

func NewTCPUARTServerMDNS(service string, port int, uart UART, md *mdns.MDNS) *TCPUARTServer {
	server := &TCPUARTServer{
		Service: service,
		md:      md,
		uart:    uart,
	}
	if err := md.Register(service, port); err != nil {
		return nil
	}
	return server
}

func (s *TCPUARTServer) Pair(timeout time.Duration) (conn *net.Conn, err error) {
	var ln *net.TCPListener

	if err = s.uart.Open(); err != nil {
		return
	}
	tcpAddr := net.TCPAddr{
		Port: s.Port,
	}
	if ln, err = net.ListenTCP("tcp", &tcpAddr); err != nil {
		s.uart.Close()
		return
	}
	defer ln.Close()
	if s.conn, err = ln.Accept(); err != nil {
		s.uart.Close()
		return
	}
	return
}

func (s *TCPUARTServer) Unpair() (err error) {
	s.uart.Close()
	err = s.conn.Close()
	return
}
