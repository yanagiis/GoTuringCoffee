package hardware

// TemperatureSensor interface, such as MAX31856, MAX31865, etc...
type TemperatureSensor interface {
	Connect() error
	Disconnect() error
	GetTemperature() (float64, error)
}
