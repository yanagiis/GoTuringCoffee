package outtemp

***REMOVED***
	"context"
	"math"
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
***REMOVED***

type Service struct {
	ScanInterval time.Duration
	Sensor       hardware.TemperatureSensor
***REMOVED***

func NewService(dev hardware.TemperatureSensor, scanInterval time.Duration***REMOVED*** *Service {
	return &Service{
		ScanInterval: scanInterval,
		Sensor:       dev,
***REMOVED***
***REMOVED***

func (o *Service***REMOVED*** Run(ctx context.Context, nc *nats.EncodedConn***REMOVED*** (err error***REMOVED*** {
	var reqSub *nats.Subscription
	var reqCh chan *nats.Msg

	reqCh = make(chan *nats.Msg***REMOVED***
	reqSub, err = nc.BindRecvChan("output.temperature", reqCh***REMOVED***
***REMOVED***
		return err
***REMOVED***
	defer func(***REMOVED*** {
		err = reqSub.Unsubscribe(***REMOVED***
		close(reqCh***REMOVED***
***REMOVED***(***REMOVED***

	var sensorErr error
	temperature := lib.TempRecord{
		Temp: math.NaN(***REMOVED***,
		Time: time.Time{***REMOVED***,
***REMOVED***

	timer := time.NewTimer(o.ScanInterval***REMOVED***

	for {
		select {
		case msg := <-reqCh:
			var resp lib.TempResponse
			if sensorErr != nil {
				resp = lib.TempResponse{
					Response: lib.Response{
						Code: lib.CodeFailure,
						Msg:  sensorErr.Error(***REMOVED***,
				***REMOVED***,
					Payload: lib.TempRecord{***REMOVED***,
			***REMOVED***
		***REMOVED*** else {
				resp = lib.TempResponse{
					Response: lib.Response{
						Code: lib.CodeSuccess,
						Msg:  sensorErr.Error(***REMOVED***,
				***REMOVED***,
					Payload: temperature,
			***REMOVED***
		***REMOVED***
			nc.Publish(msg.Reply, resp***REMOVED***
		case <-timer.C:
			var temp float64
			timer = time.NewTimer(o.ScanInterval***REMOVED***
			if sensorErr = o.Sensor.Connect(***REMOVED***; sensorErr != nil {
				log.Error(***REMOVED***.Msg(sensorErr.Error(***REMOVED******REMOVED***
				continue
		***REMOVED***
			if temp, sensorErr = o.Sensor.GetTemperature(***REMOVED***; sensorErr != nil {
				log.Error(***REMOVED***.Msg(sensorErr.Error(***REMOVED******REMOVED***
				o.Sensor.Disconnect(***REMOVED***
				continue
		***REMOVED***
			temperature.Temp = temp
			temperature.Time = time.Now(***REMOVED***
		case <-ctx.Done(***REMOVED***:
			err = ctx.Err(***REMOVED***
			return
	***REMOVED***
***REMOVED***
***REMOVED***

func GetTemperature(ctx context.Context, nc *nats.EncodedConn***REMOVED*** (resp lib.TempResponse, err error***REMOVED*** {
	payload := lib.Request{
		Code: lib.CodeGet,
***REMOVED***
	err = nc.RequestWithContext(ctx, "output.temperature", payload, &resp***REMOVED***
	return
***REMOVED***
