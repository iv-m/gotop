// +build linux

package devices

import (
	"strings"

	psHost "github.com/shirou/gopsutil/host"
)

func init() {
	RegisterTemp(getTemps)
	RegisterDeviceList(Temperatures, devs)
}

func getTemps(temps map[string]int) map[string]error {
	sensors, err := psHost.SensorsTemperatures()
	if err != nil {
		return map[string]error{"psHost": err}
	}
	for _, sensor := range sensors {
		// removes '_input' from the end of the sensor name
		idx := strings.Index(sensor.SensorKey, "_input")
		if idx >= 0 {
			label := sensor.SensorKey[:idx]
			if _, ok := temps[label]; ok {
				temps[label] = int(sensor.Temperature)
			}
		}
	}
	return nil
}

func devs() []string {
	sensors, err := psHost.SensorsTemperatures()
	if err != nil {
		return []string{}
	}
	rv := make([]string, 0, len(sensors))
	for _, sensor := range sensors {
		// only sensors with input in their name are giving us live temp info
		if strings.Contains(sensor.SensorKey, "input") && sensor.Temperature != 0 {
			// removes '_input' from the end of the sensor name
			label := sensor.SensorKey[:strings.Index(sensor.SensorKey, "_input")]
			rv = append(rv, label)
		}
	}
	return rv
}
