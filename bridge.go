package main

import (
	"github.com/brutella/hc/model"
	"github.com/brutella/hc/model/accessory"
)

func newHomeKitBridge(name string) (a *accessory.Accessory) {
	info := model.Info{
		Name:         name,
		Manufacturer: "Marco",
	}

	bridge := accessory.New(info, 2)

	return bridge
}
