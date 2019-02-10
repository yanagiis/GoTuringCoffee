package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"GoTuringCoffee/internal/hardware"
	"GoTuringCoffee/internal/hardware/uartwrap"
	"GoTuringCoffee/internal/service/barista"
	"GoTuringCoffee/internal/service/distance"
	"GoTuringCoffee/internal/service/heater"
	"GoTuringCoffee/internal/service/lib"
	"GoTuringCoffee/internal/service/mdns"
	"GoTuringCoffee/internal/service/outtemp"
	"GoTuringCoffee/internal/service/replenisher"
	"GoTuringCoffee/internal/service/tankmeter"
	"GoTuringCoffee/internal/service/tanktemp"
	"GoTuringCoffee/internal/service/uartserver"
	"GoTuringCoffee/internal/service/web"
	"GoTuringCoffee/internal/service/web/model"

	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	ErrWrongConfig = errors.New("wrong configuration")
)

type services interface {
	Run(ctx context.Context, nc *nats.EncodedConn, fin chan<- struct{}) error
	Stop() error
}

type ServiceError struct {
	name string
	msg  string
	err  error
}

func NewServiceError(name string, msg string, err error) *ServiceError {
	return &ServiceError{
		name: name,
		msg:  msg,
		err:  err,
	}
}

func NewMissingFieldError(name string, field string) *ServiceError {
	return NewServiceError(name, fmt.Sprintf("%q is required"), ErrWrongConfig)
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("%s: %s: %s", e.err.Error(), e.name, e.msg)
}

type ServiceManager struct {
	services map[string]services
	cancels  map[string]context.CancelFunc
	fins     map[string]chan struct{}
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		services: make(map[string]services),
		cancels:  make(map[string]context.CancelFunc),
		fins:     make(map[string]chan struct{}),
	}
}

func (s *ServiceManager) Load(viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS) error {
	if !viper.IsSet("services") {
		return NewServiceError("services", "no services field", ErrWrongConfig)
	}

	services := viper.GetStringMapString("services")
	for name := range services {
		services := viper.Sub(fmt.Sprintf("services.%s", name))
		if err := s.AddService(name, services, hwm, m); err != nil {
			log.Fatal().Msgf("load service %q failed - %v", name, err)
			return err
		}
		log.Info().Msgf("load service %q", name)
	}

	return nil
}

func (s *ServiceManager) AddService(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS) error {
	if _, ok := s.services[name]; ok {
		return NewServiceError(name, "name is used", ErrWrongConfig)
	}
	if !viper.IsSet("enable") {
		return NewServiceError(name, "miss enable field", ErrWrongConfig)
	}

	if !viper.GetBool("enable") {
		log.Debug().Msgf("%s is disable", name)
		return nil
	}

	switch name {
	case "output_temp_service":
		if err := s.parseOutTempService(name, viper, hwm, m); err != nil {
			return err
		}
	case "tank_temp_service":
		if err := s.parseTankTempService(name, viper, hwm, m); err != nil {
			return err
		}
	case "tank_meter_service":
		if err := s.parseTankMeterService(name, viper, hwm, m); err != nil {
			return err
		}
	case "replenish_service":
		if err := s.parseReplenisher(name, viper, hwm, m); err != nil {
			return err
		}
	case "heater":
		if err := s.parseHeater(name, viper, hwm, m); err != nil {
			return err
		}
	case "barista":
		if err := s.parseBarista(name, viper, hwm, m); err != nil {
			return err
		}
	case "uartserver":
		if err := s.parseUARTServer(name, viper, hwm, m); err != nil {
			return err
		}
	case "distance_ranging_service":
		if err := s.parseDistanceRangingServer(name, viper, hwm, m); err != nil {
			return err
		}

	// case "uartclient":
	// 	if err := s.parseUARTClient(name, viper, hwm, m); err != nil {
	// 		return err
	// 	}
	case "web":
		if err := s.parseWeb(name, viper, hwm, m); err != nil {
			return err
		}
	default:
		log.Error().Msgf("Unknown service '%s'", name)
		return fmt.Errorf("Unknown service '%s'", name)
	}

	return nil
}

func (s *ServiceManager) parseBarista(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS) error {
	var conf barista.BaristaConfig
	var err error
	var tmp interface{}

	if err = checkFields(viper, []string{"smoothie", "extruder"}); err != nil {
		return err
	}

	smoothieName := viper.GetString("smoothie")
	if tmp, err = hwm.GetDevice(smoothieName); err != nil {
		return err
	}
	smoothie := tmp.(*hardware.Smoothie)

	extruderName := viper.GetString("extruder")
	if tmp, err = hwm.GetDevice(extruderName); err != nil {
		return err
	}
	extruder := tmp.(*hardware.Extruder)

	if err := viper.Unmarshal(&conf); err != nil {
		return NewServiceError(name, err.Error(), ErrWrongConfig)
	}

	barista := barista.NewBarista(conf, &barista.SEController{
		Smoothie: smoothie,
		Extruder: extruder,
	})
	s.services[name] = barista
	return nil
}

