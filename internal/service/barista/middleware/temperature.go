package middleware

import (
	"context"
	"math"
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
	"github.com/yanagiis/GoTuringCoffee/internal/service/tanktemp"
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
	inChan         chan struct{}
	outChan        chan lib.TempRecord
	doneChan       chan struct{}
	cancel         context.CancelFunc
}

func NewTempMiddleware(ctx context.Context, nc *nats.EncodedConn, pid lib.PID, maxAccWater float64) *TempMiddleware {
	reqCtx, cancel := context.WithCancel(ctx)
	reqInCh := make(chan struct{})
	reqOutCh := make(chan lib.TempRecord)
	reqDoneCh := make(chan struct{})
	go requestTemp(reqCtx, nc, reqInCh, reqOutCh, reqDoneCh)
	return &TempMiddleware{
		pid:            pid,
		lastMeasure:    time.Time{},
		temp:           math.NaN(),
		accWater:       0,
		maxAccWater:    maxAccWater,
		idealPercent:   math.NaN(),
		currentPercent: math.NaN(),
		inChan:         reqInCh,
		outChan:        reqOutCh,
		doneChan:       reqDoneCh,
		cancel:         cancel,
	}
}

func requestTemp(ctx context.Context, nc *nats.EncodedConn, inCh <-chan struct{}, outCh chan<- lib.TempRecord, doneCh chan<- struct{}) {
	select {
	case <-inCh:
		for {
			r, err := tanktemp.GetTemperature(ctx, nc)
			if err != nil {
				continue
			}
			if r.IsFailure() {
				continue
			}
			outCh <- r.Payload
			break
		}
	case <-ctx.Done():
		break
	}
	doneCh <- struct{}{}
	close(outCh)
	close(doneCh)
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
		select {
		case m.inChan <- struct{}{}:
		default:
		}
		select {
		case record := <-m.outChan:
			duration := record.Time.Sub(m.lastMeasure)
			offset := m.pid.Compute(record.Temp, duration)
			m.currentPercent = m.idealPercent + offset
			m.lastMeasure = record.Time
			m.accWater = 0
		default:
		}
	}

	if p.E != nil && *p.E != 0 {
		*p.E1 = *p.E * m.currentPercent
		*p.E2 = *p.E - *p.E1
		m.accWater += *p.E
	}
}

func (m *TempMiddleware) Free() {
	m.cancel()
	close(m.inChan)
	<-m.doneChan
}
