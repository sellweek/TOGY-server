package config

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"models/action"
	"time"
	"util"
)

const (
	OverrideOn  = 1
	OverrideOff = -1
	NoOverride  = 0
)

//Config stores central configuration in Datastore.
//There is always only one Config record in Datastore.
//Because of that, Key field is not needed.
type Config struct {
	StandardOn     time.Time
	StandardOff    time.Time
	OverrideState  int
	UpdateInterval int
	Weekends       bool
	Timestamp      int64
}

//GetKey returns key of the Datastore Config record.
//It will always return the key of the single Datastore record,
//even if its called like
//	Config{}.GetKey
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

	c.StandardOn = util.NormalizeTime(c.StandardOn)
	c.StandardOff = util.NormalizeTime(c.StandardOff)

	c.Timestamp = time.Now().Unix()

	c.forceUTC()
	_, err = datastore.Put(ctx, key, c)
	c.forceLocal()
	if err != nil {
		return fmt.Errorf("Error when putting: %v", err)
	}

	action.DeleteFor(&Config{}, ctx)
	return
}

//Get fetches the Config record from Datastore and
//returns its data.
func Get(ctx appengine.Context) (c Config, err error) {
	_, err = datastore.NewQuery("Config").Run(ctx).Next(&c)
	if err != nil {
		return
	}
	c.forceLocal()
	return
}

func (c *Config) forceUTC() {
	c.force(time.UTC)
}

func (c *Config) forceLocal() {
	c.force(util.C.Tz)
}

func (c *Config) force(loc *time.Location) {
	c.StandardOff = time.Date(1, 1, 1, c.StandardOff.Hour(), c.StandardOff.Minute(), c.StandardOff.Second(), c.StandardOff.Nanosecond(), loc)
	c.StandardOn = time.Date(1, 1, 1, c.StandardOn.Hour(), c.StandardOn.Minute(), c.StandardOn.Second(), c.StandardOn.Nanosecond(), loc)
}
