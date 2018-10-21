package uartserver

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"

	"GoTuringCoffee/internal/hardware/uartwrap"
	"GoTuringCoffee/internal/service/mdns"

	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
)

type Service struct {
	uart    uartwrap.UART
	service string
	port    int
	md      *mdns.MDNS
	ln      net.Listener
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
	var conn net.Conn

	s.ln, err = net.ListenTCP("tcp", &net.TCPAddr{
		Port: s.port,
	})
	defer s.ln.Close()

	for {
		select {
		case <-ctx.Done():
			defer func() { fin <- struct{}{} }()
			return nil
		default:
		}

		conn, err = s.ln.Accept()
		if err != nil {
			log.Error().Err(err).Msg("Accept conn failed")
			continue
		}

		if err = s.uart.Open(ctx); err != nil {
			log.Error().Err(err).Msg("Open uart failed")
			conn.Close()
			log.Error().Msg(err.Error())
			continue
		}

		wg := sync.WaitGroup{}
		wg.Add(2)

		newctx, cancel := context.WithCancel(ctx)

		go func() {
			for {
				err := tuctxcopy(newctx, conn, s.uart)
				if err != nil {
					break
				}
			}
			cancel()
			wg.Done()
		}()

		go func() {
			for {
				err := utctxcopy(newctx, s.uart, conn)
				if err != nil {
					break
				}
			}
			cancel()
			wg.Done()
		}()

		wg.Wait()
	}

	return nil
}

func (s *Service) Stop() error {
	return s.ln.Close()
}

type readerFunc func(p []byte) (int, error)

func (r readerFunc) Read(p []byte) (int, error) {
	return r(p)
}

func utctxcopy(ctx context.Context, writer io.Writer, reader io.Reader) error {
	_, err := io.Copy(writer, readerFunc(func(p []byte) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			n, e := reader.Read(p)
			if e != nil {
				return 0, errors.New("closed")
			}
			return n, e
		}
	}))
	return err
}

func tuctxcopy(ctx context.Context, writer io.Writer, reader io.Reader) error {
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
