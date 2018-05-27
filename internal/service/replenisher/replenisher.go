package replenisher

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
	Dev          hardware.PWM
	PWMConf      hardware.PWMConfig
	devErr       error
	stop         bool
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

	r.Dev.Connect(***REMOVED***
	defer r.Dev.Disconnect(***REMOVED***
	timer := time.NewTimer(r.ScanInterval***REMOVED***

	for {
		select {
		case msg := <-reqCh:
			var req lib.ReplenisherRequest
			if decodeErr := jsoniter.Unmarshal(msg.Data, &req***REMOVED***; decodeErr != nil {
				nc.Publish(msg.Reply, lib.HeaterResponse{
					Response: lib.Response{
						Code: lib.CodeFailure,
						Msg:  decodeErr.Error(***REMOVED***,
				***REMOVED***,
			***REMOVED******REMOVED***
		***REMOVED***
			if req.IsGet(***REMOVED*** {
				resp := r.handleReplenishStatus(***REMOVED***
				nc.Publish(msg.Reply, resp***REMOVED***
		***REMOVED***
			if req.IsPut(***REMOVED*** {
				resp := r.handleControlReplenish(req.Stop***REMOVED***
				nc.Publish(msg.Reply, resp***REMOVED***
		***REMOVED***
		case <-timer.C:
			timer = time.NewTimer(r.ScanInterval***REMOVED***
			r.scan(***REMOVED***
		case <-ctx.Done(***REMOVED***:
			err = ctx.Err(***REMOVED***
			return
	***REMOVED***
***REMOVED***
***REMOVED***

func (r *Service***REMOVED*** handleReplenishStatus(***REMOVED*** lib.ReplenisherResponse {
	if r.devErr != nil {
		return lib.ReplenisherResponse{
			Response: lib.Response{
				Code: lib.CodeFailure,
				Msg:  r.devErr.Error(***REMOVED***,
		***REMOVED***,
			Payload: lib.ReplenisherRecord{***REMOVED***,
	***REMOVED***
***REMOVED*** else {
		return lib.ReplenisherResponse{
			Response: lib.Response{
				Code: lib.CodeSuccess,
		***REMOVED***,
			Payload: lib.ReplenisherRecord{
				IsReplenishing: !r.stop,
				Time:           time.Now(***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***
***REMOVED***

func (r *Service***REMOVED*** handleControlReplenish(stop bool***REMOVED*** lib.ReplenisherResponse {
	r.stop = stop
	return lib.ReplenisherResponse{
		Response: lib.Response{
			Code: lib.CodeSuccess,
	***REMOVED***,
***REMOVED***
***REMOVED***

func (r *Service***REMOVED*** scan(***REMOVED*** {
	duty := r.PWMConf.Duty
	if r.stop {
		duty = 0
***REMOVED***
	r.Dev.PWM(duty, r.PWMConf.Period***REMOVED***
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
