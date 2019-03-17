package model

import (
	"context"
	"time"

	"GoTuringCoffee/internal/service/distance"
	"GoTuringCoffee/internal/service/heater"
	"GoTuringCoffee/internal/service/lib"
	"GoTuringCoffee/internal/service/outtemp"
	"GoTuringCoffee/internal/service/tankmeter"
	"GoTuringCoffee/internal/service/tanktemp"

	nats "github.com/nats-io/go-nats"
	"github.com/rs/zerolog/log"
)

type Machine struct {
	nc  *nats.EncodedConn
	ctx context.Context
}

// type MachineStatus struct {
// 	Heater    lib.HeaterRecord      `json:"heater"`
// 	Output    lib.TempRecord        `json:"output"`
// 	Replenish lib.ReplenisherRecord `json:"replenish"`
// 	TankMeter lib.FullRecord        `json:"tankmeter"`
// 	TankTemp  lib.TempRecord        `json:"tanktemp"`
// }

type HeaterStatus struct {
	DutyCycle  *float64 `json:"duty_cycle"`
	TargetTemp *float64 `json:"target_temperature"`
}

type MachineStatus struct {
	Output     *float64      `json:"output_temperature"`
	TankTemp   *float64      `json:"tank_temperature"`
	Heater     *HeaterStatus `json:"heater_status"`
	WaterLevel *bool         `json:"water_level"`
	Distance   *int64        `json:"distance"`
}

// type MachineStatus struct {
// 	Heater    lib.HeaterRecord      `json:"heater"`
// 	Output    lib.TempRecord        `json:"output"`
// 	Replenish lib.ReplenisherRecord `json:"replenish"`
// 	TankMeter lib.FullRecord        `json:"tankmeter"`
// 	TankTemp  lib.TempRecord        `json:"tanktemp"`
// }

func NewMachine(ctx context.Context, nc *nats.EncodedConn) *Machine {
	return &Machine{
		nc:  nc,
		ctx: ctx,
	}
}

func (m *Machine) GetMachineStatus() (status MachineStatus, err error) {
	var heaterResp lib.HeaterResponse
	var outResp lib.TempResponse
	// var replenResp lib.ReplenisherResponse
	var tankMeterResp lib.FullResponse
	var tankTempResp lib.TempResponse
	var ctx context.Context
	var cancel context.CancelFunc
	var distanceResp lib.DistanceResponse

	ctx, cancel = context.WithTimeout(m.ctx, time.Second)
	if heaterResp, err = heater.GetHeaterInfo(ctx, m.nc); err != nil {
		cancel()
	}
	if !heaterResp.IsFailure() {
		log.Info().Msgf("Heater: %v\n", heaterResp)

		status.Heater = new(HeaterStatus)
		status.Heater.DutyCycle = new(float64)
		*status.Heater.DutyCycle = heaterResp.Payload.Duty
		status.Heater.TargetTemp = new(float64)
		*status.Heater.TargetTemp = heaterResp.Payload.Target
	}

	ctx, cancel = context.WithTimeout(m.ctx, time.Second)
	if outResp, err = outtemp.GetTemperature(ctx, m.nc); err != nil {
		cancel()
	}
	if !outResp.IsFailure() {
		log.Info().Msgf("Output: %v\n", outResp)
		status.Output = new(float64)
		*status.Output = outResp.Payload.Temp
	}

	ctx, cancel = context.WithTimeout(m.ctx, time.Second)
	if distanceResp, err = distance.GetDistance(ctx, m.nc); err != nil {
		cancel()
	}
	if !distanceResp.IsFailure() {
		log.Info().Msgf("Output: %v\n", distanceResp)
		status.Distance = new(int64)
		*status.Distance = int64(distanceResp.Payload.Distance)
	}

	// 	ctx, cancel = context.WithTimeout(m.ctx, time.Second)
	// 	if replenResp, err = replenisher.GetReplenishInfo(ctx, m.nc); err != nil {
	// 		cancel()
	// 	}
	// 	if !replenResp.IsFailure() {
	// 		log.Info().Msgf("Replenish: %v\n", replenResp)
	// 		status.Replenish = replenResp.Payload
	// 	}

	ctx, cancel = context.WithTimeout(m.ctx, time.Second)
	defer cancel()

	tankMeterResp, err = tankmeter.GetMeterInfo(ctx, m.nc)
	if !tankMeterResp.IsFailure() {
		log.Info().Msgf("TankMeter: %v\n", tankMeterResp)
		status.WaterLevel = new(bool)
		*status.WaterLevel = tankMeterResp.Payload.IsFull
	}

	ctx, cancel = context.WithTimeout(m.ctx, time.Second)
	tankTempResp, err = tanktemp.GetTemperature(ctx, m.nc)
	defer cancel()

	if !tankTempResp.IsFailure() {
		log.Info().Msgf("TankTemp: %v\n", tankTempResp)
		status.TankTemp = new(float64)
		*status.TankTemp = tankTempResp.Payload.Temp
	}
	return
}
