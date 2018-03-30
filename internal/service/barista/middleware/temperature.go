package middleware

***REMOVED***
	"math"
	"time"

	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
***REMOVED***

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
***REMOVED***

func NewTempMiddleware(pid lib.PID, maxAccWater float64***REMOVED*** *TempMiddleware {
	return &TempMiddleware{
		pid:            pid,
		lastMeasure:    time.Time{***REMOVED***,
		temp:           math.NaN(***REMOVED***,
		accWater:       0,
		maxAccWater:    maxAccWater,
		idealPercent:   math.NaN(***REMOVED***,
		currentPercent: math.NaN(***REMOVED***,
		tempChan:       nil,
***REMOVED***
***REMOVED***

func (m *TempMiddleware***REMOVED*** Transform(p *lib.Point***REMOVED*** {
	if p.T != nil && *p.T != m.temp {
		m.temp = *p.T
		m.accWater = 0
		m.idealPercent = (m.temp - m.lowTemp***REMOVED*** / (m.highTemp - m.lowTemp***REMOVED***
		m.pid.SetBound(-m.idealPercent, 1-m.idealPercent***REMOVED***
		m.pid.Reset(***REMOVED***
***REMOVED***

	if m.accWater > m.maxAccWater {
		if m.tempChan == nil {
			// m.tempChan = get_output_temperature(***REMOVED***
	***REMOVED***
		if record, ok := <-*m.tempChan; ok {
			duration := record.Time.Sub(m.lastMeasure***REMOVED***
			offset := m.pid.Compute(record.Temp, duration***REMOVED***
			m.currentPercent = m.idealPercent + offset
			m.lastMeasure = record.Time
			m.accWater = 0
	***REMOVED***
***REMOVED***

	if p.E != nil && *p.E != 0 {
		*p.E1 = *p.E * m.currentPercent
		*p.E2 = *p.E - *p.E1
		m.accWater += *p.E
***REMOVED***
***REMOVED***
