package uartserver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
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

	ln, err = net.ListenTCP("tcp", &net.TCPAddr{
		Port: s.port,
	})
	defer ln.Close()

	for {
		fmt.Printf("Accept uart\n")
		conn, err = ln.Accept()
		if err != nil {
			log.Error().Msg(err.Error())
			continue
		}

		fmt.Printf("Open uart\n")
		if err = s.uart.Open(); err != nil {
			conn.Close()
			log.Error().Msg(err.Error())
			continue
		}

		fmt.Printf("Start txrx\n")
		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			for {
				err := tuctxcopy(ctx, conn, s.uart)
				if err != nil {
					break
				}
			}
			s.uart.Close()
			wg.Done()
		}()

		go func() {
			for {
				err := utctxcopy(ctx, s.uart, conn)
				if err != nil {
					break
				}
			}
			conn.Close()
			wg.Done()
		}()

		wg.Wait()

		select {
		case <-ctx.Done():
			defer func() { fin <- struct{}{} }()
			return nil
		default:
		}
	}

	return nil
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
				return 0, errors.New("close")
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
