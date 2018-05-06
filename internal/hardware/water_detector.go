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
	IsWaterFull(***REMOVED*** (bool, error***REMOVED***
***REMOVED***

type GPIOWaterDetector struct {
	Pin  uint32
	gpio gpio.PinIO
***REMOVED***

func (w *GPIOWaterDetector***REMOVED*** Connect(***REMOVED*** error {
	gpioName := fmt.Sprintf("GPIO%d", w.Pin***REMOVED***
	w.gpio = gpioreg.ByName(gpioName***REMOVED***
	if w.gpio == nil {
		return fmt.Errorf("Cannot open %s", gpioName***REMOVED***
***REMOVED***
	if err := w.gpio.In(gpio.PullNoChange, gpio.NoEdge***REMOVED***; err != nil {
		return err
***REMOVED***
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
