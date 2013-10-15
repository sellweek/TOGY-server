package configuration

import (
	"appengine"
	"encoding/json"
	"models/configuration/config"
	"models/configuration/timeConfig"
)

const (
	jsonTimeFormat = "15:04"
	jsonDateFormat = "2006-1-2"
)

//JSON returns JSON representation of config and time override settings
//for use in clients.
func JSON(c appengine.Context) (js []byte, err error) {
	j := make(map[string]interface{})
	conf, err := config.Get(c)
	if err != nil {
		return
	}

	tcs, err := timeConfig.GetAll(c)
	if err != nil {
		return
	}

	j["StandardTimeSettings"] = map[string]string{
		"TurnOn":  conf.StandardOn.Format(jsonTimeFormat),
		"TurnOff": conf.StandardOff.Format(jsonTimeFormat),
	}
	j["UpdateInterval"] = conf.UpdateInterval
	switch conf.OverrideState {
	case config.NoOverride:
		j["OverrideOn"] = false
		j["OverrideOff"] = false
	case config.OverrideOn:
		j["OverrideOn"] = true
		j["OverrideOff"] = false
	case config.OverrideOff:
		j["OverrideOn"] = false
		j["OverrideOff"] = true
	}

	j["Weekends"] = conf.Weekends

	for _, tc := range tcs {
		timeMap := make(map[string]string)
		timeMap["TurnOn"] = tc.On.Format(jsonTimeFormat)
		timeMap["TurnOff"] = tc.Off.Format(jsonTimeFormat)

		j["OverrideDays"] = make(map[string]map[string]string)
		j["OverrideDays"].(map[string]map[string]string)[tc.Date.Format(jsonDateFormat)] = timeMap
	}

	js, err = json.Marshal(j)
	return
}
