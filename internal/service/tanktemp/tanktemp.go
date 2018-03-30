package tanktemp

***REMOVED***
	"context"
	"math"
	"time"

	nats "github.com/nats-io/go-nats"
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

func (t *Service***REMOVED*** Run(ctx context.Context, nc *nats.EncodedConn***REMOVED*** (err error***REMOVED*** {
	var reqSub *nats.Subscription
	var reqCh chan *nats.Msg

	reqCh = make(chan *nats.Msg***REMOVED***
	reqSub, err = nc.BindRecvChan("tank.temperature", reqCh***REMOVED***
***REMOVED***
		return err
***REMOVED***
	defer func(***REMOVED*** {
		err = reqSub.Unsubscribe(***REMOVED***
		close(reqCh***REMOVED***
***REMOVED***(***REMOVED***

	var sensorErr error = nil
	temperature := lib.TempRecord{
		Temp: math.NaN(***REMOVED***,
		Time: time.Time{***REMOVED***,
***REMOVED***

	timer := time.NewTimer(t.ScanInterval***REMOVED***

	for {
		select {
		case msg := <-reqCh:
			var resp lib.TempResponse
			if sensorErr != nil {
				resp = lib.TempResponse{
					Response: lib.Response{
						Code: 1,
						Msg:  sensorErr.Error(***REMOVED***,
				***REMOVED***,
			***REMOVED***
		***REMOVED*** else {
				resp = lib.TempResponse{
					Response: lib.Response{
						Code: 0,
						Msg:  "",
				***REMOVED***,
					Payload: temperature,
			***REMOVED***
		***REMOVED***
			nc.Publish(msg.Reply, resp***REMOVED***
		case <-timer.C:
			if sensorErr = t.Sensor.Connect(***REMOVED***; sensorErr != nil {
				continue
		***REMOVED***
			if temperature.Temp, sensorErr = t.Sensor.GetTemperature(***REMOVED***; err != nil {
				t.Sensor.Disconnect(***REMOVED***
				continue
		***REMOVED***
			temperature.Time = time.Now(***REMOVED***
			timer = time.NewTimer(t.ScanInterval***REMOVED***
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
	err = nc.RequestWithContext(ctx, "tank.temperature", payload, &resp***REMOVED***
	return
***REMOVED***
