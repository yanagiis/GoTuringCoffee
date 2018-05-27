package tankmeter

import (
	"context"
	"time"

	jsoniter "github.com/json-iterator/go"
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

func (t *Service) Run(ctx context.Context, nc *nats.EncodedConn) (err error) {
	var reqSub *nats.Subscription
	var reqCh chan *nats.Msg

	reqCh = make(chan *nats.Msg)
	reqSub, err = nc.BindRecvChan("tank.meter", reqCh)
	if err != nil {
		return err
	}
	defer func() {
		err = reqSub.Unsubscribe()
		close(reqCh)
	}()

	var sensorErr error
	fullRecord := lib.FullRecord{
		IsFull: false,
		Time:   time.Time{},
	}

	timer := time.NewTimer(t.ScanInterval)

	for {
		select {
		case msg := <-reqCh:
			var resp lib.FullResponse
			if sensorErr != nil {
				resp = lib.FullResponse{
					Response: lib.Response{
						Code: 1,
						Msg:  sensorErr.Error(),
					},
				}
			} else {
				resp = lib.FullResponse{
					Response: lib.Response{
						Code: 0,
					},
					Payload: fullRecord,
				}
			}
			nc.Publish(msg.Reply, resp)
		case <-timer.C:
			timer = time.NewTimer(t.ScanInterval)
			if sensorErr = t.Sensor.Connect(); sensorErr != nil {
				log.Error().Msg(sensorErr.Error())
				continue
			}
			fullRecord.IsFull, sensorErr = t.Sensor.IsWaterFull()
			fullRecord.Time = time.Now()
		case <-ctx.Done():
			err = ctx.Err()
			return
		}
	}
}

func GetMeterInfo(ctx context.Context, nc *nats.EncodedConn) (resp lib.FullResponse, err error) {
	payload, _ := jsoniter.Marshal(lib.Request{
		Code: lib.CodeGet,
	})
	err = nc.RequestWithContext(ctx, "tank.meter", payload, &resp)
	return
}
