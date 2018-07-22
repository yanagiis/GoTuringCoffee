package hardware

import (
	rpio "github.com/stianeikeland/go-rpio"
)

type WaterDetector interface {
	Connect() error
	Disconnect() error
	IsWaterFull() (bool, error)
}

type GPIOWaterDetector struct {
	Pin  uint32
	gpio rpio.Pin
}

func (w *GPIOWaterDetector) Connect() error {
	w.gpio = rpio.Pin(w.Pin)
	w.gpio.Input()
	w.gpio.PullDown()
	return nil
}

func (w *GPIOWaterDetector) Disconnect() error {
	return nil
}

func (w *GPIOWaterDetector) IsWaterFull() (bool, error) {
	return w.gpio.Read() == rpio.High, nil
}
