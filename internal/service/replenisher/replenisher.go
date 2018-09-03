package replenisher

import (
	"context"
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
	"github.com/yanagiis/GoTuringCoffee/internal/service/tankmeter"
)

type Service struct {
	ScanInterval time.Duration
	Dev          hardware.PWM
	PWMConf      hardware.PWMConfig
	devErr       error
	stop         bool
}

func NewService(dev hardware.PWM, scanInterval time.Duration, pwmConf hardware.PWMConfig) *Service {
	return &Service{
		ScanInterval: scanInterval,
		Dev:          dev,
		PWMConf:      pwmConf,
		stop:         false,
	}
}

func (r *Service) Run(ctx context.Context, nc *nats.EncodedConn, fin chan<- struct{}) (err error) {
	nc.Subscribe("tank.replenisher", func(subj, reply string, req lib.ReplenisherRequest) {
		if req.IsGet() {
			resp := r.handleReplenishStatus()
			nc.Publish(reply, resp)
		}
		if req.IsPut() {
			resp := r.handleControlReplenish(req.Stop)
			nc.Publish(reply, resp)
		}
	})

	if err = r.Dev.Connect(); err != nil {
		log.Info().Msg("Replenisher device connect failed")
		return
	}

	timer := time.NewTimer(r.ScanInterval)
	for {
		select {
		case <-timer.C:
			r.scan(ctx, nc)
			timer = time.NewTimer(r.ScanInterval)
		case <-ctx.Done():
			log.Info().Msg("stoping replenisher service")
			r.Dev.Disconnect()
			err = ctx.Err()
			defer func() { fin <- struct{}{} }()
			log.Info().Msg("stop replenisher service")
			return
		}
	}

}

func (r *Service) handleReplenishStatus() lib.ReplenisherResponse {
	if r.devErr != nil {
		return lib.ReplenisherResponse{
			Response: lib.Response{
				Code: lib.CodeFailure,
				Msg:  r.devErr.Error(),
			},
			Payload: lib.ReplenisherRecord{},
		}
	} else {
		return lib.ReplenisherResponse{
			Response: lib.Response{
				Code: lib.CodeSuccess,
			},
			Payload: lib.ReplenisherRecord{
				IsReplenishing: !r.stop,
				Time:           time.Now(),
			},
		}
	}
}

func (r *Service) handleControlReplenish(stop bool) lib.ReplenisherResponse {
	r.stop = stop
	return lib.ReplenisherResponse{
		Response: lib.Response{
			Code: lib.CodeSuccess,
		},
	}
}

func (r *Service) scan(ctx context.Context, nc *nats.EncodedConn) {
	var resp lib.FullResponse
	var err error

	duty := r.PWMConf.Duty
	if r.stop {
		duty = 0
	} else {
		timeCtx, _ := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
		if resp, err = tankmeter.GetMeterInfo(timeCtx, nc); err != nil {
			duty = 0
			log.Error().Msg(err.Error())
		}
		if resp.Payload.IsFull {
			duty = 0
		}
	}

	if err := r.Dev.PWM(duty, r.PWMConf.Freq); err != nil {
		log.Error().Msg(err.Error())
	}
}

func GetReplenishInfo(ctx context.Context, nc *nats.EncodedConn) (resp lib.ReplenisherResponse, err error) {
	req := lib.ReplenisherRequest{
		Request: lib.Request{
			Code: lib.CodeGet,
		},
	}
	err = nc.RequestWithContext(ctx, "tank.replenisher", &req, &resp)
	if err != nil {
		return
	}
	return
}

func StopReplenish(ctx context.Context, nc *nats.EncodedConn) (lib.ReplenisherResponse, error) {
	return toggleReplenish(ctx, nc, true)
}

func StartReplenish(ctx context.Context, nc *nats.EncodedConn) (lib.ReplenisherResponse, error) {
	return toggleReplenish(ctx, nc, false)
}

func toggleReplenish(ctx context.Context, nc *nats.EncodedConn, stop bool) (resp lib.ReplenisherResponse, err error) {
	req := lib.ReplenisherRequest{
		Request: lib.Request{
			Code: lib.CodeGet,
		},
		Stop: stop,
	}
	err = nc.RequestWithContext(ctx, "tank.replenisher", &req, &resp)
	return
}
