package otii

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func EnableMainPower() error {
	log.Println("Enabling main power")
	if err := Run("otii.create_project():enable_main_power(true)"); err != nil {
		return err
	}
	log.Println("Waiting ten seconds after enabling main power to let device boot up.")
	time.Sleep(10 * time.Second)
	return nil
}

func DisableMainPower() error {
	return Run("otii.create_project():enable_main_power(false)")
}

func Calibrate() error {
	log.Println("Calibrating...")
	return Run(`
		local devices = otii.get_devices("Arc")
		assert(#devices > 0, "No available devices")
		local box = otii.open_device(devices[1].id)
		assert(box ~= nil, "No available otii")
		box:calibrate()
	`)
}

func Record(duration time.Duration) error {
	log.Println("Recording started")
	err := Run(strings.Replace(recordScript, "DURATION", strconv.Itoa(int(duration/time.Millisecond)), -1))
	log.Println("Recording complete")
	return err
}

func Run(script string) error {
	f, err := ioutil.TempFile("", "otii-script.lua")
	if err != nil {
		log.Println("Error opening temporary file:", err)
		return err
	}
	defer os.Remove(f.Name())

	if _, err := f.WriteString(script); err != nil {
		log.Println("Error writing script:", err)
		return err
	}
	if err := f.Close(); err != nil {
		log.Println("Error closing script:", err)
		return err
	}

	scriptPath, err := filepath.Abs(f.Name())
	if err != nil {
		log.Println("Error abs path:", err)
		return err
	}
	out, err := exec.Command("/Applications/otii.app/Contents/MacOS/otiicli", "--no-banner", scriptPath).CombinedOutput()
	if err != nil {
		log.Printf("Error running otii script: %v\n%s", err, out)
		return err
	}

	if len(out) > 0 {
		log.Println("Otii script output:", string(out))
	}
	return nil
}

const recordScript = `
local project = otii.create_project()
assert(project ~= nil, "Cannot create project")

local devices = otii.get_devices("Arc")
assert(#devices > 0, "No available devices")
local box = otii.open_device(devices[1].id)
assert(box ~= nil, "No available otii")

box:set_main_voltage(3.3)
box:set_range("high")
box:set_max_current(0.5)
box:enable_channel("mc", true)
box:enable_channel("mv", true)

project:start()
otii.msleep(DURATION)
project:stop()

local filename = string.format("captures/capture_%s.otii", os.date("%Y-%m-%dT%H:%M:%S"))
project:save(filename)
project:close()
box:close()
`
