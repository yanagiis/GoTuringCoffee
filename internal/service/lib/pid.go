package lib

import (
	"time"
)

type PID interface {
	GetParam() (float64, float64, float64)
	SetParam(p float64, i float64, d float64)
	SetPoint(point float64)
	SetBound(lower float64, upper float64)
	Compute(measure float64, duration time.Duration) float64
	Reset()
}

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
}

func NewNormalPID(p float64, i float64, d float64) *NormalPID {
	pid := &NormalPID{
		P: p, I: i, D: d,
	}
	pid.Reset()
	return pid
}

func (pid *NormalPID) GetParam() (float64, float64, float64) {
	return pid.P, pid.I, pid.D
}

func (pid *NormalPID) SetParam(p float64, i float64, d float64) {
	pid.P = p
	pid.I = i
	pid.D = d
}

func (pid *NormalPID) SetPoint(point float64) {
	pid.setpoint = point
}

func (pid *NormalPID) SetBound(lower float64, upper float64) {
	pid.lower = lower
	pid.upper = upper
}

func (pid *NormalPID) Compute(measure float64, duration time.Duration) float64 {
	if pid.reset {
		pid.lastMeasure = measure
		pid.reset = false
	}

	time := float64(duration.Nanoseconds()) / 1000000000
	err := pid.setpoint - measure
	p := pid.P * err

	if time == 0 {
		return pid.limitValue(p)
	}

	dMeasure := measure - pid.lastMeasure
	pid.lastMeasure = measure
	i := pid.I*err*time + pid.iterm
	d := pid.D * (dMeasure / time)
	pid.iterm = pid.limitValue(i)

	return pid.limitValue(p + pid.iterm + d)
}

func (pid *NormalPID) Reset() {
	pid.reset = true
	pid.iterm = 0
}

func (pid *NormalPID) limitValue(value float64) float64 {
	switch {
	case value > pid.upper:
		return pid.upper
	case value < pid.lower:
		return pid.lower
	default:
		return value
	}
}
