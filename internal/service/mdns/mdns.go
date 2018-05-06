package mdns

***REMOVED***
	"context"
***REMOVED***
	"net"
	"time"

	"github.com/grandcat/zeroconf"
***REMOVED***

type Request struct {
	serviceName string
	reply       chan Response
	timeout     time.Duration
***REMOVED***

type Response struct {
	IPv4 []net.IP
	Port int
	Err  error
***REMOVED***

type Config struct {
	Instance string `mapstructure:"instance"`
	Domain   string `mapstructure:"domain"`
***REMOVED***

type MDNS struct {
	Conf     Config
	services map[string]*zeroconf.Server
	resolver *zeroconf.Resolver
***REMOVED***

func NewMDNS(conf Config***REMOVED*** *MDNS {
	return &MDNS{
		Conf:     conf,
		***REMOVED*** make(map[string]*zeroconf.Server***REMOVED***,
***REMOVED***
***REMOVED***

func (m *MDNS***REMOVED*** Register(service string, port int***REMOVED*** (err error***REMOVED*** {
	var server *zeroconf.Server
	if server, err = zeroconf.Register(m.Conf.Instance, service, m.Conf.Domain, port, nil, nil***REMOVED***; err != nil {
		return
***REMOVED***
	m.services[service] = server
	return
***REMOVED***

func (m *MDNS***REMOVED*** Stop(***REMOVED*** {
	for name, service := range m.services {
		service.Shutdown(***REMOVED***
		delete(m.services, name***REMOVED***
***REMOVED***
***REMOVED***

func (m *MDNS***REMOVED*** Lookup(serviceName string, timeout time.Duration***REMOVED*** ([]net.IP, int, error***REMOVED*** {
	resolver, err := zeroconf.NewResolver(nil***REMOVED***
***REMOVED***
		return nil, 0, err
***REMOVED***

	ctx, cancel := context.WithTimeout(context.Background(***REMOVED***, timeout***REMOVED***
	defer cancel(***REMOVED***

	entries := make(chan *zeroconf.ServiceEntry***REMOVED***
	if err = resolver.Browse(ctx, serviceName, "local.", entries***REMOVED***; err != nil {
		return nil, 0, err
***REMOVED***
	<-ctx.Done(***REMOVED***

	for entry := range entries {
		return entry.AddrIPv4, entry.Port, nil
***REMOVED***

	return nil, 0, fmt.Errorf("Cannot find %q address", serviceName***REMOVED***
***REMOVED***
