package barista

import (
	"context"
	"runtime"
	"time"

	"GoTuringCoffee/internal/service/barista/middleware"
	"GoTuringCoffee/internal/service/lib"
	"GoTuringCoffee/internal/service/outtemp"
	"GoTuringCoffee/internal/service/replenisher"

	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
)

type Position struct {
	X float64 `mapstructure:"x"`
	Y float64 `mapstructure:"i"`
	Z float64 `mapstructure:"z"`
}

type BaristaConfig struct {
	PID                lib.NormalPID `mapstructure:"pid"`
	DrainPosition      Position      `mapstructure:"drain_position" validate:"nonzero"`
	DefaultMovingSpeed float64       `mapstructure:"default_moving_speed" validate:"nonzero"`
}

type Barista struct {
	conf       BaristaConfig
	middles    []middleware.Middleware
	controller Controller
	cooking    bool
}

func NewBarista(conf BaristaConfig, controller Controller) *Barista {
	return &Barista{
		conf:       conf,
		controller: controller,
		cooking:    false,
	}
}

func (b *Barista) Run(ctx context.Context, nc *nats.EncodedConn, fin chan<- struct{}) (err error) {
	var doneCh chan struct{}

	doneCh = make(chan struct{}, 0)

	nc.Subscribe("barista.brewing", func(subj, reply string, req lib.BaristaRequest) {
		if b.cooking {
			response(nc, reply, lib.CodeFailure, "Busy", nil)
			return
		}
		b.cooking = true
		response(nc, reply, lib.CodeSuccess, "OK", nil)
		go b.cook(ctx, nc, doneCh, req.Points)
	})

CONNECT:
	for {
		err := b.controller.Connect(ctx)
		if err == nil {
			defer b.controller.Disconnect()
			break CONNECT
		}
		select {
		case <-ctx.Done():
			break CONNECT
		case <-time.After(1000 * time.Millisecond):
			break
		}
	}

	timer := time.NewTimer(100 * time.Millisecond)

	for {
		select {
		case <-doneCh:
			doneCh = nil
		case <-ctx.Done():
			log.Info().Msg("stoping barista service")
			fin <- struct{}{}
			err = ctx.Err()
			log.Info().Msg("stop barista service")
			return
		case <-timer.C:
			timer = time.NewTimer(100 * time.Millisecond)
		}
	}
}

func (b *Barista) cook(ctx context.Context, nc *nats.EncodedConn, doneCh chan<- struct{}, points []lib.Point) {

	log.Debug().Msgf("Let's start cooking")
	b.middles = []middleware.Middleware{
		middleware.NewTempMiddleware(ctx, nc, &b.conf.PID, 20),
		middleware.NewTimeMiddleware(),
	}

	replenisher.StopReplenish(ctx, nc)
	runtime.LockOSThread()

	log.Debug().Msgf("Lock os thread")

	for i := range points {
		point := points[i]
		log.Debug().Msgf("%v", point)
		select {
		case <-ctx.Done():
			break
		default:
			if err := b.handlePoint(ctx, nc, &point); err != nil {
				log.Error().Err(err).Msgf("cook by point failed: point %v", point)
			}
		}
	}

	runtime.UnlockOSThread()
	replenisher.StartReplenish(ctx, nc)

	log.Debug().Msgf("Unlock os thread")

	for i := range b.middles {
		b.middles[i].Free()
	}
	b.middles = nil
	b.cooking = false

	doneCh <- struct{}{}
	log.Debug().Msgf("Cook finish")
}

func (b *Barista) handlePoint(ctx context.Context, nc *nats.EncodedConn, point *lib.Point) error {
	switch point.Type {
	case lib.WaitT:
		return b.handleWaitT(ctx, nc, point)
	case lib.MixT:
		return b.handleMixT(ctx, nc, point)
	default:
		return b.handlePointT(ctx, nc, point)
	}
}

func (b *Barista) handlePointT(ctx context.Context, nc *nats.EncodedConn, point *lib.Point) error {
	for _, middleware := range b.middles {
		middleware.Transform(point)
	}
	return b.controller.Do(point)
}

func (b *Barista) handleWaitT(ctx context.Context, nc *nats.EncodedConn, point *lib.Point) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Duration(*point.Time) * time.Second):
	}
	return nil
}

func (b *Barista) handleMixT(ctx context.Context, nc *nats.EncodedConn, point *lib.Point) error {
	b.moveToDrainPosition(nc, ctx)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(1 * time.Second):
	}

	e := float64(0.4)
	pointTime := float64(0.1)
	count := 0
	for j := 0; j < 50; j++ {
		for k := 0; k < 10; k++ {
			err := b.handlePointT(ctx, nc, &lib.Point{
				E:    &e,
				T:    point.T,
				Time: &pointTime,
			})
			if err != nil {
				return err
			}
		}
		r, err := outtemp.GetTemperature(ctx, nc)
		if err != nil {
			continue
		}
		if r.IsFailure() {
			continue
		}
		diff := r.Payload.Temp - *point.T
		if diff > 1 || diff < -1 {
			count = 0
			continue
		}
		if count < 3 {
			count++
			continue
		}
		break
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(1 * time.Second):
	}

	return nil
}

func (b *Barista) moveToDrainPosition(nc *nats.EncodedConn, ctx context.Context) {
	b.handlePointT(ctx, nc, &lib.Point{
		Z: &b.conf.DrainPosition.Z,
		F: &b.conf.DefaultMovingSpeed,
	})
	b.handlePointT(ctx, nc, &lib.Point{
		X: &b.conf.DrainPosition.X,
		Y: &b.conf.DrainPosition.Y,
		F: &b.conf.DefaultMovingSpeed,
	})
}

func (b *Barista) Stop() error {
	return nil
}

func response(nc *nats.EncodedConn, reply string, code uint8, msg string, payload interface{}) {
	resp := lib.Response{
		Code: code,
		Msg:  msg,
	}
	nc.Publish(reply, resp)
}

func Brew(ctx context.Context, nc *nats.EncodedConn, points []lib.Point) (resp lib.Response, err error) {
	req := lib.BaristaRequest{
		Request: lib.Request{
			Code: lib.CodePut,
		},
		Points: points,
	}
	err = nc.RequestWithContext(ctx, "barista.brewing", &req, &resp)
	return
}
