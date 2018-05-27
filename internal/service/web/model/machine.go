package model

***REMOVED***
	"context"
	"errors"

	nats "github.com/nats-io/go-nats"
	"github.com/yanagiis/GoTuringCoffee/internal/service/heater"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
	"github.com/yanagiis/GoTuringCoffee/internal/service/outtemp"
	"github.com/yanagiis/GoTuringCoffee/internal/service/replenisher"
	"github.com/yanagiis/GoTuringCoffee/internal/service/tankmeter"
	"github.com/yanagiis/GoTuringCoffee/internal/service/tanktemp"
***REMOVED***

type Machine struct {
	nc  *nats.EncodedConn
	ctx context.Context
***REMOVED***

type MachineStatus struct {
	Heater    lib.HeaterRecord      `json:"heater"`
	Output    lib.TempRecord        `json:"output"`
	Replenish lib.ReplenisherRecord `json:"replenish"`
	TankMeter lib.FullRecord        `json:"tankmeter"`
	TankTemp  lib.TempRecord        `json:"tanktemp"`
***REMOVED***

func NewMachine(ctx context.Context, nc *nats.EncodedConn***REMOVED*** *Machine {
	return &Machine{
		nc:  nc,
		ctx: ctx,
***REMOVED***
***REMOVED***

func (m *Machine***REMOVED*** GetMachineStatus(***REMOVED*** (status MachineStatus, err error***REMOVED*** {
	var heaterResp lib.HeaterResponse
	var outResp lib.TempResponse
	var replenResp lib.ReplenisherResponse
	var tankMeterResp lib.FullResponse
	var tankTempResp lib.TempResponse

	if heaterResp, err = heater.GetHeaterInfo(m.ctx, m.nc***REMOVED***; err != nil {
		return
***REMOVED***
	if heaterResp.IsFailure(***REMOVED*** {
		err = errors.New(heaterResp.Msg***REMOVED***
		return
***REMOVED***
	status.Heater = heaterResp.Payload

	if outResp, err = outtemp.GetTemperature(m.ctx, m.nc***REMOVED***; err != nil {
		return
***REMOVED***
	if outResp.IsFailure(***REMOVED*** {
		err = errors.New(outResp.Msg***REMOVED***
		return
***REMOVED***
	status.Output = outResp.Payload

	if replenResp, err = replenisher.GetReplenishInfo(m.ctx, m.nc***REMOVED***; err != nil {
		return
***REMOVED***
	if replenResp.IsFailure(***REMOVED*** {
		err = errors.New(replenResp.Msg***REMOVED***
		return
***REMOVED***
	status.Replenish = replenResp.Payload

	if tankMeterResp, err = tankmeter.GetMeterInfo(m.ctx, m.nc***REMOVED***; err != nil {
		return
***REMOVED***
	if tankMeterResp.IsFailure(***REMOVED*** {
		err = errors.New(tankMeterResp.Msg***REMOVED***
		return
***REMOVED***
	status.TankMeter = tankMeterResp.Payload

	if tankTempResp, err = tanktemp.GetTemperature(m.ctx, m.nc***REMOVED***; err != nil {
		return
***REMOVED***
	if tankTempResp.IsFailure(***REMOVED*** {
		err = errors.New(tankTempResp.Msg***REMOVED***
		return
***REMOVED***
	status.TankTemp = tankTempResp.Payload
	return
***REMOVED***