func (s *ServiceManager) parseReplenisher(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS) error {
	var pwmConfig hardware.PWMConfig

	devName := viper.GetString("dev")
	dev, err := hwm.GetDevice(devName)
	if err != nil {
		return NewServiceError(name, err.Error(), ErrWrongConfig)
	}
	pwm := dev.(hardware.PWM)

	scanInterval := time.Duration(viper.GetInt64("scan_interval_ms")) * time.Millisecond
	if err = viper.Unmarshal(&pwmConfig); err != nil {
		return NewServiceError(name, err.Error(), ErrWrongConfig)
	}

	s.services[name] = replenisher.NewService(pwm, scanInterval, pwmConfig)
	return nil
}

func (s *ServiceManager) parseOutTempService(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS) error {
	devName := viper.GetString("dev")
	dev, err := hwm.GetDevice(devName)
	if err != nil {
		return NewServiceError(name, err.Error(), ErrWrongConfig)
	}
	sensor := dev.(hardware.TemperatureSensor)
	scanInterval := time.Duration(viper.GetInt64("scan_interval_ms")) * time.Millisecond
	s.services[name] = outtemp.NewService(sensor, scanInterval)
	return nil
}

func (s *ServiceManager) parseTankTempService(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS) error {
	devName := viper.GetString("dev")
	dev, err := hwm.GetDevice(devName)
	if err != nil {
		return NewServiceError(name, err.Error(), ErrWrongConfig)
	}
	sensor := dev.(hardware.TemperatureSensor)
	scanInterval := time.Duration(viper.GetInt64("scan_interval_ms")) * time.Millisecond
	s.services[name] = tanktemp.NewService(sensor, scanInterval)
	return nil
}

func (s *ServiceManager) parseDistanceRangingServer(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS) error {
	devName := viper.GetString("dev")
	dev, err := hwm.GetDevice(devName)
	if err != nil {
		return NewServiceError(name, err.Error(), ErrWrongConfig)
	}
	sensor := dev.(hardware.DistanceRangingSensor)
	scanInterval := time.Duration(viper.GetInt64("scan_interval_ms")) * time.Millisecond
	s.services[name] = distance.NewService(sensor, scanInterval)
	return nil
}

func (s *ServiceManager) parseTankMeterService(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS) error {
	devName := viper.GetString("dev")
	dev, err := hwm.GetDevice(devName)
	if err != nil {
		return NewServiceError(name, err.Error(), ErrWrongConfig)
	}
	sensor := dev.(hardware.WaterDetector)
	scanInterval := time.Duration(viper.GetInt64("scan_interval_ms")) * time.Millisecond
	s.services[name] = tankmeter.NewService(sensor, scanInterval)
	return nil
}

func (s *ServiceManager) parseHeater(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS) error {
	var pid lib.NormalPID
	devName := viper.GetString("pwm_dev")
	dev, err := hwm.GetDevice(devName)
	if err != nil {
		return NewServiceError(name, err.Error(), ErrWrongConfig)
	}
	pwm := dev.(hardware.PWM)
	scanInterval := time.Duration(viper.GetInt64("scan_interval_ms")) * time.Millisecond
	viper.Sub("pid").Unmarshal(&pid)
	s.services[name] = heater.NewService(pwm, scanInterval, &pid)
	return nil
}

func (s *ServiceManager) parseUARTServer(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS) (err error) {
	var device interface{}
	if err = checkFields(viper, []string{"mdns_service", "port", "uartdev"}); err != nil {
		return
	}
	service := viper.GetString("mdns_service")
	port := viper.GetInt("port")
	if device, err = hwm.GetDevice(viper.GetString("uartdev")); err != nil {
		return
	}

	s.services[name] = uartserver.NewService(service, port, device.(uartwrap.UART), m)
	return
}

func (s *ServiceManager) parseWeb(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS) (err error) {
	if err = checkFields(viper, []string{"port", "static_files", "mongodb"}); err != nil {
		return
	}
	port := viper.GetInt("port")
	staticFiles := viper.GetString("static_files")
	mongodbMap := viper.GetStringMapString("mongodb")

	service := &web.Service{
		DB: model.MongoDBConfig{
			Url: mongodbMap["url"],
		},
		Web: web.WebConfig{
			StaticFilePath: staticFiles,
			Port:           port,
		},
	}

	s.services[name] = service
	return
}

func (s *ServiceManager) RunServices(ctx context.Context, nc *nats.EncodedConn) error {
	for name, service := range s.services {
		var newctx context.Context
		newctx, s.cancels[name] = context.WithCancel(ctx)
		s.fins[name] = make(chan struct{})
		log.Info().Msgf("Run '%s' service", name)
		go service.Run(newctx, nc, s.fins[name])
	}
	return nil
}

func (s *ServiceManager) StopServices() error {
	for _, s := range s.services {
		s.Stop()
	}
	for _, cancel := range s.cancels {
		cancel()
	}
	for _, fin := range s.fins {
		<-fin
	}
	return nil
}

func checkFields(viper *viper.Viper, fields []string) error {
	for _, field := range fields {
		if !viper.IsSet(field) {
			return NewMissingFieldError("", field)
		}
	}
	return nil
}
