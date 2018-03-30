package hardware

***REMOVED***
	"errors"
***REMOVED***
	"time"

	"github.com/yanagiis/periph/conn/gpio"
	"github.com/yanagiis/periph/conn/gpio/gpioreg"
***REMOVED***

type PWM interface {
	Connect(***REMOVED*** error
	Disconnect(***REMOVED*** error
	PWM(duty int64, period time.Duration***REMOVED*** error
***REMOVED***

type PWMDevice struct {
	Pin int32 `mapstructure:"pwm"`
	pwm gpio.PinPWM
***REMOVED***

type PWMConfig struct {
	Duty   float64       `mapstructure:"duty_cycle"`
	Period time.Duration `mapstructure:"period"`
***REMOVED***

func (p *PWMDevice***REMOVED*** Connect(***REMOVED*** error {
	if p.pwm != nil {
		return nil
***REMOVED***
	p.pwm = gpioreg.ByName(fmt.Sprint("PWM%d", p.Pin***REMOVED******REMOVED***.(gpio.PinPWM***REMOVED***
	return nil
***REMOVED***

func (p *PWMDevice***REMOVED*** Disconnect(***REMOVED*** error {
	p.pwm = nil
	return nil
***REMOVED***

func (p *PWMDevice***REMOVED*** PWM(duty int64, period time.Duration***REMOVED*** error {
	if p.pwm == nil {
		return errors.New("pwm not connected"***REMOVED***
***REMOVED***
	return p.pwm.PWM(gpio.Duty(duty***REMOVED***, period***REMOVED***
***REMOVED***
