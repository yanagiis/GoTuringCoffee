***REMOVED***

***REMOVED***
	"flag"
***REMOVED***

	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware/max31856"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware/spiwrap"
	"github.com/yanagiis/periph/conn/spi"
***REMOVED***

func init(***REMOVED*** {
	hardware.Init(***REMOVED***
***REMOVED***

func main(***REMOVED*** {
	var tc max31856.Type
	var temp float64
	var err error
	var spidev spiwrap.SPIDevice

	pathPtr := flag.String("path", "/dev/spidev0.0", "SPI device path"***REMOVED***
	speedPtr := flag.Int64("speed", 100000, "SPI speed"***REMOVED***
	modePtr := flag.Int64("mode", 1, "SPI mode"***REMOVED***
	bitsPtr := flag.Int("bits", 8, "SPI bits"***REMOVED***
	tcPtr := flag.String("tc", "T", "Sensor type"***REMOVED***
	flag.Parse(***REMOVED***

	tc, err = max31856.ParseType(*tcPtr***REMOVED***
***REMOVED***
		panic(err***REMOVED***
***REMOVED***

	spidev.Conf.Path = *pathPtr
	spidev.Conf.Speed = *speedPtr
	spidev.Conf.Mode = spi.Mode(*modePtr***REMOVED***
	spidev.Conf.Bits = *bitsPtr

	conf := max31856.Config{
		TC:   tc,
		Avg:  max31856.SampleAvg1,
		Mode: max31856.ModeAutomatic,
***REMOVED***

	sensor := max31856.New(&spidev, conf***REMOVED***
	if err := sensor.Connect(***REMOVED***; err != nil {
		panic(err***REMOVED***
***REMOVED***

	temp, err = sensor.GetTemperature(***REMOVED***
***REMOVED***
		panic(err***REMOVED***
***REMOVED***

	fmt.Printf("Temperature: %f\n", temp***REMOVED***

	sensor.Disconnect(***REMOVED***
***REMOVED***
