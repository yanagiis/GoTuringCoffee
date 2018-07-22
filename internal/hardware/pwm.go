package hardware

import (
	"math"

	"github.com/rs/zerolog/log"
	rpio "github.com/stianeikeland/go-rpio"
)

type PWM interface {
	Connect() error
	Disconnect() error
	PWM(duty float64, period int64) error
}

type PWMConfig struct {
	Duty float64 `mapstructure:"duty_cycle"`
	Freq int64   `mapstructure:"frequency"`
}

type PWMDevice struct {
	Pin  int32 `mapstructure:"pwm"`
	pwm  rpio.Pin
	conf PWMConfig
}

func (p *PWMDevice) Connect() error {
	p.pwm = rpio.Pin(p.Pin)
	p.pwm.Pwm()
	return nil
}

func (p *PWMDevice) Disconnect() error {
	log.Info().Msg("Disconnecting PWM")
	rpio.StopPwm()
	p.pwm.Input()
	return nil
}

func (p *PWMDevice) PWM(duty float64, freq int64) error {
	if math.Abs(p.conf.Duty-duty) <= 1e-6 && p.conf.Freq == freq {
		return nil
	}
	p.conf.Duty = duty
	p.conf.Freq = freq
	rpio.StopPwm()
	dutyCount := uint32(duty * 100)
	cycleCount := uint32(100) - dutyCount
	log.Debug().Msgf("%v/%v", dutyCount, cycleCount)
	p.pwm.DutyCycle(dutyCount, cycleCount)
	p.pwm.Freq(int(freq))
	return nil
}
