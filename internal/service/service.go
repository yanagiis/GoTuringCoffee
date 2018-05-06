package service

***REMOVED***
	"context"
	"errors"
***REMOVED***
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware/uartwrap"
	"github.com/yanagiis/GoTuringCoffee/internal/service/barista"
	"github.com/yanagiis/GoTuringCoffee/internal/service/heater"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
	"github.com/yanagiis/GoTuringCoffee/internal/service/mdns"
	"github.com/yanagiis/GoTuringCoffee/internal/service/outtemp"
	"github.com/yanagiis/GoTuringCoffee/internal/service/replenisher"
	"github.com/yanagiis/GoTuringCoffee/internal/service/tankmeter"
	"github.com/yanagiis/GoTuringCoffee/internal/service/tanktemp"
	"github.com/yanagiis/GoTuringCoffee/internal/service/uartserver"
	"github.com/yanagiis/GoTuringCoffee/internal/service/web"
***REMOVED***

var (
	ErrWrongConfig = errors.New("wrong configuration"***REMOVED***
***REMOVED***

type services interface {
	Run(ctx context.Context, nc *nats.EncodedConn***REMOVED*** error
***REMOVED***

type ServiceError struct {
	name string
	msg  string
	err  error
***REMOVED***

func NewServiceError(name string, msg string, err error***REMOVED*** *ServiceError {
	return &ServiceError{
		name: name,
		msg:  msg,
		err:  err,
***REMOVED***
***REMOVED***

func NewMissingFieldError(name string, field string***REMOVED*** *ServiceError {
	return NewServiceError(name, fmt.Sprintf("%q is required"***REMOVED***, ErrWrongConfig***REMOVED***
***REMOVED***

func (e *ServiceError***REMOVED*** Error(***REMOVED*** string {
	return fmt.Sprintf("%s: %s: %s", e.err.Error(***REMOVED***, e.name, e.msg***REMOVED***
***REMOVED***

type ServiceManager struct {
	services map[string]services
	cancels  map[string]context.CancelFunc
***REMOVED***

func NewServiceManager(***REMOVED*** *ServiceManager {
	return &ServiceManager{
		***REMOVED*** make(map[string]services***REMOVED***,
		cancels:  make(map[string]context.CancelFunc***REMOVED***,
***REMOVED***
***REMOVED***

func (s *ServiceManager***REMOVED*** Load(viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS***REMOVED*** error {
	if !viper.IsSet("services"***REMOVED*** {
		return NewServiceError("services", "no services field", ErrWrongConfig***REMOVED***
***REMOVED***

	services := viper.GetStringMapString("services"***REMOVED***
	for name := range services {
		services := viper.Sub(fmt.Sprintf("services.%s", name***REMOVED******REMOVED***
		if err := s.AddService(name, services, hwm, m***REMOVED***; err != nil {
			log.Fatal(***REMOVED***.Msgf("load service %q failed - %v", name, err***REMOVED***
			return err
	***REMOVED***
		log.Info(***REMOVED***.Msgf("load service %q", name***REMOVED***
***REMOVED***

	return nil
***REMOVED***

func (s *ServiceManager***REMOVED*** AddService(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS***REMOVED*** error {
	if _, ok := s.services[name]; ok {
		return NewServiceError(name, "name is used", ErrWrongConfig***REMOVED***
***REMOVED***
	if !viper.IsSet("enable"***REMOVED*** {
		return NewServiceError(name, "miss enable field", ErrWrongConfig***REMOVED***
***REMOVED***

	if !viper.GetBool("enable"***REMOVED*** {
		log.Debug(***REMOVED***.Msgf("%s is disable", name***REMOVED***
		return nil
***REMOVED***

	switch name {
	case "output_temp_service":
		if err := s.parseOutTempService(name, viper, hwm, m***REMOVED***; err != nil {
			return err
	***REMOVED***
	case "tank_temp_service":
		if err := s.parseTankTempService(name, viper, hwm, m***REMOVED***; err != nil {
			return err
	***REMOVED***
	case "tank_meter_service":
		if err := s.parseTankMeterService(name, viper, hwm, m***REMOVED***; err != nil {
			return err
	***REMOVED***
	case "replenisher_service":
		if err := s.parseReplenisher(name, viper, hwm, m***REMOVED***; err != nil {
			return err
	***REMOVED***
	case "heater":
		if err := s.parseHeater(name, viper, hwm, m***REMOVED***; err != nil {
			return err
	***REMOVED***
	case "barista":
		if err := s.parseBarista(name, viper, hwm, m***REMOVED***; err != nil {
			return err
	***REMOVED***
	case "uartserver":
		if err := s.parseUARTServer(name, viper, hwm, m***REMOVED***; err != nil {
			return err
	***REMOVED***
	case "web":
		if err := s.parseWeb(name, viper, hwm, m***REMOVED***; err != nil {
			return err
	***REMOVED***
***REMOVED***

	return nil
***REMOVED***

func (s *ServiceManager***REMOVED*** parseBarista(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS***REMOVED*** error {
	var conf barista.BaristaConfig
	var err error
	var tmp interface{***REMOVED***

	if err = checkFields(viper, []string{"smoothie", "extruder"***REMOVED******REMOVED***; err != nil {
		return err
***REMOVED***

	smoothieName := viper.GetString("smoothie"***REMOVED***
	if tmp, err = hwm.GetDevice(smoothieName***REMOVED***; err != nil {
		return err
***REMOVED***
	smoothie := tmp.(*hardware.Smoothie***REMOVED***

	extruderName := viper.GetString("extruder"***REMOVED***
	if tmp, err = hwm.GetDevice(extruderName***REMOVED***; err != nil {
		return err
***REMOVED***
	extruder := tmp.(*hardware.Extruder***REMOVED***

	if err := viper.Unmarshal(&conf***REMOVED***; err != nil {
		return NewServiceError(name, err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
***REMOVED***

	barista := barista.NewBarista(conf, &barista.SEController{
		Smoothie: smoothie,
		Extruder: extruder,
***REMOVED******REMOVED***
	s.services[name] = barista
	return nil
***REMOVED***

func (s *ServiceManager***REMOVED*** parseReplenisher(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS***REMOVED*** error {
	var pwmConfig hardware.PWMConfig

	devName := viper.GetString("dev"***REMOVED***
	dev, err := hwm.GetDevice(devName***REMOVED***
***REMOVED***
		return NewServiceError(name, err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
***REMOVED***
	pwm := dev.(hardware.PWM***REMOVED***

	scanInterval := time.Duration(viper.GetInt64("scan_interval_ms"***REMOVED******REMOVED*** * time.Millisecond
	if err = viper.Unmarshal(&pwmConfig***REMOVED***; err != nil {
		return NewServiceError(name, err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
***REMOVED***

	s.services[name] = replenisher.NewService(pwm, scanInterval, pwmConfig***REMOVED***
	return nil
***REMOVED***

func (s *ServiceManager***REMOVED*** parseOutTempService(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS***REMOVED*** error {
	devName := viper.GetString("dev"***REMOVED***
	dev, err := hwm.GetDevice(devName***REMOVED***
***REMOVED***
		return NewServiceError(name, err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
***REMOVED***
	sensor := dev.(hardware.TemperatureSensor***REMOVED***
	scanInterval := time.Duration(viper.GetInt64("scan_interval_ms"***REMOVED******REMOVED*** * time.Millisecond
	s.services[name] = outtemp.NewService(sensor, scanInterval***REMOVED***
	return nil
***REMOVED***

func (s *ServiceManager***REMOVED*** parseTankTempService(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS***REMOVED*** error {
	devName := viper.GetString("dev"***REMOVED***
	dev, err := hwm.GetDevice(devName***REMOVED***
***REMOVED***
		return NewServiceError(name, err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
***REMOVED***
	sensor := dev.(hardware.TemperatureSensor***REMOVED***
	scanInterval := time.Duration(viper.GetInt64("scan_interval_ms"***REMOVED******REMOVED*** * time.Millisecond
	s.services[name] = tanktemp.NewService(sensor, scanInterval***REMOVED***
	return nil
***REMOVED***

func (s *ServiceManager***REMOVED*** parseTankMeterService(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS***REMOVED*** error {
	devName := viper.GetString("dev"***REMOVED***
	dev, err := hwm.GetDevice(devName***REMOVED***
***REMOVED***
		return NewServiceError(name, err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
***REMOVED***
	sensor := dev.(hardware.WaterDetector***REMOVED***
	scanInterval := time.Duration(viper.GetInt64("scan_interval_ms"***REMOVED******REMOVED*** * time.Millisecond
	s.services[name] = tankmeter.NewService(sensor, scanInterval***REMOVED***
	return nil
***REMOVED***

func (s *ServiceManager***REMOVED*** parseHeater(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS***REMOVED*** error {
	var pid lib.NormalPID
	devName := viper.GetString("pwm_dev"***REMOVED***
	dev, err := hwm.GetDevice(devName***REMOVED***
***REMOVED***
		return NewServiceError(name, err.Error(***REMOVED***, ErrWrongConfig***REMOVED***
***REMOVED***
	pwm := dev.(hardware.PWM***REMOVED***
	scanInterval := time.Duration(viper.GetInt64("scan_interval_ms"***REMOVED******REMOVED*** * time.Millisecond
	viper.Sub("pid"***REMOVED***.Unmarshal(&pid***REMOVED***
	s.services[name] = heater.NewService(pwm, scanInterval, &pid***REMOVED***
	return nil
***REMOVED***

func (s *ServiceManager***REMOVED*** parseUARTServer(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS***REMOVED*** (err error***REMOVED*** {
	var device interface{***REMOVED***
	if err = checkFields(viper, []string{"mdns_service", "port", "uartdev"***REMOVED******REMOVED***; err != nil {
		return
***REMOVED***
	service := viper.GetString("mdns_service"***REMOVED***
	port := viper.GetInt("port"***REMOVED***
	if device, err = hwm.GetDevice(viper.GetString("uartdev"***REMOVED******REMOVED***; err != nil {
		return
***REMOVED***

	s.services[name] = uartserver.NewService(service, port, device.(uartwrap.UART***REMOVED***, m***REMOVED***
	return
***REMOVED***

func (s *ServiceManager***REMOVED*** parseWeb(name string, viper *viper.Viper, hwm *hardware.HWManager, m *mdns.MDNS***REMOVED*** (err error***REMOVED*** {
	if err = checkFields(viper, []string{"port", "static_files"***REMOVED******REMOVED***; err != nil {
		return
***REMOVED***
	port := viper.GetInt("port"***REMOVED***
	staticFiles := viper.GetString("static_files"***REMOVED***
	s.services[name] = web.NewService(port, staticFiles, m***REMOVED***
	return
***REMOVED***

func (s *ServiceManager***REMOVED*** RunServices(nc *nats.EncodedConn***REMOVED*** error {
	for name, service := range s.services {
		var ctx context.Context
		ctx, s.cancels[name] = context.WithCancel(context.Background(***REMOVED******REMOVED***
		go service.Run(ctx, nc***REMOVED***
***REMOVED***
	return nil
***REMOVED***

func (s *ServiceManager***REMOVED*** StopServices(***REMOVED*** error {
	for _, cancel := range s.cancels {
		cancel(***REMOVED***
***REMOVED***
	return nil
***REMOVED***

func checkFields(viper *viper.Viper, fields []string***REMOVED*** error {
	for _, field := range fields {
		if !viper.IsSet(field***REMOVED*** {
			return NewMissingFieldError("", field***REMOVED***
	***REMOVED***
***REMOVED***
	return nil
***REMOVED***
