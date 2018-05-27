package replenisher

import (
	"context"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
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
	}
}

func (r *Service) Run(ctx context.Context, nc *nats.EncodedConn) (err error) {
	var reqSub *nats.Subscription
	var reqCh chan *nats.Msg

	reqCh = make(chan *nats.Msg)
	reqSub, err = nc.BindRecvChan("tank.replenisher", reqCh)
	if err != nil {
		return err
	}
	defer func() {
		err = reqSub.Unsubscribe()
		close(reqCh)
	}()

	r.Dev.Connect()
	defer r.Dev.Disconnect()
	timer := time.NewTimer(r.ScanInterval)

	for {
		select {
		case msg := <-reqCh:
			var req lib.ReplenisherRequest
			if decodeErr := jsoniter.Unmarshal(msg.Data, &req); decodeErr != nil {
				nc.Publish(msg.Reply, lib.HeaterResponse{
					Response: lib.Response{
						Code: lib.CodeFailure,
						Msg:  decodeErr.Error(),
					},
				})
			}
			if req.IsGet() {
				resp := r.handleReplenishStatus()
				nc.Publish(msg.Reply, resp)
			}
			if req.IsPut() {
				resp := r.handleControlReplenish(req.Stop)
				nc.Publish(msg.Reply, resp)
			}
		case <-timer.C:
			r.scan(ctx, nc)
			timer = time.NewTimer(r.ScanInterval)
		case <-ctx.Done():
			err = ctx.Err()
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

	if err = r.Dev.Connect(); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	duty := r.PWMConf.Duty
	if r.stop {
		duty = 0
	}
	timeCtx, _ := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	if resp, err = tankmeter.GetMeterInfo(timeCtx, nc); err != nil {
		duty = 0
		log.Error().Msg(err.Error())
	}
	if resp.Payload.IsFull {
		duty = 0
	}

	fmt.Printf("%v\n", duty)
	if err := r.Dev.PWM(duty, r.PWMConf.Period); err != nil {
		fmt.Printf("%v\n", err)
	}
}

func GetReplenishInfo(ctx context.Context, nc *nats.EncodedConn) (resp lib.ReplenisherResponse, err error) {
	payload := lib.ReplenisherRequest{
		Request: lib.Request{
			Code: lib.CodeGet,
		},
	}
	err = nc.RequestWithContext(ctx, "tank.meter", payload, &resp)
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
	payload := lib.ReplenisherRequest{
		Request: lib.Request{
			Code: lib.CodeGet,
		},
		Stop: stop,
	}
	err = nc.RequestWithContext(ctx, "tank.replenisher", payload, &resp)
	return
}
