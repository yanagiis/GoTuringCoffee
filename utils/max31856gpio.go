package main

import (
	"flag"
	"fmt"

	"GoTuringCoffee/internal/hardware"
	"GoTuringCoffee/internal/hardware/max31856"
	"GoTuringCoffee/internal/hardware/spiwrap"

	"periph.io/x/periph/conn/spi"
)

func init() {
	hardware.Init()
}

func main() {
	var temp float64
	var err error
	var spidev spiwrap.SPIGPIO

	pathPtr := flag.String("path", "/dev/spidev0.0", "SPI device path")
	speedPtr := flag.Int64("speed", 100000, "SPI speed")
	modePtr := flag.Int64("mode", 1, "SPI mode")
	bitsPtr := flag.Int("bits", 8, "SPI bits")
	flag.Parse()

	spidev.Conf.Path = *pathPtr
	spidev.Conf.Speed = *speedPtr
	spidev.Conf.Mode = spi.Mode(*modePtr)
	spidev.Conf.Bits = *bitsPtr
	spidev.Pins.CLK = 0
	spidev.Pins.MOSI = 1
	spidev.Pins.MISO = 2
	spidev.Pins.CS = 3

	conf := max31856.Config{
		TC:   max31856.TypeT,
		Mode: max31856.ModeAutomatic,
		Avg:  max31856.SampleAvg1,
	}

	sensor := max31856.New(&spidev, conf)
	if err := sensor.Connect(); err != nil {
		panic(err)
	}

	temp, err = sensor.GetTemperature()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Temperature: %f\n", temp)

	sensor.Disconnect()
}
