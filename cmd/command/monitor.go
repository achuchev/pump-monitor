package command

import (
	"fmt"
	"time"

	c "github.com/achuchev/pump-monitor/cmd/common"
	"github.com/achuchev/pump-monitor/cmd/mail"
	"github.com/achuchev/pump-monitor/cmd/pump"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
}

var (
	debugCmd = &cobra.Command{
		Use:   "monitor",
		Short: "Starts monitoring of the pump",
		Long:  `...`,
		RunE:  doMonitor,
	}
)

func doMonitor(cmd *cobra.Command, args []string) error {
	ticker := time.NewTicker(time.Duration(c.FlagsConfig.MonitorEveryMinutes) * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				go func() {
					log.Info("")
					monitor()
					log.Info("")
				}()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	log.Info("Started")

	log.Info("Press Ctrl+C to stop")
	monitor()
	for {
		time.Sleep(10 * time.Second)
		fmt.Print(".")
	}
}
func monitor() {
	log.Info("Checking the status of the pump...")

	err := pump.Authenticate()
	if err != nil {
		log.Fatalf("authentication failed: %v", err)
	}

	result, err := pump.GetLiveData(pump.DeviceID, []string{pump.CurrentBufferStorageTemperatureSection, pump.TargetBufferStorageTemperatureSection, pump.CurrentElectricMeterWattsSection, pump.CurrentCompressorPercentageSection},
		[]string{pump.CurrentBufferStorageTemperaturePosition, pump.TargetBufferStorageTemperaturePosition, pump.CurrentElectricMeterWattsPosition, pump.CurrentCompressorPercentagePosition})
	if err != nil {
		log.Fatalf("failed to get live data: %v", err)
	}

	var currentBufferStorageTemperature float32
	var targetBufferStorageTemperature float32
	var currentElectricMeterWatts float32
	var currentCompressorPercentage float32

	for i, v := range result.Data {
		position := result.Positions[i]
		if position == pump.CurrentBufferStorageTemperaturePosition {
			log.Infof(" Current buffer storage temperature: %v\u00B0C\n", v)
			currentBufferStorageTemperature = v
		}
		if position == pump.TargetBufferStorageTemperaturePosition {
			log.Infof(" Target buffer storage temperature: %v\u00B0C\n", v)
			targetBufferStorageTemperature = v
		}
		if position == pump.CurrentElectricMeterWattsPosition {
			log.Infof(" Current electric meter: %vW\n", v)
			currentElectricMeterWatts = v
		}
		if position == pump.CurrentCompressorPercentagePosition {
			log.Infof(" Current compressor: %v%%\n", v)
			currentCompressorPercentage = v
		}
	}
	force := false
	if force || ((targetBufferStorageTemperature-currentBufferStorageTemperature) > float32(c.FlagsConfig.MonitorTempDifferenceThreshold)) && ((currentElectricMeterWatts == 0) || (currentCompressorPercentage == 0)) {
		err := mail.Notify("Heat Pump Alert", fmt.Sprintf(`Seems like the heat pump failed. Consider restart.

Buffer storage temperature is %v, target is %v.
Current electric meter (W): %v
Current compressor (%%): %v`, currentBufferStorageTemperature, targetBufferStorageTemperature, currentElectricMeterWatts, currentCompressorPercentage))

		if err != nil {
			log.Errorf("failed to send mail: %v", err)
		}
	} else {
		log.Info("Status OK")
	}
}
