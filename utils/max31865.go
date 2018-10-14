package main

import (
	"flag"
	"fmt"

	"GoTuringCoffee/internal/hardware"
	"GoTuringCoffee/internal/hardware/max31865"
	"GoTuringCoffee/internal/hardware/spiwrap"

	"periph.io/x/periph/conn/spi"
)

func init() {
	hardware.Init()
}

func main() {
	var wire max31865.Wire
	var temp float64
	var err error
	var spidev spiwrap.SPIDevice

	pathPtr := flag.String("path", "/dev/spidev0.0", "SPI device path")
	speedPtr := flag.Int64("speed", 100000, "SPI speed")
	modePtr := flag.Int64("mode", 1, "SPI mode")
	bitsPtr := flag.Int("bits", 8, "SPI bits")
	wirePtr := flag.String("wire", "3", "PT100 wiring")
	flag.Parse()

	wire, err = max31865.ParseWire(*wirePtr)
	if err != nil {
		panic(err)
	}

	spidev.Conf.Path = *pathPtr
	spidev.Conf.Speed = *speedPtr
	spidev.Conf.Mode = spi.Mode(*modePtr)
	spidev.Conf.Bits = *bitsPtr

	conf := max31865.Config{
		Mode: max31865.ModeAutomatic,
		Wire: wire,
	}

	sensor := max31865.New(&spidev, conf)
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
