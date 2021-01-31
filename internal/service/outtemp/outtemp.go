package outtemp

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

func (o *Service) Run(ctx context.Context, nc *nats.EncodedConn, fin chan<- struct{}) (err error) {
	var sensorErr error
	temperature := lib.TempRecord{
		Temp: math.NaN(),
		Time: time.Time{},
	}

	nc.Subscribe("output.temperature", func(subj, reply string, req lib.Request) {
		var resp lib.TempResponse
		if sensorErr != nil {
			resp = lib.TempResponse{
				Response: lib.Response{
					Code: lib.CodeFailure,
					Msg:  sensorErr.Error(),
				},
				Payload: lib.TempRecord{},
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

	timer := time.NewTimer(o.ScanInterval)

	for {
		select {
		case <-timer.C:
			var temp float64
			if sensorErr = o.Sensor.Connect(); sensorErr != nil {
				log.Error().Msg(sensorErr.Error())
				timer = time.NewTimer(o.ScanInterval)
				continue
			}
			if temp, sensorErr = o.Sensor.GetTemperature(); sensorErr != nil {
				log.Error().Msg(sensorErr.Error())
				o.Sensor.Disconnect()
				timer = time.NewTimer(o.ScanInterval)
				continue
			}
			temperature.Temp = temp
			temperature.Time = time.Now()
			timer = time.NewTimer(o.ScanInterval)
		case <-ctx.Done():
			log.Info().Msg("stoping output temperature service")
			err = ctx.Err()
			defer func() { fin <- struct{}{} }()
			log.Info().Msg("stop output temperature service")
			return
		}
	}
}

func (o *Service) Stop() error {
	return nil
}

func GetTemperature(ctx context.Context, nc *nats.EncodedConn) (resp lib.TempResponse, err error) {
	req := lib.Request{
		Code: lib.CodeGet,
	}
	err = nc.RequestWithContext(ctx, "output.temperature", &req, &resp)
	if err != nil {
		log.Error().Err(err).Msg("Get output temperature failed")
	}
	return
}
