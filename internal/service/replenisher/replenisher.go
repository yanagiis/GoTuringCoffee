package replenisher

***REMOVED***
	"context"
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
***REMOVED***

type Service struct {
	ScanInterval time.Duration
	Dev          hardware.PWM
	PWMConf      hardware.PWMConfig
***REMOVED***

func NewService(dev hardware.PWM, scanInterval time.Duration, pwmConf hardware.PWMConfig***REMOVED*** *Service {
	return &Service{
		ScanInterval: scanInterval,
		Dev:          dev,
		PWMConf:      pwmConf,
***REMOVED***
***REMOVED***

func (r *Service***REMOVED*** Run(ctx context.Context, nc *nats.EncodedConn***REMOVED*** (err error***REMOVED*** {
	var reqSub *nats.Subscription
	var reqCh chan *nats.Msg

	reqCh = make(chan *nats.Msg***REMOVED***
	reqSub, err = nc.BindRecvChan("tank.replenisher", reqCh***REMOVED***
***REMOVED***
		return err
***REMOVED***
	defer func(***REMOVED*** {
		err = reqSub.Unsubscribe(***REMOVED***
		close(reqCh***REMOVED***
***REMOVED***(***REMOVED***

	var devErr error
	replenishRecord := lib.ReplenisherRecord{
		IsReplenishing: false,
		Time:           time.Time{***REMOVED***,
***REMOVED***

	r.Dev.Connect(***REMOVED***
	defer r.Dev.Disconnect(***REMOVED***
	timer := time.NewTimer(r.ScanInterval***REMOVED***

	for {
		select {
		case msg := <-reqCh:
			var resp lib.ReplenisherResponse
			if devErr != nil {
				resp = lib.ReplenisherResponse{
					Response: lib.Response{
						Code: lib.CodeFailure,
						Msg:  devErr.Error(***REMOVED***,
				***REMOVED***,
					Payload: lib.ReplenisherRecord{***REMOVED***,
			***REMOVED***
		***REMOVED*** else {
				resp = lib.ReplenisherResponse{
					Response: lib.Response{
						Code: lib.CodeSuccess,
						Msg:  "",
				***REMOVED***,
					Payload: replenishRecord,
			***REMOVED***
		***REMOVED***
			nc.Publish(msg.Reply, resp***REMOVED***
		case <-timer.C:
			timer = time.NewTimer(r.ScanInterval***REMOVED***
		case <-ctx.Done(***REMOVED***:
			err = ctx.Err(***REMOVED***
			return
	***REMOVED***
***REMOVED***
***REMOVED***

func GetReplenishInfo(ctx context.Context, nc *nats.EncodedConn***REMOVED*** (resp lib.ReplenisherResponse, err error***REMOVED*** {
	payload := lib.ReplenisherRequest{
		Request: lib.Request{
			Code: lib.CodeGet,
	***REMOVED***,
***REMOVED***
	err = nc.RequestWithContext(ctx, "tank.meter", payload, &resp***REMOVED***
***REMOVED***
		return
***REMOVED***
	return
***REMOVED***

func StopReplenish(ctx context.Context, nc *nats.EncodedConn***REMOVED*** (lib.ReplenisherResponse, error***REMOVED*** {
	return toggleReplenish(ctx, nc, true***REMOVED***
***REMOVED***

func StartReplenish(ctx context.Context, nc *nats.EncodedConn***REMOVED*** (lib.ReplenisherResponse, error***REMOVED*** {
	return toggleReplenish(ctx, nc, false***REMOVED***
***REMOVED***

func toggleReplenish(ctx context.Context, nc *nats.EncodedConn, stop bool***REMOVED*** (resp lib.ReplenisherResponse, err error***REMOVED*** {
	payload := lib.ReplenisherRequest{
		Request: lib.Request{
			Code: lib.CodeGet,
	***REMOVED***,
		Stop: stop,
***REMOVED***
	err = nc.RequestWithContext(ctx, "tank.meter", payload, &resp***REMOVED***
	return
***REMOVED***
