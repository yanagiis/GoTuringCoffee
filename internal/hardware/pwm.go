package hardware

import (
	"context"

	"github.com/brian-armstrong/gpio"
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
	// if math.Abs(p.conf.Duty-duty) <= 1e-6 && p.conf.Freq == freq {
	// 	return nil
	// }
	p.conf.Duty = duty
	dutyCount := uint32(duty * 100)
	cycleCount := uint32(100)
	p.pwm.DutyCycle(dutyCount, cycleCount)
	if freq != 0 && p.conf.Freq != freq {
		p.conf.Freq = freq
		p.pwm.Freq(int(freq))
	}
	return nil
}

type PWMSoftware struct {
	PinNum uint `mapstructure:"pwm"`
	pin    gpio.Pin
	conf   PWMConfig
	cancel context.CancelFunc
	duty   float64
	freq   float64
}

func (p *PWMSoftware) Connect() error {
	var ctx context.Context
	p.pin = gpio.NewOutput(p.PinNum, false)
	ctx, p.cancel = context.WithCancel(context.Background())
	go func() {
	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			}
		}
		log.Debug().Msgf("Close PWM(software)")
		p.pin.Close()
	}()
	return nil
}

func (p *PWMSoftware) Disconnect() error {
	p.cancel()
	return nil
}

func (p *PWMSoftware) PWM(duty float64, freq int64) error {
	// if math.Abs(p.conf.Duty-duty) <= 1e-6 && p.conf.Freq == freq {
	// 	return nil
	// }

	log.Debug().Msgf("PWM(Soft) duty %f", duty)
	// FIXME: implement real pwm
	if duty > 0.05 {
		log.Debug().Msg("PWM(Soft) high")
		p.pin.High()
	} else {
		log.Debug().Msg("PWM(Soft) low")
		p.pin.Low()
	}
	return nil
}
