package hardware

***REMOVED***
	"errors"
***REMOVED***
	"log"

	"github.com/spf13/viper"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware/max31856"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware/max31865"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware/spiwrap"
	"github.com/yanagiis/periph/host"
***REMOVED***

var (
	ErrDevNotFound = errors.New("device not found"***REMOVED***
	ErrWrongConfig = errors.New("wrong configuration"***REMOVED***
***REMOVED***

type HWManager struct {
	hw map[string]interface{***REMOVED***
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

func (e *HardwareError***REMOVED*** Error(***REMOVED*** string {
	return fmt.Sprintf("%s: %s: %s", e.err.Error(***REMOVED***, e.name, e.msg***REMOVED***
***REMOVED***

func Init(***REMOVED*** {
	_, err := host.Init(***REMOVED***
***REMOVED***
		log.Fatalf("failed to initialize periph: %v", err***REMOVED***
***REMOVED***
***REMOVED***

func NewHWManager(***REMOVED*** *HWManager {
	return &HWManager{
		hw: make(map[string]interface{***REMOVED******REMOVED***,
***REMOVED***
***REMOVED***

func (m *HWManager***REMOVED*** Load(viper *viper.Viper***REMOVED*** error {
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
			if err := m.AddDevice(name, hardware***REMOVED***; err != nil {
				herr := err.(*HardwareError***REMOVED***
				switch herr.err {
				case ErrDevNotFound:
					unresolved = append(unresolved, name***REMOVED***
				default:
					return err
			***REMOVED***
		***REMOVED*** else {
				hasNewDevice = true
		***REMOVED***
	***REMOVED***
		names = unresolved
		if !(len(names***REMOVED*** > 0 && hasNewDevice***REMOVED*** {
			break
	***REMOVED***
***REMOVED***

	return nil
***REMOVED***

func (m *HWManager***REMOVED*** AddDevice(name string, viper *viper.Viper***REMOVED*** error {
	if _, ok := m.hw[name]; ok {
		return NewHardwareError(name, "name is used", ErrWrongConfig***REMOVED***
***REMOVED***
	if !viper.IsSet("type"***REMOVED*** {
		return NewHardwareError(name, "miss type field", ErrWrongConfig***REMOVED***
***REMOVED***

	t := viper.GetString("type"***REMOVED***
	switch t {
	case "spi":
		var spi spiwrap.SPIDevice
		err := viper.Unmarshal(&spi***REMOVED***
	***REMOVED***
			return NewHardwareError(name, err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
	***REMOVED***
		m.hw[name] = &spi
	case "uart":
		var uart UARTDevice
		if err := viper.Unmarshal(&uart***REMOVED***; err != nil {
			return NewHardwareError(name, err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
	***REMOVED***
		m.hw[name] = &uart
	case "tcp":
		var tcp TCP
		if err := viper.Unmarshal(&tcp***REMOVED***; err != nil {
			return NewHardwareError(name, err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
	***REMOVED***
		m.hw[name] = &tcp
	case "pwm":
		var pwm PWMDevice
		if err := viper.Unmarshal(&pwm***REMOVED***; err != nil {
			return NewHardwareError(name, err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
	***REMOVED***
		m.hw[name] = &pwm
	case "max31856":
		var conf max31856.Config

		if err := viper.Unmarshal(&conf***REMOVED***; err != nil {
			return NewHardwareError(name, err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
	***REMOVED***

		devName := viper.GetString("dev"***REMOVED***
		dev, ok := m.hw[devName]
		if !ok {
			msg := fmt.Sprintf("cannot find %s", devName***REMOVED***
			return NewHardwareError(name, msg, ErrDevNotFound***REMOVED***
	***REMOVED***
		spi := dev.(spiwrap.SPI***REMOVED***
		m.hw[name] = max31856.NewMAX31856(spi, conf***REMOVED***
	case "max31865":
		var conf max31865.Config

		if err := viper.Unmarshal(&conf***REMOVED***; err != nil {
			return NewHardwareError(name, err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
	***REMOVED***
		devName := viper.GetString("dev"***REMOVED***
		dev, ok := m.hw[devName]
		if !ok {
			msg := fmt.Sprintf("cannot find %s", devName***REMOVED***
			return NewHardwareError(name, msg, ErrDevNotFound***REMOVED***
	***REMOVED***
		spi := dev.(spiwrap.SPI***REMOVED***
		m.hw[name] = max31865.NewMAX31865(spi, conf***REMOVED***
	case "smoothie":
		var port SmoothiePort

		devName := viper.GetString("dev"***REMOVED***
		dev, ok := m.hw[devName]
		if !ok {
			msg := fmt.Sprintf("cannot find %s", devName***REMOVED***
			return NewHardwareError(name, msg, ErrDevNotFound***REMOVED***
	***REMOVED***
		port = dev.(SmoothiePort***REMOVED***
		m.hw[name] = NewSmoothie(port***REMOVED***
	case "extruder":
		var dev interface{***REMOVED***
		var port ExtruderPort
		var ok bool

		devName := viper.GetString("dev"***REMOVED***
		dev, ok = m.hw[devName]
		if !ok {
			msg := fmt.Sprintf("cannot find %s", devName***REMOVED***
			return NewHardwareError(name, msg, ErrDevNotFound***REMOVED***
	***REMOVED***
		port = dev.(ExtruderPort***REMOVED***
		m.hw[name] = NewExtruder(port***REMOVED***
	case "water_detector":
		var gpioWaterDetector GPIOWaterDetector
		if err := viper.Unmarshal(&gpioWaterDetector***REMOVED***; err != nil {
			return NewHardwareError(name, err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
	***REMOVED***
		m.hw[name] = &gpioWaterDetector
	default:
		return NewHardwareError(name, fmt.Sprintf("%s is not support yet", t***REMOVED***, ErrWrongConfig***REMOVED***
***REMOVED***

	return nil
***REMOVED***

func (m *HWManager***REMOVED*** GetDevice(name string***REMOVED*** (interface{***REMOVED***, error***REMOVED*** {
	if device, ok := m.hw[name]; ok {
		return device, nil
***REMOVED***
	return nil, fmt.Errorf("Cannot find %s device", name***REMOVED***
***REMOVED***
