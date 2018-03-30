***REMOVED***

***REMOVED***
***REMOVED***
	"time"

	nats "github.com/nats-io/go-nats"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/service"
***REMOVED***

type NatsConfig struct {
	Addr string `mapstructure:"host"`
	Port uint32 `mapstructure:"port"`
***REMOVED***

func main(***REMOVED*** {
	var nc *nats.Conn
	var err error
	var configFile string
	flag.StringVar(&configFile, "config", "config", "configuration file"***REMOVED***
	flag.Parse(***REMOVED***

	viper.SetConfigName(configFile***REMOVED***
	viper.AddConfigPath("."***REMOVED***
	if err := viper.ReadInConfig(***REMOVED***; err != nil {
		panic(err***REMOVED***
***REMOVED***

	var natsConf NatsConfig
	if err = viper.UnmarshalKey("nats", &natsConf***REMOVED***; err != nil {
		panic(err***REMOVED***
***REMOVED***

	addr := fmt.Sprintf("***REMOVED***//%s:%d", natsConf.Addr, natsConf.Port***REMOVED***
	nc, err = nats.Connect(addr***REMOVED***
***REMOVED***
		panic(err***REMOVED***
***REMOVED***
	conn, _ := nats.NewEncodedConn(nc, "jsoniter"***REMOVED***
	defer conn.Close(***REMOVED***

	hwm := hardware.NewHWManager(***REMOVED***
	err = hwm.Load(viper.GetViper(***REMOVED******REMOVED***
***REMOVED***
		panic(err***REMOVED***
***REMOVED***

	svm := service.NewServiceManager(***REMOVED***
	err = svm.Load(viper.GetViper(***REMOVED***, hwm***REMOVED***
***REMOVED***
		panic(err***REMOVED***
***REMOVED***

	svm.RunServices(conn***REMOVED***
	defer svm.StopServices(***REMOVED***

	for {
		time.Sleep(1 * time.Second***REMOVED***
***REMOVED***
***REMOVED***
