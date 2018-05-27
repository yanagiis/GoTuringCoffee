***REMOVED***

***REMOVED***
***REMOVED***
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/service"
	"github.com/yanagiis/GoTuringCoffee/internal/service/mdns"
***REMOVED***

type NatsConfig struct {
	Register bool   `mapstructure:"register"`
	Service  string `mapstructure:"service"`
	Port     int    `mapstructure:"port"`
***REMOVED***

func init(***REMOVED*** {
	hardware.Init(***REMOVED***
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr***REMOVED******REMOVED***
***REMOVED***

func main(***REMOVED*** {
	var err error
	var configFile string
	var natsConf NatsConfig
	var mdnsConf mdns.Config
	var m *mdns.MDNS

	flag.StringVar(&configFile, "config", "config", "configuration file"***REMOVED***
	flag.Parse(***REMOVED***

	viper.SetConfigName(configFile***REMOVED***
	viper.AddConfigPath("."***REMOVED***
	if err := viper.ReadInConfig(***REMOVED***; err != nil {
		log.Fatal(***REMOVED***.Msg(err.Error(***REMOVED******REMOVED***
***REMOVED***

	if err = viper.UnmarshalKey("mdns", &mdnsConf***REMOVED***; err != nil {
		log.Fatal(***REMOVED***.Msg(err.Error(***REMOVED******REMOVED***
***REMOVED***
	if err = viper.UnmarshalKey("nats", &natsConf***REMOVED***; err != nil {
		log.Fatal(***REMOVED***.Msg(err.Error(***REMOVED******REMOVED***
***REMOVED***

	m = mdns.NewMDNS(mdnsConf***REMOVED***
	if natsConf.Register {
		if err = m.Register(natsConf.Service, natsConf.Port***REMOVED***; err != nil {
			log.Fatal(***REMOVED***.Msg(err.Error(***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

	log.Info(***REMOVED***.Msg("Load hardware configurations ..."***REMOVED***
	hwm := hardware.NewHWManager(***REMOVED***
	err = hwm.Load(viper.GetViper(***REMOVED***, m***REMOVED***
***REMOVED***
		log.Fatal(***REMOVED***.Err(err***REMOVED***
		panic(err***REMOVED***
***REMOVED***
	log.Info(***REMOVED***.Msg("Load hardware configurations successfully"***REMOVED***

	log.Info(***REMOVED***.Msg("Load services configurations ..."***REMOVED***
	svm := service.NewServiceManager(***REMOVED***
	err = svm.Load(viper.GetViper(***REMOVED***, hwm, m***REMOVED***
***REMOVED***
		log.Fatal(***REMOVED***.Err(err***REMOVED***
		panic(err***REMOVED***
***REMOVED***
	log.Info(***REMOVED***.Msg("Load services configurations successfully"***REMOVED***

	conn := connectNats(natsConf, m***REMOVED***
	defer conn.Close(***REMOVED***

	log.Info(***REMOVED***.Msg("Run services ..."***REMOVED***
	svm.RunServices(conn***REMOVED***
	log.Info(***REMOVED***.Msg("Run services successfully"***REMOVED***
	defer func(***REMOVED*** {
		log.Info(***REMOVED***.Msg("Stop services ..."***REMOVED***
		svm.StopServices(***REMOVED***
		log.Info(***REMOVED***.Msg("Stop services successfully"***REMOVED***
***REMOVED***(***REMOVED***

	sigs := make(chan os.Signal, 1***REMOVED***
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM***REMOVED***
	<-sigs
***REMOVED***

func connectNats(conf NatsConfig, m *mdns.MDNS***REMOVED*** *nats.EncodedConn {
	var addrs []net.IP
	var port int
	var nc *nats.Conn
	var err error

	for {
		if addrs, port, err = m.Lookup(conf.Service, 5*time.Second***REMOVED***; err != nil {
			log.Info(***REMOVED***.Msgf(err.Error(***REMOVED******REMOVED***
			continue
	***REMOVED***
		if len(addrs***REMOVED*** == 0 {
			log.Info(***REMOVED***.Msgf("Cannot find %q address", conf.Service***REMOVED***
			continue
	***REMOVED***

		log.Info(***REMOVED***.Msg("Connect to Nats server ..."***REMOVED***
		url := fmt.Sprintf("nats://%s:%d", addrs[0], port***REMOVED***
		opt := nats.Option(func(opts *nats.Options***REMOVED*** error {
			opts.ReconnectedCB = func(conn *nats.Conn***REMOVED*** {
				log.Info(***REMOVED***.Msg("Reconnect to Nats server successfully"***REMOVED***
		***REMOVED***
			opts.DisconnectedCB = func(conn *nats.Conn***REMOVED*** {
				log.Info(***REMOVED***.Msg("Nats server is disconnected"***REMOVED***
		***REMOVED***
			return nil
	***REMOVED******REMOVED***
		if nc, err = nats.Connect(url, opt***REMOVED***; err != nil {
			log.Info(***REMOVED***.Err(err***REMOVED***
			continue
	***REMOVED***
		conn, _ := nats.NewEncodedConn(nc, "jsoniter"***REMOVED***
		log.Info(***REMOVED***.Msg("Connect to Nats server successfully"***REMOVED***
		return conn
***REMOVED***
***REMOVED***
