package hardware

import (
	"errors"
	"fmt"

	"GoTuringCoffee/internal/hardware/max31856"
	"GoTuringCoffee/internal/hardware/max31865"
	"GoTuringCoffee/internal/hardware/spiwrap"
	"GoTuringCoffee/internal/hardware/uartwrap"
	"GoTuringCoffee/internal/service/mdns"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	rpio "github.com/stianeikeland/go-rpio"
	"periph.io/x/periph/host"
)

var (
	ErrDevNotFound = errors.New("device not found")
	ErrWrongConfig = errors.New("wrong configuration")
)

var hardwareFuncs = map[string]ParseHardwareFunc{
	"spi":            ParseSPI,
	"uart":           ParseUART,
	"tcpuartclient":  ParseTcpUartClient,
	"tcpuartserver":  ParseTcpUartServer,
	"pwm":            ParsePWM,
	"max31856":       ParseMAX31856,
	"max31865":       ParseMAX31865,
	"smoothie":       ParseSmoothie,
	"extruder":       ParseExtruder,
	"water_detector": ParseWaterDetector,
}

type ParseHardwareFunc func(m *HWManager, viper *viper.Viper, md *mdns.MDNS) (interface{}, error)

type HardwareType byte

type HWManager struct {
	hws map[string]interface{}
}

type HardwareError struct {
	name string
	msg  string
	err  error
}

func NewHardwareError(name string, msg string, err error) *HardwareError {
	return &HardwareError{
		name: name,
		msg:  msg,
		err:  err,
	}
}

func NewMissingFieldError(name string, field string) *HardwareError {
	return NewHardwareError(name, fmt.Sprintf("%q is required"), ErrWrongConfig)
}

func (e *HardwareError) Error() string {
	return fmt.Sprintf("%s: %s: %s", e.err.Error(), e.name, e.msg)
}

func Init() {
	rpio.Open()
	_, err := host.Init()
	if err != nil {
		log.Fatal().Msgf("failed to initialize periph: %v", err)
	}
}

func NewHWManager() *HWManager {
	return &HWManager{
		hws: make(map[string]interface{}),
	}
}

func (m *HWManager) Load(viper *viper.Viper, md *mdns.MDNS) error {
	if !viper.IsSet("hardwares") {
		return NewHardwareError("hardware", "no hardware field", ErrWrongConfig)
	}

	hardwares := viper.GetStringMapString("hardwares")
	names := make([]string, 0, len(hardwares))
	for name := range hardwares {
		names = append(names, name)
	}

	for {
		hasNewDevice := false
		unresolved := make([]string, 0, len(names))
		for _, name := range names {
			hardware := viper.Sub(fmt.Sprintf("hardwares.%s", name))
			if err := m.AddDevice(name, hardware, md); err != nil {
				herr := err.(*HardwareError)
				switch herr.err {
				case ErrDevNotFound:
					unresolved = append(unresolved, name)
					continue
				default:
					log.Fatal().Msgf("load hardware %q failed - %v", name, err)
					return err
				}
			}
			log.Info().Msgf("load hardware %q", name)
			hasNewDevice = true
		}
		names = unresolved
		if len(names) > 0 && !hasNewDevice {
			log.Warn().Msgf("%v are not loaded", names)
			break
		}
		if len(names) == 0 {
			break
		}
	}

	return nil
}

func (m *HWManager) AddDevice(name string, viper *viper.Viper, md *mdns.MDNS) error {
	var hw interface{}

	if _, ok := m.hws[name]; ok {
		return NewHardwareError(name, "name is used", ErrWrongConfig)
	}
	if !viper.IsSet("type") {
		return NewHardwareError(name, "miss type field", ErrWrongConfig)
	}

	t := viper.GetString("type")
	fn, ok := hardwareFuncs[t]
	if !ok {
		return NewHardwareError(name, fmt.Sprintf("%s is not support yet", t), ErrWrongConfig)
	}

	hw, err := fn(m, viper, md)
	if err != nil {
		return err
	}

	m.hws[name] = hw
	return nil
}

func (m *HWManager) GetDevice(name string) (interface{}, error) {
	if device, ok := m.hws[name]; ok {
		return device, nil
	}
	return nil, fmt.Errorf("Cannot find %s device", name)
}

func ParseSPI(m *HWManager, viper *viper.Viper, md *mdns.MDNS) (interface{}, error) {
	var spi spiwrap.SPIDevice
	if err := viper.Unmarshal(&spi.Conf); err != nil {
		return nil, NewHardwareError("SPI", err.Error(), ErrWrongConfig)
	}
	return &spi, nil
}

func ParseUART(m *HWManager, viper *viper.Viper, md *mdns.MDNS) (interface{}, error) {
	var uart uartwrap.UARTDevice
	if err := viper.Unmarshal(&uart.Conf); err != nil {
		return nil, NewHardwareError("UART", err.Error(), ErrWrongConfig)
	}
	return &uart, nil
}

