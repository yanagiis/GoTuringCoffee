package main

import (
	"flag"
	"fmt"

	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware/max31856"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware/spiwrap"
	"periph.io/x/periph/conn/spi"
)

func init() {
	hardware.Init()
}

func main() {
	var tc max31856.Type
	var temp float64
	var err error
	var spidev spiwrap.SPIDevice

	pathPtr := flag.String("path", "/dev/spidev0.0", "SPI device path")
	speedPtr := flag.Int64("speed", 100000, "SPI speed")
	modePtr := flag.Int64("mode", 1, "SPI mode")
	bitsPtr := flag.Int("bits", 8, "SPI bits")
	tcPtr := flag.String("tc", "T", "Sensor type")
	flag.Parse()

	tc, err = max31856.ParseType(*tcPtr)
	if err != nil {
		panic(err)
	}

	spidev.Conf.Path = *pathPtr
	spidev.Conf.Speed = *speedPtr
	spidev.Conf.Mode = spi.Mode(*modePtr)
	spidev.Conf.Bits = *bitsPtr

	conf := max31856.Config{
		TC:   tc,
		Avg:  max31856.SampleAvg1,
		Mode: max31856.ModeAutomatic,
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
