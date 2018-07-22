package uartserver

import (
	"context"
	"io"
	"net"
	"sync"

	nats "github.com/nats-io/go-nats"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware/uartwrap"
	"github.com/yanagiis/GoTuringCoffee/internal/service/mdns"
)

type Service struct {
	uart    uartwrap.UART
	service string
	port    int
	md      *mdns.MDNS
}

func NewService(serviceName string, port int, uart uartwrap.UART, md *mdns.MDNS) *Service {
	md.Register(serviceName, port)
	return &Service{
		uart:    uart,
		service: serviceName,
		port:    port,
	}
}

func (s *Service) Run(ctx context.Context, nc *nats.EncodedConn, fin chan<- struct{}) (err error) {
	var ln net.Listener
	var conn net.Conn

LOOP:
	for {
		ln, err = net.ListenTCP("tcp", &net.TCPAddr{
			Port: s.port,
		})

		conn, err = ln.Accept()
		if err != nil {
			continue
		}

		if err = s.uart.Open(); err != nil {
			conn.Close()
			continue
		}

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			ctxcopy(ctx, conn, s.uart)
			wg.Done()
		}()
		go func() {
			ctxcopy(ctx, s.uart, conn)
			wg.Done()
		}()

		wg.Wait()
		s.uart.Close()
		conn.Close()

		select {
		case <-ctx.Done():
			defer func() { fin <- struct{}{} }()
			break LOOP
		default:
		}
	}
	return nil
}

type readerFunc func(p []byte) (int, error)

func (r readerFunc) Read(p []byte) (int, error) {
	return r(p)
}

func ctxcopy(ctx context.Context, writer io.Writer, reader io.Reader) error {
	_, err := io.Copy(writer, readerFunc(func(p []byte) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return reader.Read(p)
		}
	}))
	return err
}
