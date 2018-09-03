package barista

import (
	"context"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
	"github.com/yanagiis/GoTuringCoffee/internal/service/barista/middleware"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
	"github.com/yanagiis/GoTuringCoffee/internal/service/tanktemp"
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
	var cookCtx context.Context
	var cookCancel context.CancelFunc
	var doneCh chan struct{}

	nc.Subscribe("barista.brewing", func(subj, reply string, req lib.BaristaRequest) {
		if b.cooking {
			response(nc, reply, lib.CodeFailure, "Budy", nil)
			return
		}
		b.cooking = true
		response(nc, reply, lib.CodeSuccess, "OK", nil)
		cookCtx, cookCancel = context.WithCancel(context.Background())
		go b.cook(cookCtx, nc, doneCh, req.Points)
	})

	if err := b.controller.Connect(); err != nil {
		return err
	}
	defer b.controller.Disconnect()

	timer := time.NewTimer(100 * time.Millisecond)

	for {
		select {
		case <-doneCh:
			cookCtx = nil
			cookCancel = nil
			doneCh = nil
		case <-ctx.Done():
			log.Info().Msg("stoping barista service")
			if cookCancel != nil {
				cookCancel()
				cookCancel = nil
			}
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

	for i := range points {
		point := points[i]

		select {
		case <-ctx.Done():
			break
		default:
			switch point.Type {
			case lib.WaitT:
				select {
				case <-ctx.Done():
				case <-time.After(time.Duration(*point.Time) * time.Second):
				}
			case lib.MixT:
				b.moveToDrainPosition(ctx)
				e := float64(0.2)
				time := float64(0.1)
				for j := 0; j < 100; j++ {
					for k := 0; k < 10; k++ {
						b.handlePoint(ctx, &lib.Point{
							E:    &e,
							T:    point.T,
							Time: &time,
						})
					}
					r, err := tanktemp.GetTemperature(ctx, nc)
					if err != nil {
						continue
					}
					if r.IsFailure() {
						continue
					}
					diff := r.Payload.Temp - *point.T
					if diff > 1 || diff < -1 {
						continue
					}
					break
				}
			case lib.PointT:
				b.handlePoint(ctx, &points[i])
			}
		}
	}
	for i := range b.middles {
		b.middles[i].Free()
	}
	b.middles = nil
	b.cooking = false
}

func (b *Barista) handlePoint(ctx context.Context, point *lib.Point) {
	for _, middleware := range b.middles {
		middleware.Transform(point)
	}
	b.controller.Do(point)
}

func (b *Barista) moveToDrainPosition(ctx context.Context) {
	b.handlePoint(ctx, &lib.Point{
		X: &b.conf.DrainPosition.X,
		Y: &b.conf.DrainPosition.Y,
		Z: &b.conf.DrainPosition.Z,
		F: &b.conf.DefaultMovingSpeed,
	})
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
