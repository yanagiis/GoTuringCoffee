package hardware

***REMOVED***
	"errors"
***REMOVED***

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware/max31856"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware/max31865"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware/spiwrap"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware/uartwrap"
	"github.com/yanagiis/GoTuringCoffee/internal/service/mdns"
	"github.com/yanagiis/periph/host"
***REMOVED***

var (
	ErrDevNotFound = errors.New("device not found"***REMOVED***
	ErrWrongConfig = errors.New("wrong configuration"***REMOVED***
***REMOVED***

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
***REMOVED***

type ParseHardwareFunc func(m *HWManager, viper *viper.Viper, md *mdns.MDNS***REMOVED*** (interface{***REMOVED***, error***REMOVED***

type HardwareType byte

type HWManager struct {
	hws map[string]interface{***REMOVED***
***REMOVED***

type HardwareError struct {
	name string
	msg  string
	err  error
***REMOVED***

func NewHardwareError(name string, msg string, err error***REMOVED*** *HardwareError {
	return &HardwareError{
		name: name,
		msg:  msg,
		err:  err,
***REMOVED***
***REMOVED***

func NewMissingFieldError(name string, field string***REMOVED*** *HardwareError {
	return NewHardwareError(name, fmt.Sprintf("%q is required"***REMOVED***, ErrWrongConfig***REMOVED***
***REMOVED***

func (e *HardwareError***REMOVED*** Error(***REMOVED*** string {
	return fmt.Sprintf("%s: %s: %s", e.err.Error(***REMOVED***, e.name, e.msg***REMOVED***
***REMOVED***

func Init(***REMOVED*** {
	_, err := host.Init(***REMOVED***
***REMOVED***
		log.Fatal(***REMOVED***.Msgf("failed to initialize periph: %v", err***REMOVED***
***REMOVED***
***REMOVED***

func NewHWManager(***REMOVED*** *HWManager {
	return &HWManager{
		hws: make(map[string]interface{***REMOVED******REMOVED***,
***REMOVED***
***REMOVED***

func (m *HWManager***REMOVED*** Load(viper *viper.Viper, md *mdns.MDNS***REMOVED*** error {
	if !viper.IsSet("hardwares"***REMOVED*** {
		return NewHardwareError("hardware", "no hardware field", ErrWrongConfig***REMOVED***
***REMOVED***

	hardwares := viper.GetStringMapString("hardwares"***REMOVED***
	names := make([]string, 0, len(hardwares***REMOVED******REMOVED***
	for name := range hardwares {
		names = append(names, name***REMOVED***
***REMOVED***

	for {
		hasNewDevice := false
		unresolved := make([]string, 0, len(names***REMOVED******REMOVED***
		for _, name := range names {
			hardware := viper.Sub(fmt.Sprintf("hardwares.%s", name***REMOVED******REMOVED***
			if err := m.AddDevice(name, hardware, md***REMOVED***; err != nil {
				herr := err.(*HardwareError***REMOVED***
				switch herr.err {
				case ErrDevNotFound:
					unresolved = append(unresolved, name***REMOVED***
					continue
				default:
					log.Fatal(***REMOVED***.Msgf("load hardware %q failed - %v", name, err***REMOVED***
					return err
			***REMOVED***
		***REMOVED***
			log.Info(***REMOVED***.Msgf("load hardware %q", name***REMOVED***
			hasNewDevice = true
	***REMOVED***
		names = unresolved
		if len(names***REMOVED*** > 0 && !hasNewDevice {
			log.Warn(***REMOVED***.Msgf("%v are not loaded", names***REMOVED***
			break
	***REMOVED***
		if len(names***REMOVED*** == 0 {
			break
	***REMOVED***
***REMOVED***

	return nil
***REMOVED***

func (m *HWManager***REMOVED*** AddDevice(name string, viper *viper.Viper, md *mdns.MDNS***REMOVED*** error {
	var hw interface{***REMOVED***

	if _, ok := m.hws[name]; ok {
		return NewHardwareError(name, "name is used", ErrWrongConfig***REMOVED***
***REMOVED***
	if !viper.IsSet("type"***REMOVED*** {
		return NewHardwareError(name, "miss type field", ErrWrongConfig***REMOVED***
***REMOVED***

	t := viper.GetString("type"***REMOVED***
	fn, ok := hardwareFuncs[t]
	if !ok {
		return NewHardwareError(name, fmt.Sprintf("%s is not support yet", t***REMOVED***, ErrWrongConfig***REMOVED***
***REMOVED***

	hw, err := fn(m, viper, md***REMOVED***
***REMOVED***
		return err
***REMOVED***

	m.hws[name] = hw
	return nil
***REMOVED***

func (m *HWManager***REMOVED*** GetDevice(name string***REMOVED*** (interface{***REMOVED***, error***REMOVED*** {
	if device, ok := m.hws[name]; ok {
		return device, nil
***REMOVED***
	return nil, fmt.Errorf("Cannot find %s device", name***REMOVED***
***REMOVED***

func ParseSPI(m *HWManager, viper *viper.Viper, md *mdns.MDNS***REMOVED*** (interface{***REMOVED***, error***REMOVED*** {
	var spi spiwrap.SPIDevice
	if err := viper.Unmarshal(&spi.Conf***REMOVED***; err != nil {
		return nil, NewHardwareError("SPI", err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
***REMOVED***
	return &spi, nil
***REMOVED***

func ParseUART(m *HWManager, viper *viper.Viper, md *mdns.MDNS***REMOVED*** (interface{***REMOVED***, error***REMOVED*** {
	var uart uartwrap.UARTDevice
	if err := viper.Unmarshal(&uart.Conf***REMOVED***; err != nil {
		return nil, NewHardwareError("UART", err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
***REMOVED***
	return &uart, nil
***REMOVED***

func ParseTcpUartClient(m *HWManager, viper *viper.Viper, md *mdns.MDNS***REMOVED*** (interface{***REMOVED***, error***REMOVED*** {
	fields := []string{"mdns_service"***REMOVED***
	if err := checkFields(viper, fields***REMOVED***; err != nil {
		return nil, err
***REMOVED***
	service := viper.GetString("mdns_service"***REMOVED***
	client := uartwrap.NewTCPUARTClientMDNS(service, md***REMOVED***
	return client, nil
***REMOVED***

func ParseTcpUartServer(m *HWManager, viper *viper.Viper, md *mdns.MDNS***REMOVED*** (interface{***REMOVED***, error***REMOVED*** {
	var uartdev interface{***REMOVED***
	var err error
	fields := []string{"service", "port", "uartdev"***REMOVED***
	if err = checkFields(viper, fields***REMOVED***; err != nil {
		return nil, err
***REMOVED***
	service := viper.GetString("service"***REMOVED***
	port := viper.GetInt("port"***REMOVED***
	if uartdev, err = m.GetDevice(viper.GetString("uartdev"***REMOVED******REMOVED***; err != nil {
		return nil, NewHardwareError("tcpuarts", err.Error(***REMOVED***, ErrDevNotFound***REMOVED***
***REMOVED***
	server := uartwrap.NewTCPUARTServerMDNS(service, port, uartdev.(uartwrap.UART***REMOVED***, md***REMOVED***
	return server, nil
***REMOVED***

func ParsePWM(m *HWManager, viper *viper.Viper, md *mdns.MDNS***REMOVED*** (interface{***REMOVED***, error***REMOVED*** {
	var pwm PWMDevice
	if err := viper.Unmarshal(&pwm***REMOVED***; err != nil {
		return nil, err
***REMOVED***
	return &pwm, nil
***REMOVED***

func ParseMAX31856(m *HWManager, viper *viper.Viper, md *mdns.MDNS***REMOVED*** (interface{***REMOVED***, error***REMOVED*** {
	var device interface{***REMOVED***
	var err error
	var tc max31856.Type
	var mode max31856.Mode
	var sample max31856.Sample

	fields := []string{"dev", "tc", "mode", "sample"***REMOVED***
	if err = checkFields(viper, fields***REMOVED***; err != nil {
		return nil, err
***REMOVED***

	if device, err = m.GetDevice(viper.GetString("dev"***REMOVED******REMOVED***; err != nil {
		return nil, NewHardwareError("max31856", err.Error(***REMOVED***, ErrDevNotFound***REMOVED***
***REMOVED***

	if tc, err = max31856.ParseType(viper.GetString("tc"***REMOVED******REMOVED***; err != nil {
		return nil, err
***REMOVED***
	if mode, err = max31856.ParseMode(viper.GetString("mode"***REMOVED******REMOVED***; err != nil {
		return nil, err
***REMOVED***
	if sample, err = max31856.ParseSample(viper.GetString("sample"***REMOVED******REMOVED***; err != nil {
		return nil, err
***REMOVED***

	sensor := max31856.New(
		device.(spiwrap.SPI***REMOVED***,
		max31856.Config{
			TC:   tc,
			Avg:  sample,
			Mode: mode,
	***REMOVED***,
	***REMOVED***

	return sensor, nil
***REMOVED***

func ParseMAX31865(m *HWManager, viper *viper.Viper, md *mdns.MDNS***REMOVED*** (interface{***REMOVED***, error***REMOVED*** {
	var device interface{***REMOVED***
	var err error
	var wire max31865.Wire
	var mode max31865.Mode

	fields := []string{"dev", "wire", "mode"***REMOVED***
	if err = checkFields(viper, fields***REMOVED***; err != nil {
		return nil, err
***REMOVED***
	if device, err = m.GetDevice(viper.GetString("dev"***REMOVED******REMOVED***; err != nil {
		return nil, NewHardwareError("max31865", err.Error(***REMOVED***, ErrDevNotFound***REMOVED***
***REMOVED***
	if wire, err = max31865.ParseWire(viper.GetString("wire"***REMOVED******REMOVED***; err != nil {
		return nil, err
***REMOVED***
	if mode, err = max31865.ParseMode(viper.GetString("mode"***REMOVED******REMOVED***; err != nil {
		return nil, err
***REMOVED***

	sensor := max31865.New(
		device.(spiwrap.SPI***REMOVED***,
		max31865.Config{
			Wire: wire,
			Mode: mode,
	***REMOVED***,
	***REMOVED***

	return sensor, nil
***REMOVED***

func ParseSmoothie(m *HWManager, viper *viper.Viper, md *mdns.MDNS***REMOVED*** (interface{***REMOVED***, error***REMOVED*** {
	var err error
	var device interface{***REMOVED***
	if err = checkFields(viper, []string{"dev"***REMOVED******REMOVED***; err != nil {
		return nil, err
***REMOVED***
	if device, err = m.GetDevice(viper.GetString("dev"***REMOVED******REMOVED***; err != nil {
		return nil, NewHardwareError("smoothie", err.Error(***REMOVED***, ErrDevNotFound***REMOVED***
***REMOVED***
	smoothie := NewSmoothie(device.(SmoothiePort***REMOVED******REMOVED***
	return smoothie, nil
***REMOVED***

func ParseExtruder(m *HWManager, viper *viper.Viper, md *mdns.MDNS***REMOVED*** (interface{***REMOVED***, error***REMOVED*** {
	var err error
	var device interface{***REMOVED***
	if err = checkFields(viper, []string{"dev"***REMOVED******REMOVED***; err != nil {
		return nil, err
***REMOVED***
	if device, err = m.GetDevice(viper.GetString("dev"***REMOVED******REMOVED***; err != nil {
		return nil, NewHardwareError("extruder", err.Error(***REMOVED***, ErrDevNotFound***REMOVED***
***REMOVED***
	extruder := NewExtruder(device.(ExtruderPort***REMOVED******REMOVED***
	return extruder, nil
***REMOVED***

func ParseWaterDetector(m *HWManager, viper *viper.Viper, md *mdns.MDNS***REMOVED*** (interface{***REMOVED***, error***REMOVED*** {
	var pin uint32
	var err error
	if err = checkFields(viper, []string{"gpio"***REMOVED******REMOVED***; err != nil {
		return nil, err
***REMOVED***
	pin = uint32(viper.GetInt("gpio"***REMOVED******REMOVED***
	gpio := &GPIOWaterDetector{
		Pin: pin,
***REMOVED***
	return gpio, nil
***REMOVED***

func checkFields(viper *viper.Viper, fields []string***REMOVED*** error {
	for _, field := range fields {
		if !viper.IsSet(field***REMOVED*** {
			return NewMissingFieldError("", field***REMOVED***
	***REMOVED***
***REMOVED***
	return nil
***REMOVED***
