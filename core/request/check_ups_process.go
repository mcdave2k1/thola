// +build !client

package request

import (
	"context"
	"github.com/inexio/go-monitoringplugin"
	"github.com/inexio/thola/core/device"
	"github.com/inexio/thola/core/utility"
	"github.com/inexio/thola/core/value"
)

func (r *CheckUPSRequest) process(ctx context.Context) (Response, error) {
	r.init()

	readUPSResponse, err := r.getData(ctx)
	if r.mon.UpdateStatusOnError(err, monitoringplugin.UNKNOWN, "error while processing read ups request", true) {
		r.mon.PrintPerformanceData(false)
		return &CheckResponse{r.mon.GetInfo()}, nil
	}

	err = addCheckUPSPerformanceData(readUPSResponse.UPS, r.mon)
	if r.mon.UpdateStatusOnError(err, monitoringplugin.UNKNOWN, "error while adding performance data", true) {
		r.mon.PrintPerformanceData(false)
		return &CheckResponse{r.mon.GetInfo()}, nil
	}

	if readUPSResponse.UPS.MainsVoltageApplied != nil {
		r.mon.UpdateStatusIfNot(*readUPSResponse.UPS.MainsVoltageApplied, monitoringplugin.CRITICAL, "Mains voltage is not applied")
	}

	if !r.BatteryCurrentThresholds.isEmpty() {
		if readUPSResponse.UPS.BatteryCurrent == nil {
			r.mon.UpdateStatus(monitoringplugin.UNKNOWN, "battery current value is empty")
		} else if status := r.BatteryCurrentThresholds.checkValue(value.New(*readUPSResponse.UPS.BatteryCurrent)); status != monitoringplugin.OK {
			r.mon.UpdateStatus(status, "battery current is outside of threshold")
		}
	}

	if !r.BatteryTemperatureThresholds.isEmpty() {
		if readUPSResponse.UPS.BatteryTemperature == nil {
			r.mon.UpdateStatus(monitoringplugin.UNKNOWN, "battery temperature value is empty")
		} else if status := r.BatteryTemperatureThresholds.checkValue(value.New(*readUPSResponse.UPS.BatteryTemperature)); status != monitoringplugin.OK {
			r.mon.UpdateStatus(status, "battery temperature is outside of threshold")
		}
	}

	if !r.CurrentLoadThresholds.isEmpty() {
		if readUPSResponse.UPS.CurrentLoad == nil {
			r.mon.UpdateStatus(monitoringplugin.UNKNOWN, "current load value is empty")
		} else if status := r.CurrentLoadThresholds.checkValue(value.New(*readUPSResponse.UPS.CurrentLoad)); status != monitoringplugin.OK {
			r.mon.UpdateStatus(status, "current load is outside of threshold")
		}
	}

	if !r.RectifierCurrentThresholds.isEmpty() {
		if readUPSResponse.UPS.RectifierCurrent == nil {
			r.mon.UpdateStatus(monitoringplugin.UNKNOWN, "rectifier current value is empty")
		} else if status := r.RectifierCurrentThresholds.checkValue(value.New(*readUPSResponse.UPS.RectifierCurrent)); status != monitoringplugin.OK {
			r.mon.UpdateStatus(status, "rectifier current is outside of threshold")
		}
	}

	if !r.SystemVoltageThresholds.isEmpty() {
		if readUPSResponse.UPS.SystemVoltage == nil {
			r.mon.UpdateStatus(monitoringplugin.UNKNOWN, "system voltage value is empty")
		} else if status := r.SystemVoltageThresholds.checkValue(value.New(*readUPSResponse.UPS.SystemVoltage)); status != monitoringplugin.OK {
			r.mon.UpdateStatus(status, "system voltage is outside of threshold")
		}
	}

	return &CheckResponse{r.mon.GetInfo()}, nil
}

func (r *CheckUPSRequest) getData(ctx context.Context) (*ReadUPSResponse, error) {
	readUPSRequest := ReadUPSRequest{ReadRequest{r.BaseRequest}}
	response, err := readUPSRequest.process(ctx)
	if err != nil {
		return nil, err
	}

	readUPSResponse := response.(*ReadUPSResponse)
	return readUPSResponse, nil
}

func addCheckUPSPerformanceData(ups device.UPSComponent, r *monitoringplugin.Response) error {
	if ups.AlarmLowVoltageDisconnect != nil {
		err := r.AddPerformanceDataPoint(monitoringplugin.NewPerformanceDataPoint("alarm_low_voltage_disconnect", *ups.AlarmLowVoltageDisconnect, ""))
		if err != nil {
			return err
		}
	}

	if ups.BatteryAmperage != nil {
		err := r.AddPerformanceDataPoint(monitoringplugin.NewPerformanceDataPoint("batt_amperage", *ups.BatteryAmperage, ""))
		if err != nil {
			return err
		}
	}

	if ups.BatteryRemainingTime != nil {
		err := r.AddPerformanceDataPoint(monitoringplugin.NewPerformanceDataPoint("batt_remaining_time", *ups.BatteryRemainingTime, ""))
		if err != nil {
			return err
		}
	}

	if ups.BatteryCapacity != nil {
		err := r.AddPerformanceDataPoint(monitoringplugin.NewPerformanceDataPoint("batt_capacity", *ups.BatteryCapacity, ""))
		if err != nil {
			return err
		}
	}

	if ups.BatteryCurrent != nil {
		err := r.AddPerformanceDataPoint(monitoringplugin.NewPerformanceDataPoint("batt_current", *ups.BatteryCurrent, ""))
		if err != nil {
			return err
		}
	}

	if ups.BatteryTemperature != nil {
		err := r.AddPerformanceDataPoint(monitoringplugin.NewPerformanceDataPoint("batt_temperature", *ups.BatteryTemperature, ""))
		if err != nil {
			return err
		}
	}

	if ups.BatteryVoltage != nil {
		err := r.AddPerformanceDataPoint(monitoringplugin.NewPerformanceDataPoint("batt_voltage", *ups.BatteryVoltage, ""))
		if err != nil {
			return err
		}
	}

	if ups.CurrentLoad != nil {
		err := r.AddPerformanceDataPoint(monitoringplugin.NewPerformanceDataPoint("current_load", *ups.CurrentLoad, ""))
		if err != nil {
			return err
		}
	}

	if ups.MainsVoltageApplied != nil {
		err := r.AddPerformanceDataPoint(monitoringplugin.NewPerformanceDataPoint("mains_voltage_applied", utility.IfThenElse(*ups.MainsVoltageApplied, 1, 0), ""))
		if err != nil {
			return err
		}
	}

	if ups.RectifierCurrent != nil {
		err := r.AddPerformanceDataPoint(monitoringplugin.NewPerformanceDataPoint("rectifier_current", *ups.RectifierCurrent, ""))
		if err != nil {
			return err
		}
	}

	if ups.SystemVoltage != nil {
		err := r.AddPerformanceDataPoint(monitoringplugin.NewPerformanceDataPoint("sys_voltage", *ups.SystemVoltage, ""))
		if err != nil {
			return err
		}
	}
	return nil
}
