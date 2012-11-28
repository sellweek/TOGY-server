package models

import (
	"appengine"
	"appengine/datastore"
	"time"
)

const confTimeFormat = "15:04"
const confDateFormat = "2006-1-2"

//Config stores central configuration in Datastore.
//There is always only one Config record in Datastore.
type Config struct {
	StandardOn     time.Time
	StandardOff    time.Time
	OverrideOn     bool
	OverrideOff    bool
	UpdateInterval int
}

//GetKey returns key of the Datastore Config record.
func (c Config) GetKey(ctx appengine.Context) (k *datastore.Key, err error) {
	k, err = datastore.NewQuery("Config").KeysOnly().Run(ctx).Next(nil)
	return
}

//SaveConfig saves data provided into the Config record.
func (c Config) SaveConfig(ctx appengine.Context) (err error) {
	var key *datastore.Key
	key, err = c.GetKey(ctx)
	if err == datastore.Done {
		key = datastore.NewIncompleteKey(ctx, "Config", nil)
	} else if err != nil {
		return
	}

	c.StandardOn = normalizeTime(standardOn)
	c.StandardOff = normalizeTime(standardOff)

	_, err = datastore.Put(ctx, key, c)
	if err != nil {
		return
	}

	DeleteQueryTimesFor(Config{}, ctx)
	return
}

//GetConfig fetches the Config record from Datastore and
//returns its data.
func GetConfig(ctx appengine.Context) (c Config, err error) {
	var conf Config
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
	Key  datastore.Key `datastore:"-"`
}

func NewTimeConfig(date, on, off time.Time) (tc TimeConfig) {
	tc = new(TimeConfig)
	tc.Date = date
	tc.Off = off
	tc.On = on
	return
}

func MakeTimeConfig(date, on, off time.Time, c appengine.Context) (tc TimeConfig, err error) {
	tc = NewTimeConfig(date, on, off)
	err = tc.Save()
	return
}

func (tc *TimeConfig) Save(c appengine.Context) (err error) {
	var k *datastore.Key
	if tc.Key == nil {
		confKey, err := Config{}.GetKey(c)
		if err != nil {
			return
		}

		k, err = datastore.NewIncompleteKey(c, "TimeConfig", confKey)
		if err != nil {
			return
		}
		tc.Key = k.Encode()
	} else {
		k, err = datastore.DecodeKey(tc.Key)
		if err != nil {
			return
		}
	}
	_, err = datastore.Put(c, k, tc)
	return
}

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, Tz)
}

func normalizeTime(t time.Time) time.Time {
	return time.Date(0, 0, 0, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), Tz)
}
