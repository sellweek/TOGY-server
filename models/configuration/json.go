package configuration

import (
	"appengine"
	"encoding/json"
	"models/configuration/config"
	"models/configuration/timeConfig"
)

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
		"TurnOn":  conf.StandardOn.Format(config.ConfTimeFormat),
		"TurnOff": conf.StandardOff.Format(config.ConfTimeFormat),
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
	for _, tc := range tcs {
		timeMap := make(map[string]string)
		timeMap["TurnOn"] = tc.On.Format(config.ConfTimeFormat)
		timeMap["TurnOff"] = tc.Off.Format(config.ConfTimeFormat)

		j["OverrideDays"] = make(map[string]map[string]string)
		j["OverrideDays"].(map[string]map[string]string)[tc.Date.Format(config.ConfDateFormat)] = timeMap
	}

	js, err = json.Marshal(j)
	return
}
