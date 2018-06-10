package heater

import (
	"context"
	"errors"
	"time"

	jsoniter "github.com/json-iterator/go"
	nats "github.com/nats-io/go-nats"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
	"github.com/yanagiis/GoTuringCoffee/internal/service/tanktemp"
)

type Service struct {
	ScanInterval   time.Duration
	pwm            hardware.PWM
	pid            lib.PID
	targetTemp     float64
	sensorErr      error
	record         lib.HeaterRecord
	lastTempRecord lib.TempRecord
}

func NewService(dev hardware.PWM, scanInterval time.Duration, pid lib.PID) *Service {
	return &Service{
		ScanInterval: scanInterval,
		pwm:          dev,
		pid:          pid,
	}
}

func (h *Service) Run(ctx context.Context, nc *nats.EncodedConn) (err error) {
	var reqSub *nats.Subscription

	reqCh := make(chan *nats.Msg)
	reqSub, err = nc.BindRecvChan("tank.heater", reqCh)
	if err != nil {
		return err
	}
	defer func() {
		err = reqSub.Unsubscribe()
		close(reqCh)
	}()

	timer := time.NewTimer(h.ScanInterval)

	for {
		select {
		case msg := <-reqCh:
			var req lib.HeaterRequest
			if err := jsoniter.Unmarshal(msg.Data, &req); err != nil {
				nc.Publish(msg.Reply, lib.HeaterResponse{
					Response: lib.Response{
						Code: lib.CodeFailure,
						Msg:  err.Error(),
					},
				})
			}
			if req.IsGet() {
				resp := h.handleHeaterStatus()
				nc.Publish(msg.Reply, resp)
			}
			if req.IsPut() {
				resp := h.handleSetTemperature(req.Temp)
				nc.Publish(msg.Reply, resp)
			}
		case <-timer.C:
			h.adjustTemperature(ctx, nc)
			timer = time.NewTimer(h.ScanInterval)
		case <-ctx.Done():
			err = ctx.Err()
		}
	}
}

func (h *Service) handleHeaterStatus() lib.HeaterResponse {
	var resp lib.HeaterResponse
	if h.sensorErr != nil {
		resp = lib.HeaterResponse{
			Response: lib.Response{
				Code: lib.CodeFailure,
				Msg:  h.sensorErr.Error(),
			},
		}
	} else {
		resp = lib.HeaterResponse{
			Response: lib.Response{
				Code: lib.CodeSuccess,
			},
			Payload: h.record,
		}
	}
	return resp
}

func (h *Service) handleSetTemperature(temp float64) lib.HeaterResponse {
	h.targetTemp = temp
	return lib.HeaterResponse{
		Response: lib.Response{
			Code: lib.CodeSuccess,
		},
	}
}

func (h *Service) adjustTemperature(ctx context.Context, nc *nats.EncodedConn) error {
	resp, err := tanktemp.GetTemperature(ctx, nc)
	if err != nil {
		return err
	}
	if resp.IsFailure() {
		return errors.New("Cannot get tank temperature")
	}
	duty := h.pid.Compute(resp.Payload.Temp, resp.Payload.Time.Sub(h.lastTempRecord.Time))
	if err := h.pwm.PWM(duty, time.Second); err != nil {
		return err
	}
	h.record.Duty = duty
	h.record.Time = time.Now()
	h.lastTempRecord = resp.Payload
	return nil
}

func GetHeaterInfo(ctx context.Context, nc *nats.EncodedConn) (resp lib.HeaterResponse, err error) {
	req := lib.HeaterRequest{
		Request: lib.Request{
			Code: lib.CodePut,
		},
	}
	err = nc.RequestWithContext(ctx, "output.temperature", &req, &resp)
	return
}

func SetTemperature(ctx context.Context, nc *nats.EncodedConn, temp float64) (resp lib.TempResponse, err error) {
	req := lib.HeaterRequest{
		Request: lib.Request{
			Code: lib.CodePut,
		},
		Temp: temp,
	}
	err = nc.RequestWithContext(ctx, "output.temperature", &req, &resp)
	return
}
