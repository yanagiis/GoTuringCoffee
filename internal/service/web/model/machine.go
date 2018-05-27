package model

import (
	"context"
	"errors"

	nats "github.com/nats-io/go-nats"
	"github.com/yanagiis/GoTuringCoffee/internal/service/heater"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
	"github.com/yanagiis/GoTuringCoffee/internal/service/outtemp"
	"github.com/yanagiis/GoTuringCoffee/internal/service/replenisher"
	"github.com/yanagiis/GoTuringCoffee/internal/service/tankmeter"
	"github.com/yanagiis/GoTuringCoffee/internal/service/tanktemp"
)

type Machine struct {
	nc  *nats.EncodedConn
	ctx context.Context
}

type MachineStatus struct {
	Heater    lib.HeaterRecord      `json:"heater"`
	Output    lib.TempRecord        `json:"output"`
	Replenish lib.ReplenisherRecord `json:"replenish"`
	TankMeter lib.FullRecord        `json:"tankmeter"`
	TankTemp  lib.TempRecord        `json:"tanktemp"`
}

func NewMachine(ctx context.Context, nc *nats.EncodedConn) *Machine {
	return &Machine{
		nc:  nc,
		ctx: ctx,
	}
}

func (m *Machine) GetMachineStatus() (status MachineStatus, err error) {
	var heaterResp lib.HeaterResponse
	var outResp lib.TempResponse
	var replenResp lib.ReplenisherResponse
	var tankMeterResp lib.FullResponse
	var tankTempResp lib.TempResponse

	if heaterResp, err = heater.GetHeaterInfo(m.ctx, m.nc); err != nil {
		return
	}
	if heaterResp.IsFailure() {
		err = errors.New(heaterResp.Msg)
		return
	}
	status.Heater = heaterResp.Payload

	if outResp, err = outtemp.GetTemperature(m.ctx, m.nc); err != nil {
		return
	}
	if outResp.IsFailure() {
		err = errors.New(outResp.Msg)
		return
	}
	status.Output = outResp.Payload

	if replenResp, err = replenisher.GetReplenishInfo(m.ctx, m.nc); err != nil {
		return
	}
	if replenResp.IsFailure() {
		err = errors.New(replenResp.Msg)
		return
	}
	status.Replenish = replenResp.Payload

	if tankMeterResp, err = tankmeter.GetMeterInfo(m.ctx, m.nc); err != nil {
		return
	}
	if tankMeterResp.IsFailure() {
		err = errors.New(tankMeterResp.Msg)
		return
	}
	status.TankMeter = tankMeterResp.Payload

	if tankTempResp, err = tanktemp.GetTemperature(m.ctx, m.nc); err != nil {
		return
	}
	if tankTempResp.IsFailure() {
		err = errors.New(tankTempResp.Msg)
		return
	}
	status.TankTemp = tankTempResp.Payload
	return
}
