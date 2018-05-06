package uartserver

***REMOVED***
	"context"
	"io"
	"net"
	"sync"

	nats "github.com/nats-io/go-nats"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware/uartwrap"
	"github.com/yanagiis/GoTuringCoffee/internal/service/mdns"
***REMOVED***

type Service struct {
	uart    uartwrap.UART
	service string
	port    int
	md      *mdns.MDNS
***REMOVED***

func NewService(serviceName string, port int, uart uartwrap.UART, md *mdns.MDNS***REMOVED*** *Service {
	md.Register(serviceName, port***REMOVED***
	return &Service{
		uart:    uart,
		service: serviceName,
		port:    port,
***REMOVED***
***REMOVED***

func (s *Service***REMOVED*** Run(ctx context.Context, nc *nats.EncodedConn***REMOVED*** (err error***REMOVED*** {
	var ln net.Listener
	var conn net.Conn

LOOP:
	for {
		ln, err = net.ListenTCP("tcp", &net.TCPAddr{
			Port: s.port,
	***REMOVED******REMOVED***

		conn, err = ln.Accept(***REMOVED***
	***REMOVED***
			continue
	***REMOVED***

		if err = s.uart.Open(***REMOVED***; err != nil {
			conn.Close(***REMOVED***
			continue
	***REMOVED***

		wg := sync.WaitGroup{***REMOVED***
		wg.Add(2***REMOVED***
		go func(***REMOVED*** {
			ctxcopy(ctx, conn, s.uart***REMOVED***
			wg.Done(***REMOVED***
	***REMOVED***(***REMOVED***
		go func(***REMOVED*** {
			ctxcopy(ctx, s.uart, conn***REMOVED***
			wg.Done(***REMOVED***
	***REMOVED***(***REMOVED***

		wg.Wait(***REMOVED***
		s.uart.Close(***REMOVED***
		conn.Close(***REMOVED***

		select {
		case <-ctx.Done(***REMOVED***:
			break LOOP
		default:
	***REMOVED***
***REMOVED***
	return nil
***REMOVED***

type readerFunc func(p []byte***REMOVED*** (int, error***REMOVED***

func (r readerFunc***REMOVED*** Read(p []byte***REMOVED*** (int, error***REMOVED*** {
	return r(p***REMOVED***
***REMOVED***

func ctxcopy(ctx context.Context, writer io.Writer, reader io.Reader***REMOVED*** error {
	_, err := io.Copy(writer, readerFunc(func(p []byte***REMOVED*** (int, error***REMOVED*** {
		select {
		case <-ctx.Done(***REMOVED***:
			return 0, ctx.Err(***REMOVED***
		default:
			return reader.Read(p***REMOVED***
	***REMOVED***
***REMOVED******REMOVED******REMOVED***
	return err
***REMOVED***
