package hardware

import (
	"errors"
	"fmt"
	"time"

	"github.com/yanagiis/periph/conn/gpio"
	"github.com/yanagiis/periph/conn/gpio/gpioreg"
)

type PWM interface {
	Connect() error
	Disconnect() error
	PWM(duty float64, period time.Duration) error
}

type PWMConfig struct {
	Duty   float64       `mapstructure:"duty_cycle"`
	Period time.Duration `mapstructure:"period"`
}

type PWMDevice struct {
	Pin int32 `mapstructure:"pwm"`
	pwm gpio.PinPWM
}

func (p *PWMDevice) Connect() error {
	var ok bool
	if p.pwm != nil {
		return nil
	}
	p.pwm, ok = gpioreg.ByName(fmt.Sprintf("GPIO%d", p.Pin)).(gpio.PinPWM)
	if !ok {
		return fmt.Errorf("Cannot get 'GPIO%d' device", p.Pin)
	}
	return nil
}

func (p *PWMDevice) Disconnect() error {
	p.pwm = nil
	return nil
}

func (p *PWMDevice) PWM(duty float64, period time.Duration) error {
	if p.pwm == nil {
		return errors.New("pwm not connected")
	}
	dutyInt := gpio.Duty(duty * float64(gpio.DutyMax))
	return p.pwm.PWM(gpio.Duty(dutyInt), period)
}
