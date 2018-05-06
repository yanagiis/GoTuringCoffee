package uartwrap

***REMOVED***
	"context"
***REMOVED***
	"net"
	"time"

	"github.com/yanagiis/GoTuringCoffee/internal/service/mdns"
***REMOVED***

type TCPUARTClient struct {
	service string
	md      *mdns.MDNS
	conn    net.Conn
***REMOVED***

func NewTCPUARTClientMDNS(service string, md *mdns.MDNS***REMOVED*** *TCPUARTClient {
	return &TCPUARTClient{
		service: service,
		md:      md,
***REMOVED***
***REMOVED***

func (c *TCPUARTClient***REMOVED*** Open(***REMOVED*** (err error***REMOVED*** {
	var addrs []net.IP
	var port int
	if addrs, port, err = c.md.Lookup(c.service, time.Second***REMOVED***; err != nil {
		return
***REMOVED***
	if len(addrs***REMOVED*** == 0 {
		err = fmt.Errorf("Cannot lookup %q ip and port", c.service***REMOVED***
		return
***REMOVED***

	url := fmt.Sprintf("%s:%d", addrs[0], port***REMOVED***
	if c.conn, err = net.Dial("tcp", url***REMOVED***; err != nil {
		return
***REMOVED***
	return
***REMOVED***

func (c *TCPUARTClient***REMOVED*** IsOpen(***REMOVED*** bool {
	return c.conn != nil
***REMOVED***

func (c *TCPUARTClient***REMOVED*** Close(***REMOVED*** (err error***REMOVED*** {
	if err = c.conn.Close(***REMOVED***; err != nil {
		return
***REMOVED***
	c.conn = nil
	return
***REMOVED***

func (c *TCPUARTClient***REMOVED*** Read(p []byte***REMOVED*** (int, error***REMOVED*** {
	return c.conn.Read(p***REMOVED***
***REMOVED***

func (c *TCPUARTClient***REMOVED*** Write(p []byte***REMOVED*** (int, error***REMOVED*** {
	return c.conn.Write(p***REMOVED***
***REMOVED***

type TCPUARTServer struct {
	Service string
	Port    int
	md      *mdns.MDNS
	conn    net.Conn
	uart    UART
	ctx     context.Context
***REMOVED***

func NewTCPUARTServerMDNS(service string, port int, uart UART, md *mdns.MDNS***REMOVED*** *TCPUARTServer {
	server := &TCPUARTServer{
		Service: service,
		md:      md,
		uart:    uart,
***REMOVED***
	if err := md.Register(service, port***REMOVED***; err != nil {
		return nil
***REMOVED***
	return server
***REMOVED***

func (s *TCPUARTServer***REMOVED*** Pair(timeout time.Duration***REMOVED*** (conn *net.Conn, err error***REMOVED*** {
	var ln *net.TCPListener

	if err = s.uart.Open(***REMOVED***; err != nil {
		return
***REMOVED***
	tcpAddr := net.TCPAddr{
		Port: s.Port,
***REMOVED***
	if ln, err = net.ListenTCP("tcp", &tcpAddr***REMOVED***; err != nil {
		s.uart.Close(***REMOVED***
		return
***REMOVED***
	defer ln.Close(***REMOVED***
	if s.conn, err = ln.Accept(***REMOVED***; err != nil {
		s.uart.Close(***REMOVED***
		return
***REMOVED***
	return
***REMOVED***

func (s *TCPUARTServer***REMOVED*** Unpair(***REMOVED*** (err error***REMOVED*** {
	s.uart.Close(***REMOVED***
	err = s.conn.Close(***REMOVED***
	return
***REMOVED***
