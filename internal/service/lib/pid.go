package lib

***REMOVED***
	"time"
***REMOVED***

type PID interface {
	GetParam(***REMOVED*** (float64, float64, float64***REMOVED***
	SetParam(p float64, i float64, d float64***REMOVED***
	SetPoint(point float64***REMOVED***
	SetBound(lower float64, upper float64***REMOVED***
	Compute(measure float64, duration time.Duration***REMOVED*** float64
	Reset(***REMOVED***
***REMOVED***

type NormalPID struct {
	P           float64 `mapstructure:"p"`
	I           float64 `mapstructure:"i"`
	D           float64 `mapstructure:"d"`
	iterm       float64
	setpoint    float64
	lastMeasure float64
	lower       float64
	upper       float64
	reset       bool
***REMOVED***

func NewNormalPID(p float64, i float64, d float64***REMOVED*** *NormalPID {
	pid := &NormalPID{
		P: p, I: i, D: d,
***REMOVED***
	pid.Reset(***REMOVED***
	return pid
***REMOVED***

func (pid *NormalPID***REMOVED*** GetParam(***REMOVED*** (float64, float64, float64***REMOVED*** {
	return pid.P, pid.I, pid.D
***REMOVED***

func (pid *NormalPID***REMOVED*** SetParam(p float64, i float64, d float64***REMOVED*** {
	pid.P = p
	pid.I = i
	pid.D = d
***REMOVED***

func (pid *NormalPID***REMOVED*** SetPoint(point float64***REMOVED*** {
	pid.setpoint = point
***REMOVED***

func (pid *NormalPID***REMOVED*** SetBound(lower float64, upper float64***REMOVED*** {
	pid.lower = lower
	pid.upper = upper
***REMOVED***

func (pid *NormalPID***REMOVED*** Compute(measure float64, duration time.Duration***REMOVED*** float64 {
	if pid.reset {
		pid.lastMeasure = measure
		pid.reset = false
***REMOVED***

	time := float64(duration.Nanoseconds(***REMOVED******REMOVED*** / 1000000000
	err := pid.setpoint - measure
	p := pid.P * err

	if time == 0 {
		return pid.limitValue(p***REMOVED***
***REMOVED***

	dMeasure := measure - pid.lastMeasure
	pid.lastMeasure = measure
	i := pid.I*err*time + pid.iterm
	d := pid.D * (dMeasure / time***REMOVED***
	pid.iterm = pid.limitValue(i***REMOVED***

	return pid.limitValue(p + pid.iterm + d***REMOVED***
***REMOVED***

func (pid *NormalPID***REMOVED*** Reset(***REMOVED*** {
	pid.reset = true
	pid.iterm = 0
***REMOVED***

func (pid *NormalPID***REMOVED*** limitValue(value float64***REMOVED*** float64 {
	switch {
	case value > pid.upper:
		return pid.upper
	case value < pid.lower:
		return pid.lower
	default:
		return value
***REMOVED***
***REMOVED***
