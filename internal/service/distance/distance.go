package distance

import (
	"context"
	"time"

	"GoTuringCoffee/internal/hardware"
	"GoTuringCoffee/internal/service/lib"

	nats "github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type Service struct {
	ScanInterval time.Duration
	Sensor       hardware.DistanceRangingSensor
}

func NewService(dev hardware.DistanceRangingSensor, scanInterval time.Duration) *Service {
	return &Service{
		ScanInterval: scanInterval,
		Sensor:       dev,
	}
}

func (o *Service) Run(ctx context.Context, nc *nats.EncodedConn, fin chan<- struct{}) (err error) {
	var sensorErr error
	distanceRecord := lib.DistanceRecord{
		Distance: -1,
		Time:     time.Time{},
	}

	nc.Subscribe("output.distance", func(subj, reply string, req lib.Request) {
		var resp lib.DistanceResponse
		if sensorErr != nil {
			resp = lib.DistanceResponse{
				Response: lib.Response{
					Code: lib.CodeFailure,
					Msg:  sensorErr.Error(),
				},
				Payload: lib.DistanceRecord{},
			}
		} else {
			resp = lib.DistanceResponse{
				Response: lib.Response{
					Code: lib.CodeSuccess,
				},
				Payload: distanceRecord,
			}
		}
		nc.Publish(reply, resp)
	})

	timer := time.NewTimer(o.ScanInterval)

	for {
		select {
		case <-timer.C:
			if sensorErr = o.Sensor.Open(); sensorErr != nil {
				log.Error().Msg(sensorErr.Error())
				timer = time.NewTimer(o.ScanInterval)
				continue
			}

			distanceRecord.Distance = int(o.Sensor.ReadRange())
			distanceRecord.Time = time.Now()
			timer = time.NewTimer(o.ScanInterval)
			log.Debug().Msgf("Distance: %d", distanceRecord.Distance)
		case <-ctx.Done():
			log.Info().Msg("stoping distance ranging service")
			err = ctx.Err()
			defer func() { fin <- struct{}{} }()
			log.Info().Msg("stop distance ranging service")
			return
		}
	}
}

func (o *Service) Stop() error {
	return nil
}

func GetDistance(ctx context.Context, nc *nats.EncodedConn) (resp lib.DistanceResponse, err error) {
	req := lib.Request{
		Code: lib.CodeGet,
	}
	err = nc.RequestWithContext(ctx, "output.distance", &req, &resp)
	if err != nil {
		log.Error().Err(err).Msg("Get distance ranging failed")
	}
	return
}
