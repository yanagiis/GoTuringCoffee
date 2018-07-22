package barista

import (
	"context"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
	"github.com/yanagiis/GoTuringCoffee/internal/service/barista/middleware"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
)

type Position struct {
	x float64 `mapstructure:"x"`
	y float64 `mapstructure:"y"`
	z float64 `mapstructure:"z"`
}

type BaristaConfig struct {
	PID                lib.NormalPID `mapstructure:"pid"`
	WasteWaterPosition Position      `mapstructure:"waste_water_position"`
	DefaultMovingSpeed float64       `mapstructure:"default_moving_speed"`
}

type Barista struct {
	conf       BaristaConfig
	middles    []middleware.Middleware
	controller Controller
}

func NewBarista(conf BaristaConfig, controller Controller) *Barista {
	middles := []middleware.Middleware{
		middleware.NewTempMiddleware(&conf.PID, 20),
		middleware.NewTimeMiddleware(),
	}
	return &Barista{
		conf:       conf,
		middles:    middles,
		controller: controller,
	}
}

func (b *Barista) Run(ctx context.Context, nc *nats.EncodedConn, fin chan<- struct{}) (err error) {
	var cookCtx context.Context
	var cookCancel context.CancelFunc
	var doneCh chan struct{}

	nc.Subscribe("barista.brewing", func(subj, reply string, points []lib.Point) {
		response(nc, reply, lib.CodeSuccess, "OK", nil)
		cookCtx, cookCancel = context.WithCancel(context.Background())
		go b.cook(cookCtx, doneCh, points)
	})

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

func (b *Barista) cook(ctx context.Context, doneCh chan<- struct{}, points []lib.Point) {
	for i := range points {
		if _, ok := <-ctx.Done(); ok {
			break
		}
		for _, middleware := range b.middles {
			middleware.Transform(&points[i])
		}
		b.controller.Do(&points[i])
	}
	doneCh <- struct{}{}
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
