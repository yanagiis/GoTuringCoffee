package tankmeter

***REMOVED***
	"context"
	"time"

	jsoniter "github.com/json-iterator/go"
	nats "github.com/nats-io/go-nats"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
***REMOVED***

type Service struct {
	ScanInterval time.Duration
	Sensor       hardware.WaterDetector
***REMOVED***

func NewService(dev hardware.WaterDetector, scanInterval time.Duration***REMOVED*** *Service {
	return &Service{
		ScanInterval: scanInterval,
		Sensor:       dev,
***REMOVED***
***REMOVED***

func (t *Service***REMOVED*** Run(ctx context.Context, nc *nats.EncodedConn***REMOVED*** (err error***REMOVED*** {
	var reqSub *nats.Subscription
	var reqCh chan *nats.Msg

	reqCh = make(chan *nats.Msg***REMOVED***
	reqSub, err = nc.BindRecvChan("tank.meter", reqCh***REMOVED***
***REMOVED***
		return err
***REMOVED***
	defer func(***REMOVED*** {
		err = reqSub.Unsubscribe(***REMOVED***
		close(reqCh***REMOVED***
***REMOVED***(***REMOVED***

	var sensorErr error
	fullRecord := lib.FullRecord{
		IsFull: false,
		Time:   time.Time{***REMOVED***,
***REMOVED***

	timer := time.NewTimer(t.ScanInterval***REMOVED***

	for {
		select {
		case msg := <-reqCh:
			var resp lib.FullResponse
			if sensorErr != nil {
				resp = lib.FullResponse{
					Response: lib.Response{
						Code: 1,
						Msg:  sensorErr.Error(***REMOVED***,
				***REMOVED***,
			***REMOVED***
		***REMOVED*** else {
				resp = lib.FullResponse{
					Response: lib.Response{
						Code: 0,
				***REMOVED***,
					Payload: fullRecord,
			***REMOVED***
		***REMOVED***
			nc.Publish(msg.Reply, resp***REMOVED***
		case <-timer.C:
			if sensorErr = t.Sensor.Connect(***REMOVED***; sensorErr != nil {
				continue
		***REMOVED***
			fullRecord.IsFull = t.Sensor.IsWaterFull(***REMOVED***
			fullRecord.Time = time.Now(***REMOVED***
			timer = time.NewTimer(t.ScanInterval***REMOVED***
		case <-ctx.Done(***REMOVED***:
			err = ctx.Err(***REMOVED***
			return
	***REMOVED***
***REMOVED***
***REMOVED***

func GetMeterInfo(ctx context.Context, nc *nats.EncodedConn***REMOVED*** (resp lib.FullResponse, err error***REMOVED*** {
	payload, _ := jsoniter.Marshal(lib.Request{
		Code: lib.CodeGet,
***REMOVED******REMOVED***
	err = nc.RequestWithContext(ctx, "tank.meter", payload, &resp***REMOVED***
	return
***REMOVED***
