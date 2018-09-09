package middleware

import (
	"context"
	"math"
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
	"github.com/yanagiis/GoTuringCoffee/internal/service/outtemp"
)

type TempMiddleware struct {
	pid             lib.PID
	lastMeasureTime time.Time
	temp            float64
	highTemp        float64
	lowTemp         float64
	accWater        float64
	maxAccWater     float64
	idealPercent    float64
	currentPercent  float64
	inChan          chan struct{}
	outChan         chan lib.TempRecord
	doneChan        chan struct{}
	cancel          context.CancelFunc
}

func NewTempMiddleware(ctx context.Context, nc *nats.EncodedConn, pid lib.PID, maxAccWater float64) *TempMiddleware {
	reqCtx, cancel := context.WithCancel(ctx)
	reqInCh := make(chan struct{})
	reqOutCh := make(chan lib.TempRecord)
	reqDoneCh := make(chan struct{})
	go requestTemp(reqCtx, nc, reqInCh, reqOutCh, reqDoneCh)
	return &TempMiddleware{
		pid:             pid,
		lastMeasureTime: time.Time{},
		temp:            math.NaN(),
		accWater:        0,
		maxAccWater:     maxAccWater,
		idealPercent:    math.NaN(),
		currentPercent:  0,
		inChan:          reqInCh,
		outChan:         reqOutCh,
		doneChan:        reqDoneCh,
		cancel:          cancel,
		highTemp:        90,
		lowTemp:         20,
	}
}

func requestTemp(ctx context.Context, nc *nats.EncodedConn, inCh <-chan struct{}, outCh chan<- lib.TempRecord, doneCh chan<- struct{}) {
	for {
		select {
		case <-inCh:
			for {
				r, err := outtemp.GetTemperature(ctx, nc)
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
			doneCh <- struct{}{}
			close(outCh)
			close(doneCh)
			return
		}
	}
}

func (m *TempMiddleware) Transform(p *lib.Point) {
	if p.T != nil && *p.T != m.temp {
		m.temp = *p.T
		m.accWater = 0
		m.idealPercent = (m.temp - m.lowTemp) / (m.highTemp - m.lowTemp)
		m.currentPercent = m.idealPercent
		m.lastMeasureTime = time.Time{}
		m.pid.SetPoint(m.temp)
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
			var duration time.Duration
			if m.lastMeasureTime.IsZero() {
				duration = time.Duration(0)
			} else {
				duration = record.Time.Sub(m.lastMeasureTime)
			}
			offset := m.pid.Compute(record.Temp, duration)
			m.currentPercent = m.idealPercent + offset
			m.lastMeasureTime = record.Time
			m.accWater = 0
		default:
		}
	}

	log.Debug().Msgf("percent %+v measure %+v ideal %+v", m.currentPercent, m.lastMeasureTime, m.idealPercent)

	if p.E != nil {
		p.E1 = new(float64)
		p.E2 = new(float64)
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
