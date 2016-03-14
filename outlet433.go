package main

import (
	// "github.com/boltdb/bolt"
	"fmt"
	"github.com/brutella/hc/model"
	"github.com/brutella/hc/model/accessory"
	"os"
	"os/exec"
)

func newHomeKitOutlet433(name string, systemCode string, unitCode string) (a *accessory.Switch) {

	binary := "./send"
	onArgs := []string{"-u", systemCode, unitCode, "1"}
	offArgs := []string{"-u", systemCode, unitCode, "0"}

	info := model.Info{
		Name:         name,
		Manufacturer: "Marco",
	}

	outlet433 := accessory.NewSwitch(info)

	// outlet433.SetOn(false)

	outlet433.OnStateChanged(func(on bool) {
		var err error
		var out []byte
		if on == true {
			cmd := exec.Command(binary, onArgs...)
			out, err = cmd.CombinedOutput()
			// outlet433.SetOn(true)
		} else {
			cmd := exec.Command(binary, offArgs...)
			out, err = cmd.CombinedOutput()
			// outlet433.SetOn(false)
		}
		fmt.Printf("==> Output: %s\n", string(out))

		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("433 Error: %s\n", err.Error()))
		}

	})

	return outlet433
}
