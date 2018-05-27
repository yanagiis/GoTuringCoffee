package outtemp

import (
	"context"
	"math"
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
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

func (o *Service) Run(ctx context.Context, nc *nats.EncodedConn) (err error) {
	var reqSub *nats.Subscription
	var reqCh chan *nats.Msg

	reqCh = make(chan *nats.Msg)
	reqSub, err = nc.BindRecvChan("output.temperature", reqCh)
	if err != nil {
		return err
	}
	defer func() {
		err = reqSub.Unsubscribe()
		close(reqCh)
	}()

	var sensorErr error
	temperature := lib.TempRecord{
		Temp: math.NaN(),
		Time: time.Time{},
	}

	timer := time.NewTimer(o.ScanInterval)

	for {
		select {
		case msg := <-reqCh:
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
						Msg:  sensorErr.Error(),
					},
					Payload: temperature,
				}
			}
			nc.Publish(msg.Reply, resp)
		case <-timer.C:
			var temp float64
			timer = time.NewTimer(o.ScanInterval)
			if sensorErr = o.Sensor.Connect(); sensorErr != nil {
				log.Error().Msg(sensorErr.Error())
				continue
			}
			if temp, sensorErr = o.Sensor.GetTemperature(); sensorErr != nil {
				log.Error().Msg(sensorErr.Error())
				o.Sensor.Disconnect()
				continue
			}
			temperature.Temp = temp
			temperature.Time = time.Now()
		case <-ctx.Done():
			err = ctx.Err()
			return
		}
	}
}

func GetTemperature(ctx context.Context, nc *nats.EncodedConn) (resp lib.TempResponse, err error) {
	payload := lib.Request{
		Code: lib.CodeGet,
	}
	err = nc.RequestWithContext(ctx, "output.temperature", payload, &resp)
	return
}
