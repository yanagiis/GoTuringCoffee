package tanktemp

import (
	"context"
	"math"
	"time"

	"GoTuringCoffee/internal/hardware"
	"GoTuringCoffee/internal/service/lib"

	nats "github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type Service struct {
	ScanInterval time.Duration
	Sensor       hardware.TemperatureSensor
}

func NewService(dev hardware.TemperatureSensor, scanInterval time.Duration) *Service {
	return &Service{
		ScanInterval: scanInterval,
		Sensor:       dev,
	}
}

func (t *Service) Run(ctx context.Context, nc *nats.EncodedConn, fin chan<- struct{}) (err error) {
	var sensorErr error = nil
	temperature := lib.TempRecord{
		Temp: math.NaN(),
		Time: time.Time{},
	}

	nc.Subscribe("tank.temperature", func(subj, reply string, req lib.Request) {
		var resp lib.TempResponse
		if sensorErr != nil {
			resp = lib.TempResponse{
				Response: lib.Response{
					Code: lib.CodeFailure,
					Msg:  sensorErr.Error(),
				},
			}
		} else {
			resp = lib.TempResponse{
				Response: lib.Response{
					Code: lib.CodeSuccess,
				},
				Payload: temperature,
			}
		}
		nc.Publish(reply, resp)
	})

	timer := time.NewTimer(t.ScanInterval)

	for {
		select {
		case <-timer.C:
			if sensorErr = t.Sensor.Connect(); sensorErr != nil {
				log.Error().Msg(sensorErr.Error())
				timer = time.NewTimer(t.ScanInterval)
				continue
			}
			if temperature.Temp, sensorErr = t.Sensor.GetTemperature(); err != nil {
				t.Sensor.Disconnect()
				log.Error().Msg(sensorErr.Error())
				timer = time.NewTimer(t.ScanInterval)
				continue
			}
			temperature.Time = time.Now()
			timer = time.NewTimer(t.ScanInterval)
		case <-ctx.Done():
			log.Info().Msg("stoping tank temperature service")
			err = ctx.Err()
			defer func() { fin <- struct{}{} }()
			log.Info().Msg("stop tank temperature service")
			return
		}
	}
}

func (t *Service) Stop() error {
	return nil
}

func GetTemperature(ctx context.Context, nc *nats.EncodedConn) (resp lib.TempResponse, err error) {
	req := lib.Request{
		Code: lib.CodeGet,
	}
	err = nc.RequestWithContext(ctx, "tank.temperature", &req, &resp)
	if err != nil {
		log.Error().Err(err).Msg("Get tank temperature failed")
	}
	return
}
