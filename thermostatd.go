package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/brutella/hc/model"
	"github.com/brutella/hc/model/accessory"
	"log"
	"strconv"
	"time"
)

func newHomeKitThermostat(bucketName string) (a *accessory.Thermostat) {

	var hum OpenHabItem
	var settemp OpenHabItem
	var temp OpenHabItem
	var mode OpenHabItem

	info := model.Info{
		Name:         bucketName,
		Manufacturer: "Marco",
	}

	thermostat := accessory.NewThermostat(info, 12, 0, 23, 0.5)

	thermostat.OnTargetTempChange(func(targetTemp float64) {
		setValue("localhost", "8080", settemp.Name, strconv.FormatFloat(targetTemp, 'f', 6, 64))
	})

	thermostat.OnTargetModeChange(func(targetMode model.HeatCoolModeType) {
		switch targetMode {
		case model.HeatCoolModeAuto:
			setValue("localhost", "8080", mode.Name, "ON")
		case model.HeatCoolModeOff:
			setValue("localhost", "8080", mode.Name, "OFF")
		// My system doesn't support this, so default to auto mode
		case model.HeatCoolModeHeat, model.HeatCoolModeCool:
			thermostat.SetTargetMode(model.HeatCoolModeAuto)
		}
	})

	go func() {

		for {
			// open bolt db
			db, err := bolt.Open("hk.db", 0600, nil)
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()

			db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(bucketName))
				h := b.Get([]byte("hum"))
				s := b.Get([]byte("settemp"))
				t := b.Get([]byte("temp"))
				m := b.Get([]byte("mode"))
				json.Unmarshal(h, &hum)
				json.Unmarshal(s, &settemp)
				json.Unmarshal(t, &temp)
				json.Unmarshal(m, &mode)
				return nil
			})

			db.Close()

			thermostat.SetTemperature(getFloat(temp.State))
			thermostat.SetTargetTemperature(getFloat(settemp.State))

			if mode.State == "ON" {
				thermostat.SetTargetMode(model.HeatCoolModeAuto)
			} else {
				thermostat.SetTargetMode(model.HeatCoolModeOff)
			}

			time.Sleep(30 * time.Second)
		}
	}()

	return thermostat
}

func getFloat(sValue string) (fValue float64) {
	fValue, _ = strconv.ParseFloat(sValue, 64)

	return fValue
}
