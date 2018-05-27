package tanktemp

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

func (t *Service) Run(ctx context.Context, nc *nats.EncodedConn) (err error) {
	var reqSub *nats.Subscription
	var reqCh chan *nats.Msg

	reqCh = make(chan *nats.Msg)
	reqSub, err = nc.BindRecvChan("tank.temperature", reqCh)
	if err != nil {
		return err
	}
	defer func() {
		err = reqSub.Unsubscribe()
		close(reqCh)
	}()

	var sensorErr error = nil
	temperature := lib.TempRecord{
		Temp: math.NaN(),
		Time: time.Time{},
	}

	timer := time.NewTimer(t.ScanInterval)

	for {
		select {
		case msg := <-reqCh:
			var resp lib.TempResponse
			if sensorErr != nil {
				resp = lib.TempResponse{
					Response: lib.Response{
						Code: 1,
						Msg:  sensorErr.Error(),
					},
				}
			} else {
				resp = lib.TempResponse{
					Response: lib.Response{
						Code: 0,
						Msg:  "",
					},
					Payload: temperature,
				}
			}
			nc.Publish(msg.Reply, resp)
		case <-timer.C:
			timer = time.NewTimer(t.ScanInterval)
			if sensorErr = t.Sensor.Connect(); sensorErr != nil {
				log.Error().Msg(sensorErr.Error())
				continue
			}
			if temperature.Temp, sensorErr = t.Sensor.GetTemperature(); err != nil {
				t.Sensor.Disconnect()
				continue
			}
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
	err = nc.RequestWithContext(ctx, "tank.temperature", payload, &resp)
	return
}