func ParseTcpUartClient(m *HWManager, viper *viper.Viper, md *mdns.MDNS) (interface{}, error) {
	fields := []string{"mdns_service"}
	if err := checkFields(viper, fields); err != nil {
		return nil, err
	}
	service := viper.GetString("mdns_service")
	client := uartwrap.NewTCPUARTClientMDNS(service, md)
	return client, nil
}

func ParseTcpUartServer(m *HWManager, viper *viper.Viper, md *mdns.MDNS) (interface{}, error) {
	var uartdev interface{}
	var err error
	fields := []string{"service", "port", "uartdev"}
	if err = checkFields(viper, fields); err != nil {
		return nil, err
	}
	service := viper.GetString("service")
	port := viper.GetInt("port")
	if uartdev, err = m.GetDevice(viper.GetString("uartdev")); err != nil {
		return nil, NewHardwareError("tcpuarts", err.Error(), ErrDevNotFound)
	}
	server := uartwrap.NewTCPUARTServerMDNS(service, port, uartdev.(uartwrap.UART), md)
	return server, nil
}

func ParsePWM(m *HWManager, viper *viper.Viper, md *mdns.MDNS) (interface{}, error) {
	var pwm PWMDevice
	if err := viper.Unmarshal(&pwm); err != nil {
		return nil, err
	}
	return &pwm, nil
}

func ParseMAX31856(m *HWManager, viper *viper.Viper, md *mdns.MDNS) (interface{}, error) {
	var device interface{}
	var err error
	var tc max31856.Type
	var mode max31856.Mode
	var sample max31856.Sample

	fields := []string{"dev", "tc", "mode", "sample"}
	if err = checkFields(viper, fields); err != nil {
		return nil, err
	}

	if device, err = m.GetDevice(viper.GetString("dev")); err != nil {
		return nil, NewHardwareError("max31856", err.Error(), ErrDevNotFound)
	}

	if tc, err = max31856.ParseType(viper.GetString("tc")); err != nil {
		return nil, err
	}
	if mode, err = max31856.ParseMode(viper.GetString("mode")); err != nil {
		return nil, err
	}
	if sample, err = max31856.ParseSample(viper.GetString("sample")); err != nil {
		return nil, err
	}

	sensor := max31856.New(
		device.(spiwrap.SPI),
		max31856.Config{
			TC:   tc,
			Avg:  sample,
			Mode: mode,
		},
	)

	return sensor, nil
}

func ParseMAX31865(m *HWManager, viper *viper.Viper, md *mdns.MDNS) (interface{}, error) {
	var device interface{}
	var err error
	var wire max31865.Wire
	var mode max31865.Mode

	fields := []string{"dev", "wire", "mode"}
	if err = checkFields(viper, fields); err != nil {
		return nil, err
	}
	if device, err = m.GetDevice(viper.GetString("dev")); err != nil {
		return nil, NewHardwareError("max31865", err.Error(), ErrDevNotFound)
	}
	if wire, err = max31865.ParseWire(viper.GetString("wire")); err != nil {
		return nil, err
	}
	if mode, err = max31865.ParseMode(viper.GetString("mode")); err != nil {
		return nil, err
	}

	sensor := max31865.New(
		device.(spiwrap.SPI),
		max31865.Config{
			Wire: wire,
			Mode: mode,
		},
	)

	return sensor, nil
}

func ParseSmoothie(m *HWManager, viper *viper.Viper, md *mdns.MDNS) (interface{}, error) {
	var err error
	var device interface{}
	if err = checkFields(viper, []string{"dev"}); err != nil {
		return nil, err
	}
	if device, err = m.GetDevice(viper.GetString("dev")); err != nil {
		return nil, NewHardwareError("smoothie", err.Error(), ErrDevNotFound)
	}
	smoothie := NewSmoothie(device.(SmoothiePort))
	return smoothie, nil
}

func ParseExtruder(m *HWManager, viper *viper.Viper, md *mdns.MDNS) (interface{}, error) {
	var err error
	var device interface{}
	if err = checkFields(viper, []string{"dev"}); err != nil {
		return nil, err
	}
	if device, err = m.GetDevice(viper.GetString("dev")); err != nil {
		return nil, NewHardwareError("extruder", err.Error(), ErrDevNotFound)
	}
	extruder := NewExtruder(device.(ExtruderPort))
	return extruder, nil
}

func ParseWaterDetector(m *HWManager, viper *viper.Viper, md *mdns.MDNS) (interface{}, error) {
	var pin uint32
	var err error
	if err = checkFields(viper, []string{"gpio"}); err != nil {
		return nil, err
	}
	pin = uint32(viper.GetInt("gpio"))
	gpio := &GPIOWaterDetector{
		Pin: pin,
	}
	return gpio, nil
}

func checkFields(viper *viper.Viper, fields []string) error {
	for _, field := range fields {
		if !viper.IsSet(field) {
			return NewMissingFieldError("", field)
		}
	}
	return nil
}
