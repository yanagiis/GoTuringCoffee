package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"GoTuringCoffee/internal/hardware"
	"GoTuringCoffee/internal/hardware/i2cwrap"
	"GoTuringCoffee/internal/hardware/vl6180x"

	"github.com/rs/zerolog"
)

func init() {
	hardware.Init()
}

func main() {
	var err error
	var sensor *vl6180x.Vl6180x

	zerolog.SetGlobalLevel(zerolog.Disabled)

	pathPtr := flag.String("path", "/dev/i2c-1", "I2C device path")
	addressPtr := flag.Int("address", 0x29, "the i2c address of vl6180x")
	scalingPtr := flag.Int("scaling", 1, "Scaling 1 ~ 3")
	flag.Parse()

	i2cDevice := i2cwrap.NewI2C(*pathPtr)

	sensor, err = vl6180x.New(i2cDevice, *addressPtr, *scalingPtr)
	if err != nil {
		panic(err)
	}

	if err := sensor.Open(); err != nil {
		panic(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		sensor.Close()
		os.Exit(1)
	}()

	for {
		distance := sensor.ReadRange()
		fmt.Printf("Distance: %d mm\n", distance)
		time.Sleep(1 * time.Second) // or runtime.Gosched() or similar per @misterbee
	}
}
