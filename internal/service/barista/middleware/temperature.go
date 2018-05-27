package middleware

import (
	"math"
	"time"

	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
)

type TempMiddleware struct {
	pid            lib.PID
	lastMeasure    time.Time
	temp           float64
	highTemp       float64
	lowTemp        float64
	accWater       float64
	maxAccWater    float64
	idealPercent   float64
	currentPercent float64
	tempChan       *chan lib.TempRecord
}

func NewTempMiddleware(pid lib.PID, maxAccWater float64) *TempMiddleware {
	return &TempMiddleware{
		pid:            pid,
		lastMeasure:    time.Time{},
		temp:           math.NaN(),
		accWater:       0,
		maxAccWater:    maxAccWater,
		idealPercent:   math.NaN(),
		currentPercent: math.NaN(),
		tempChan:       nil,
	}
}

func (m *TempMiddleware) Transform(p *lib.Point) {
	if p.T != nil && *p.T != m.temp {
		m.temp = *p.T
		m.accWater = 0
		m.idealPercent = (m.temp - m.lowTemp) / (m.highTemp - m.lowTemp)
		m.pid.SetBound(-m.idealPercent, 1-m.idealPercent)
		m.pid.Reset()
	}

	if m.accWater > m.maxAccWater {
		if m.tempChan == nil {
			// m.tempChan = get_output_temperature()
		}
		if record, ok := <-*m.tempChan; ok {
			duration := record.Time.Sub(m.lastMeasure)
			offset := m.pid.Compute(record.Temp, duration)
			m.currentPercent = m.idealPercent + offset
			m.lastMeasure = record.Time
			m.accWater = 0
		}
	}

	if p.E != nil && *p.E != 0 {
		*p.E1 = *p.E * m.currentPercent
		*p.E2 = *p.E - *p.E1
		m.accWater += *p.E
	}
}
