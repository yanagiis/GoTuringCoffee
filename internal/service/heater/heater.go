package heater

***REMOVED***
	"context"
	"errors"
	"time"

	jsoniter "github.com/json-iterator/go"
	nats "github.com/nats-io/go-nats"
	"github.com/yanagiis/GoTuringCoffee/internal/hardware"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
	"github.com/yanagiis/GoTuringCoffee/internal/service/tanktemp"
***REMOVED***

type Service struct {
	ScanInterval   time.Duration
	pwm            hardware.PWM
	pid            lib.PID
	targetTemp     float64
	sensorErr      error
	record         lib.HeaterRecord
	lastTempRecord lib.TempRecord
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

	reqCh := make(chan *nats.Msg***REMOVED***
	reqSub, err = nc.BindRecvChan("tank.heater", reqCh***REMOVED***
***REMOVED***
		return err
***REMOVED***
	defer func(***REMOVED*** {
		err = reqSub.Unsubscribe(***REMOVED***
		close(reqCh***REMOVED***
***REMOVED***(***REMOVED***

	timer := time.NewTimer(h.ScanInterval***REMOVED***

	for {
		select {
		case msg := <-reqCh:
			var req lib.HeaterRequest
			if err := jsoniter.Unmarshal(msg.Data, &req***REMOVED***; err != nil {
				nc.Publish(msg.Reply, lib.HeaterResponse{
					Response: lib.Response{
						Code: lib.CodeFailure,
						Msg:  err.Error(***REMOVED***,
				***REMOVED***,
			***REMOVED******REMOVED***
		***REMOVED***
			if req.IsGet(***REMOVED*** {
				resp := h.handleHeaterStatus(***REMOVED***
				nc.Publish(msg.Reply, resp***REMOVED***
		***REMOVED***
			if req.IsPut(***REMOVED*** {
				resp := h.handleSetTemperature(req.Temp***REMOVED***
				nc.Publish(msg.Reply, resp***REMOVED***
		***REMOVED***
		case <-timer.C:
			h.adjustTemperature(ctx, nc***REMOVED***
			timer = time.NewTimer(h.ScanInterval***REMOVED***
		case <-ctx.Done(***REMOVED***:
			err = ctx.Err(***REMOVED***
	***REMOVED***
***REMOVED***
***REMOVED***

func (h *Service***REMOVED*** handleHeaterStatus(***REMOVED*** lib.HeaterResponse {
	var resp lib.HeaterResponse
	if h.sensorErr != nil {
		resp = lib.HeaterResponse{
			Response: lib.Response{
				Code: lib.CodeFailure,
				Msg:  h.sensorErr.Error(***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED*** else {
		resp = lib.HeaterResponse{
			Response: lib.Response{
				Code: lib.CodeSuccess,
		***REMOVED***,
			Payload: h.record,
	***REMOVED***
***REMOVED***
	return resp
***REMOVED***

func (h *Service***REMOVED*** handleSetTemperature(temp float64***REMOVED*** lib.HeaterResponse {
	h.targetTemp = temp
	return lib.HeaterResponse{
		Response: lib.Response{
			Code: lib.CodeSuccess,
	***REMOVED***,
***REMOVED***
***REMOVED***

func (h *Service***REMOVED*** adjustTemperature(ctx context.Context, nc *nats.EncodedConn***REMOVED*** error {
	resp, err := tanktemp.GetTemperature(ctx, nc***REMOVED***
***REMOVED***
		return err
***REMOVED***
	if resp.IsFailure(***REMOVED*** {
		return errors.New("Cannot get tank temperature"***REMOVED***
***REMOVED***
	duty := h.pid.Compute(resp.Payload.Temp, resp.Payload.Time.Sub(h.lastTempRecord.Time***REMOVED******REMOVED***
	if err := h.pwm.PWM(duty, time.Second***REMOVED***; err != nil {
		return err
***REMOVED***
	h.record.Duty = duty
	h.record.Time = time.Now(***REMOVED***
	h.lastTempRecord = resp.Payload
	return nil
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
