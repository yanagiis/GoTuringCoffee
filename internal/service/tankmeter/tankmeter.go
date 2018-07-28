package tankmeter

import (
	"context"
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
)

type Service struct {
	ScanInterval time.Duration
	Sensor       hardware.WaterDetector
}

func NewService(dev hardware.WaterDetector, scanInterval time.Duration) *Service {
	return &Service{
		ScanInterval: scanInterval,
		Sensor:       dev,
	}
}

func (t *Service) Run(ctx context.Context, nc *nats.EncodedConn, fin chan<- struct{}) (err error) {
	var sensorErr error

	fullRecord := lib.FullRecord{
		IsFull: false,
		Time:   time.Time{},
	}

	nc.Subscribe("tank.meter", func(subj, reply string, req lib.Request) {
		var resp lib.FullResponse
		if sensorErr != nil {
			resp = lib.FullResponse{
				Response: lib.Response{
					Code: lib.CodeFailure,
					Msg:  sensorErr.Error(),
				},
			}
		} else {
			resp = lib.FullResponse{
				Response: lib.Response{
					Code: lib.CodeSuccess,
				},
				Payload: fullRecord,
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
			fullRecord.IsFull, sensorErr = t.Sensor.IsWaterFull()
			fullRecord.Time = time.Now()
			timer = time.NewTimer(t.ScanInterval)
		case <-ctx.Done():
			log.Info().Msg("stoping tank meter service")
			err = ctx.Err()
			defer func() { fin <- struct{}{} }()
			log.Info().Msg("stop tank meter service")
			return
		}
	}
}

func GetMeterInfo(ctx context.Context, nc *nats.EncodedConn) (resp lib.FullResponse, err error) {
	req := lib.Request{
		Code: lib.CodeGet,
	}
	err = nc.RequestWithContext(ctx, "tank.meter", &req, &resp)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return
}
