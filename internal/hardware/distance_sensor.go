package hardware

// Distance interface, such as vl6180x, etc...
type DistanceRangingSensor interface {
	Open() error
	Close() error
	ReadRange() uint8
}
