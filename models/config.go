package models

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"time"
)

const (
	ConfTimeFormat = "15:04"
	ConfDateFormat = "2006-1-2"
	OverrideOn     = 1
	OverrideOff    = -1
	NoOverride     = 0
)

var tz, _ = time.LoadLocation("UTC")

//Config stores central configuration in Datastore.
//There is always only one Config record in Datastore.
type Config struct {
	StandardOn     time.Time
	StandardOff    time.Time
	OverrideState  int
	UpdateInterval int
}

//GetKey returns key of the Datastore Config record.
func (c Config) GetKey(ctx appengine.Context) (k *datastore.Key, err error) {
	k, err = datastore.NewQuery("Config").KeysOnly().Run(ctx).Next(nil)
	return
}

//SaveConfig saves data provided into the Config record.
func (c *Config) Save(ctx appengine.Context) (err error) {
	var key *datastore.Key
	key, err = c.GetKey(ctx)
	if err == datastore.Done {
		ctx.Infof("Creating new config key")
		key = datastore.NewIncompleteKey(ctx, "Config", nil)
	} else if err != nil {
		return
	}

	c.StandardOn = normalizeTime(c.StandardOn)
	c.StandardOff = normalizeTime(c.StandardOff)

	_, err = datastore.Put(ctx, key, c)
	if err != nil {
		return fmt.Errorf("Error when putting: %v", err)
	}

	DeleteQueryTimesFor(&Config{}, ctx)
	return
}

//GetConfig fetches the Config record from Datastore and
//returns its data.
func GetConfig(ctx appengine.Context) (c Config, err error) {
	_, err = datastore.NewQuery("Config").Run(ctx).Next(&c)
	if err != nil {
		return
	}
	return
}

type TimeConfig struct {
	Date time.Time
	On   time.Time
	Off  time.Time
	Key  string `datastore:"-"`
}

func NewTimeConfig(date, on, off time.Time) (tc *TimeConfig) {
	tc = new(TimeConfig)
	tc.Date = normalizeDate(date)
	tc.Off = normalizeTime(off)
	tc.On = normalizeTime(on)
	return
}

func MakeTimeConfig(date, on, off time.Time, c appengine.Context) (tc *TimeConfig, err error) {
	tc = NewTimeConfig(date, on, off)
	err = tc.Save(c)
	return
}

func (tc *TimeConfig) Save(c appengine.Context) (err error) {
	tc.On = normalizeTime(tc.On)
	tc.Off = normalizeTime(tc.Off)
	tc.Date = normalizeDate(tc.Date)

	var k *datastore.Key
	if tc.Key == "" {
		var confKey *datastore.Key
		confKey, err = new(Config).GetKey(c)
		if err != nil {
			return
		}

		k = datastore.NewIncompleteKey(c, "TimeConfig", confKey)
		tc.Key = k.Encode()
	} else {
		k, err = datastore.DecodeKey(tc.Key)
		if err != nil {
			return
		}
	}
	_, err = datastore.Put(c, k, tc)
	DeleteQueryTimesFor(&Config{}, c)
	return
}

func (tc *TimeConfig) Delete(c appengine.Context) (err error) {
	k, err := datastore.DecodeKey(tc.Key)
	if err != nil {
		return
	}

	err = datastore.Delete(c, k)
	if err != nil {
		return
	}
	tc.Key = ""

	return
}

func GetTimeConfigs(c appengine.Context) (tcs []*TimeConfig, err error) {
	keys, err := datastore.NewQuery("TimeConfig").GetAll(c, &tcs)
	if err != nil {
		return
	}
	for i, k := range keys {
		tcs[i].Key = k.Encode()
	}
	return
}

//GetTimeConfig gets the TimeConfig with given key from Datastore and returns a pointer to it.
func GetTimeConfig(key string, c appengine.Context) (tc *TimeConfig, err error) {
	k, err := datastore.DecodeKey(key)
	if err != nil {
		return
	}
	tc = new(TimeConfig)
	err = datastore.Get(c, k, tc)
	if err != nil {
		return
	}
	tc.Key = key
	return
}

func ConfigJSON(c appengine.Context) (js []byte, err error) {
	j := make(map[string]interface{})
	conf, err := GetConfig(c)
	if err != nil {
		return
	}

	tcs, err := GetTimeConfigs(c)
	if err != nil {
		return
	}

	j["StandardTimeSettings"] = map[string]string{
		"TurnOn":  conf.StandardOn.Format(ConfTimeFormat),
		"TurnOff": conf.StandardOff.Format(ConfTimeFormat),
	}
	j["UpdateInterval"] = conf.UpdateInterval
	switch conf.OverrideState {
	case NoOverride:
		j["OverrideOn"] = false
		j["OverrideOff"] = false
	case OverrideOn:
		j["OverrideOn"] = true
		j["OverrideOff"] = false
	case OverrideOff:
		j["OverrideOn"] = false
		j["OverrideOff"] = true
	}
	for _, tc := range tcs {
		timeMap := make(map[string]string)
		timeMap["TurnOn"] = tc.On.Format(ConfTimeFormat)
		timeMap["TurnOff"] = tc.Off.Format(ConfTimeFormat)

		j["OverrideDays"] = make(map[string]map[string]string)
		j["OverrideDays"].(map[string]map[string]string)[tc.Date.Format(ConfDateFormat)] = timeMap
	}

	js, err = json.Marshal(j)
	return
}

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, tz)
}

func normalizeTime(t time.Time) time.Time {
	return time.Date(1, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), tz)
}
