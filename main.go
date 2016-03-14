package main

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/brutella/hc/hap"
	"github.com/brutella/hc/model/accessory"
	"log"
	"time"
)

func main() {

	host := flag.String("host", "localhost", "OpenHAB Host")
	port := flag.String("port", "8080", "OpenHAB Port")
	sitemap := flag.String("sitemap", "default", "OpenHAB Sitemap")
	bridgename := flag.String("name", "LBI_Bridge", "HomeKit Bridge Name")
	bridgepin := flag.String("pin", "32191123", "HomeKit Bridge PIN")

	flag.Parse()

	value, _ := querySitemap(*host, *port, *sitemap)
	getThermostats(value)

	go func() {

		for {

			value, _ := querySitemap(*host, *port, *sitemap)
			getThermostats(value)
			fmt.Println("Querying OpenHAB")
			time.Sleep(45 * time.Second)
		}
	}()

	// get all bucket names into array and: for each array element do newHomeKitThermostat giving name as argument

	var thermoNames []string
	var hkObjects []*accessory.Accessory
	// var hkObjects []*accessory.Thermostat

	// open bolt db
	db, err := bolt.Open("hk.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			thermoNames = append(thermoNames, string(name))
			return nil
		})
	})

	db.Close()

	if err != nil {
		log.Fatal(err)
	}

	for i := range thermoNames {
		x := newHomeKitThermostat(thermoNames[i]).Accessory
		hkObjects = append(hkObjects, x)
	}

	schrank := newHomeKitOutlet433("Schrank", "01111", "3").Accessory
	hkObjects = append(hkObjects, schrank)

	t, err := hap.NewIPTransport(hap.Config{Pin: *bridgepin}, newHomeKitBridge(*bridgename), hkObjects...)
	if err != nil {
		log.Fatal(err)
	}

	hap.OnTermination(func() {
		t.Stop()
	})

	t.Start()
}
