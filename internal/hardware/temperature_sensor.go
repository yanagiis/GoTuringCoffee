package hardware

// TemperatureSensor interface, such as MAX31856, MAX31865, etc...
type TemperatureSensor interface {
	Connect(***REMOVED*** error
	Disconnect(***REMOVED*** error
	GetTemperature(***REMOVED*** (float64, error***REMOVED***
***REMOVED***
