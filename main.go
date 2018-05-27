package main

import (
	"fmt"
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
)

type NatsConfig struct {
	Register bool   `mapstructure:"register"`
	Service  string `mapstructure:"service"`
	Port     int    `mapstructure:"port"`
}

func init() {
	hardware.Init()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	var err error
	var configFile string
	var natsConf NatsConfig
	var mdnsConf mdns.Config
	var m *mdns.MDNS

	flag.StringVar(&configFile, "config", "config", "configuration file")
	flag.Parse()

	viper.SetConfigName(configFile)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Msg(err.Error())
	}

	if err = viper.UnmarshalKey("mdns", &mdnsConf); err != nil {
		log.Fatal().Msg(err.Error())
	}
	if err = viper.UnmarshalKey("nats", &natsConf); err != nil {
		log.Fatal().Msg(err.Error())
	}

	m = mdns.NewMDNS(mdnsConf)
	if natsConf.Register {
		if err = m.Register(natsConf.Service, natsConf.Port); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}

	log.Info().Msg("Load hardware configurations ...")
	hwm := hardware.NewHWManager()
	err = hwm.Load(viper.GetViper(), m)
	if err != nil {
		log.Fatal().Err(err)
		panic(err)
	}
	log.Info().Msg("Load hardware configurations successfully")

	log.Info().Msg("Load services configurations ...")
	svm := service.NewServiceManager()
	err = svm.Load(viper.GetViper(), hwm, m)
	if err != nil {
		log.Fatal().Err(err)
		panic(err)
	}
	log.Info().Msg("Load services configurations successfully")

	conn := connectNats(natsConf, m)
	defer conn.Close()

	log.Info().Msg("Run services ...")
	svm.RunServices(conn)
	log.Info().Msg("Run services successfully")
	defer func() {
		log.Info().Msg("Stop services ...")
		svm.StopServices()
		log.Info().Msg("Stop services successfully")
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func connectNats(conf NatsConfig, m *mdns.MDNS) *nats.EncodedConn {
	var addrs []net.IP
	var port int
	var nc *nats.Conn
	var err error

	for {
		if addrs, port, err = m.Lookup(conf.Service, 5*time.Second); err != nil {
			log.Info().Msgf(err.Error())
			continue
		}
		if len(addrs) == 0 {
			log.Info().Msgf("Cannot find %q address", conf.Service)
			continue
		}

		log.Info().Msg("Connect to Nats server ...")
		url := fmt.Sprintf("nats://%s:%d", addrs[0], port)
		opt := nats.Option(func(opts *nats.Options) error {
			opts.ReconnectedCB = func(conn *nats.Conn) {
				log.Info().Msg("Reconnect to Nats server successfully")
			}
			opts.DisconnectedCB = func(conn *nats.Conn) {
				log.Info().Msg("Nats server is disconnected")
			}
			return nil
		})
		if nc, err = nats.Connect(url, opt); err != nil {
			log.Info().Err(err)
			continue
		}
		conn, _ := nats.NewEncodedConn(nc, "jsoniter")
		log.Info().Msg("Connect to Nats server successfully")
		return conn
	}
}
