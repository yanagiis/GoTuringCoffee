package hardware

***REMOVED***
	"errors"
***REMOVED***

	"github.com/yanagiis/periph/conn/gpio"
	"github.com/yanagiis/periph/conn/gpio/gpioreg"
***REMOVED***

type WaterDetector interface {
	Connect(***REMOVED*** error
	Disconnect(***REMOVED*** error
	IsWaterFull(***REMOVED*** bool
***REMOVED***

type GPIOWaterDetector struct {
	Pin  uint32 `mapstructure:"gpio"`
	gpio gpio.PinIO
***REMOVED***

func (w *GPIOWaterDetector***REMOVED*** Connect(***REMOVED*** error {
	w.gpio = gpioreg.ByName(fmt.Sprint("GPIO%lu", w.Pin***REMOVED******REMOVED***
	w.gpio.In(gpio.PullDown, gpio.NoEdge***REMOVED***
	return nil
***REMOVED***

func (w *GPIOWaterDetector***REMOVED*** Disconnect(***REMOVED*** error {
	w.gpio = nil
	return nil
***REMOVED***

func (w *GPIOWaterDetector***REMOVED*** IsWaterFull(***REMOVED*** (bool, error***REMOVED*** {
	if w.gpio == nil {
		return false, errors.New("gpio not connected"***REMOVED***
***REMOVED***
	return w.gpio.Read(***REMOVED*** == gpio.High, nil
***REMOVED***
