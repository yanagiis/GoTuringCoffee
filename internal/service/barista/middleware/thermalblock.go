package middleware

import (
	"context"

	"GoTuringCoffee/internal/service/lib"
	"GoTuringCoffee/internal/service/thermalblockheater"

	nats "github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type ThermalMiddleware struct {
	temp     float64
	inChan   chan float64
	doneChan chan struct{}
	cancel   context.CancelFunc
}

func NewThermalMiddleware(ctx context.Context, nc *nats.EncodedConn) *ThermalMiddleware {
	reqCtx, cancel := context.WithCancel(ctx)
	reqInCh := make(chan float64)
	reqDoneCh := make(chan struct{})

	go requestThermalHeater(reqCtx, nc, reqInCh, reqDoneCh)
	return &ThermalMiddleware{
		inChan:   reqInCh,
		doneChan: reqDoneCh,
		cancel:   cancel,
	}
}

func requestThermalHeater(ctx context.Context, nc *nats.EncodedConn, inCh <-chan float64, doneCh chan<- struct{}) {
	for {
	SELECT:
		select {
		case temperature := <-inCh:
			for {
				r, err := thermalblockheater.SetTemperature(ctx, nc, temperature)
				if err != nil {
					break SELECT
				}
				if r.IsFailure() {
					break SELECT
				}
				break
			}
		case <-ctx.Done():
			for {
				log.Debug().Msgf("set temperature to 0")
				newCtx := context.Background()
				r, err := thermalblockheater.SetTemperature(newCtx, nc, 0)
				if err != nil {
					log.Error().Err(err).Msgf("")
					continue
				}
				if r.IsFailure() {
					log.Error().Msgf("Set temperature to 0 failed")
					continue
				}
				break
			}
			doneCh <- struct{}{}
			close(doneCh)
			return
		}
	}
}

func (m *ThermalMiddleware) Transform(p *lib.Point) {
	if p.T != nil && *p.T != m.temp {
		m.temp = *p.T
		m.inChan <- m.temp
	}
}

func (m *ThermalMiddleware) Free() {
	m.cancel()
	close(m.inChan)
	<-m.doneChan
}
