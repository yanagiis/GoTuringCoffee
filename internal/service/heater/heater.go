package heater

import (
	"context"
	"errors"
	"math"
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
	"github.com/yanagiis/GoTuringCoffee/internal/service/tanktemp"
)

type Service struct {
	ScanInterval   time.Duration
	pwm            hardware.PWM
	pid            lib.PID
	sensorErr      error
	record         lib.HeaterRecord
	lastTempRecord *lib.TempRecord
}

func NewService(dev hardware.PWM, scanInterval time.Duration, pid lib.PID) *Service {
	return &Service{
		ScanInterval: scanInterval,
		pwm:          dev,
		pid:          pid,
	}
}

func (h *Service) Run(ctx context.Context, nc *nats.EncodedConn, fin chan<- struct{}) (err error) {
	nc.Subscribe("tank.heater", func(subj, reply string, req lib.HeaterRequest) {
		if req.IsGet() {
			resp := h.handleHeaterStatus()
			nc.Publish(reply, resp)
		}
		if req.IsPut() {
			resp := h.handleSetTemperature(req.Temp)
			nc.Publish(reply, resp)
		}
	})

	if err = h.pwm.Connect(); err != nil {
		log.Info().Msg("Heater device connect failed")
		return
	}

	h.record.Target = 90
	h.pid.SetBound(0, 100)
	h.pid.SetPoint(h.record.Target)

	timer := time.NewTimer(h.ScanInterval)
	for {
		select {
		case <-timer.C:
			h.adjustTemperature(ctx, nc)
			timer = time.NewTimer(h.ScanInterval)
		case <-ctx.Done():
			log.Info().Msg("stoping heater service")
			h.pwm.Disconnect()
			err = ctx.Err()
			defer func() { fin <- struct{}{} }()
			log.Info().Msg("stop heater service")
			return
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
	log.Info().Msg("Set temperature")
	h.record.Target = temp
	h.pid.SetPoint(h.record.Target)
	return lib.HeaterResponse{
		Response: lib.Response{
			Code: lib.CodeSuccess,
		},
	}
}

func (h *Service) adjustTemperature(ctx context.Context, nc *nats.EncodedConn) error {
	resp, err := tanktemp.GetTemperature(ctx, nc)
	if err != nil {
		h.pwm.PWM(0, 0)
		return err
	}
	if resp.IsFailure() {
		h.pwm.PWM(0, 0)
		return errors.New("Cannot get tank temperature")
	}
	if math.IsNaN(resp.Payload.Temp) {
		h.pwm.PWM(0, 0)
		return errors.New("Cannot get NaN temperature")
	}
	if resp.Payload.Temp <= 0 {
		h.pwm.PWM(0, 0)
		return errors.New("Cannot get negative temperature")
	}
	log.Debug().Msgf("Tank: %f Time: %v", resp.Payload.Temp, resp.Payload.Time)
	difftime := time.Second * 0
	if h.lastTempRecord != nil {
		difftime = resp.Payload.Time.Sub(h.lastTempRecord.Time)
	}
	duty := h.pid.Compute(resp.Payload.Temp, difftime) / 100
	log.Debug().Msgf("Duty: %f", duty)
	if err := h.pwm.PWM(duty, 100); err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	h.record.Duty = duty
	h.record.Time = time.Now()
	h.lastTempRecord = &resp.Payload
	return nil
}

func GetHeaterInfo(ctx context.Context, nc *nats.EncodedConn) (resp lib.HeaterResponse, err error) {
	req := lib.HeaterRequest{
		Request: lib.Request{
			Code: lib.CodeGet,
		},
	}
	err = nc.RequestWithContext(ctx, "tank.heater", &req, &resp)
	return
}

func SetTemperature(ctx context.Context, nc *nats.EncodedConn, temp float64) (resp lib.HeaterResponse, err error) {
	req := lib.HeaterRequest{
		Request: lib.Request{
			Code: lib.CodePut,
		},
		Temp: temp,
	}
	err = nc.RequestWithContext(ctx, "tank.heater", &req, &resp)
	return
}
