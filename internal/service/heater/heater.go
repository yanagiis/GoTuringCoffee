package heater

***REMOVED***
	"context"
	"errors"
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
	"github.com/yanagiis/GoTuringCoffee/internal/service/tanktemp"
***REMOVED***

type Service struct {
	ScanInterval time.Duration
	pwm          hardware.PWM
	pid          lib.PID
***REMOVED***

func NewService(dev hardware.PWM, scanInterval time.Duration, pid lib.PID***REMOVED*** *Service {
	return &Service{
		ScanInterval: scanInterval,
		pwm:          dev,
		pid:          pid,
***REMOVED***
***REMOVED***

func (h *Service***REMOVED*** Run(ctx context.Context, nc *nats.EncodedConn***REMOVED*** (err error***REMOVED*** {
	var reqSub *nats.Subscription
	var reqCh chan *nats.Msg

	reqCh = make(chan *nats.Msg***REMOVED***
	reqSub, err = nc.BindRecvChan("tank.heater", reqCh***REMOVED***
***REMOVED***
		return err
***REMOVED***
	defer func(***REMOVED*** {
		err = reqSub.Unsubscribe(***REMOVED***
		close(reqCh***REMOVED***
***REMOVED***(***REMOVED***

	var sensorErr error
	heaterRecord := lib.HeaterRecord{
		Duty:   0,
		Period: time.Duration(0***REMOVED***,
		Time:   time.Time{***REMOVED***,
***REMOVED***

	for {
		select {
		case msg := <-reqCh:
			var resp lib.HeaterResponse
			if sensorErr != nil {
				resp = lib.HeaterResponse{
					Response: lib.Response{
						Code: lib.CodeFailure,
						Msg:  sensorErr.Error(***REMOVED***,
				***REMOVED***,
			***REMOVED***
		***REMOVED*** else {
				resp = lib.HeaterResponse{
					Response: lib.Response{
						Code: lib.CodeSuccess,
				***REMOVED***,
					Payload: heaterRecord,
			***REMOVED***
		***REMOVED***
			nc.Publish(msg.Reply, resp***REMOVED***
		case <-ctx.Done(***REMOVED***:
			err = ctx.Err(***REMOVED***
			return
	***REMOVED***
***REMOVED***
***REMOVED***

func (h *Service***REMOVED*** scan(ctx context.Context, nc *nats.EncodedConn, out chan<- interface{***REMOVED******REMOVED*** {
	timer := time.NewTimer(h.ScanInterval***REMOVED***
	for {
		select {
		case <-ctx.Done(***REMOVED***:
			timer.Stop(***REMOVED***
			close(out***REMOVED***
			return
		case <-timer.C:
			resp, err := tanktemp.GetTemperature(ctx, nc***REMOVED***
			timer = time.NewTimer(h.ScanInterval***REMOVED***
		***REMOVED***
				out <- err
				continue
		***REMOVED***
			if resp.IsFailure(***REMOVED*** {
				out <- errors.New(resp.Msg.(string***REMOVED******REMOVED***
				continue
		***REMOVED***
	***REMOVED***
***REMOVED***
***REMOVED***

func GetHeaterInfo(ctx context.Context, nc *nats.EncodedConn***REMOVED*** (resp lib.HeaterResponse, err error***REMOVED*** {
	payload := lib.HeaterRequest{
		Request: lib.Request{
			Code: lib.CodePut,
	***REMOVED***,
***REMOVED***
	err = nc.RequestWithContext(ctx, "output.temperature", payload, &resp***REMOVED***
	return
***REMOVED***

func SetTemperature(ctx context.Context, nc *nats.EncodedConn, temp float64***REMOVED*** (resp lib.TempResponse, err error***REMOVED*** {
	payload := lib.HeaterRequest{
		Request: lib.Request{
			Code: lib.CodePut,
	***REMOVED***,
		Temp: temp,
***REMOVED***
	err = nc.RequestWithContext(ctx, "output.temperature", payload, &resp***REMOVED***
	return
***REMOVED***
