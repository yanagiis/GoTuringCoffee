package barista

import (
	"context"
	"time"

	"github.com/nats-io/go-nats"
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

func (b *Barista) Run(ctx context.Context, nc *nats.EncodedConn) (err error) {
	var cookSub, querySub *nats.Subscription
	var cookCh, queryCh chan *nats.Msg

	cookCh = make(chan *nats.Msg)
	cookSub, err = nc.BindRecvChan("barista.cooking", cookCh)
	if err != nil {
		return
	}
	defer func() {
		err = cookSub.Unsubscribe()
		close(cookCh)
	}()

	queryCh = make(chan *nats.Msg, 16)
	querySub, err = nc.BindRecvChan("barista.query", queryCh)
	if err != nil {
		return
	}
	defer func() {
		err = querySub.Unsubscribe()
		close(queryCh)
	}()

	var cookCtx context.Context
	var cookCancel context.CancelFunc
	var doneCh chan struct{}

	timer := time.NewTimer(100 * time.Millisecond)

	for {
		select {
		case req := <-cookCh:
			var points []lib.Point
			response(nc, req, lib.CodeSuccess, "OK", nil)
			cookCtx, cookCancel = context.WithCancel(context.Background())
			go b.cook(cookCtx, doneCh, points)
		case <-queryCh:
			// b.query(ctx, query)
		case <-doneCh:
			cookCtx = nil
			cookCancel = nil
			doneCh = nil
		case <-ctx.Done():
			if cookCancel != nil {
				cookCancel()
				cookCancel = nil
			}
			if doneCh != nil {
				<-doneCh
				doneCh = nil
				cookCtx = nil
			}
			err = ctx.Err()
			break
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

func response(nc *nats.EncodedConn, req *nats.Msg, code uint8, msg string, payload interface{}) {
	resp := lib.Response{
		Code: code,
		Msg:  msg,
	}
	nc.Publish(req.Reply, resp)
}
