package config

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"models/action"
	"reflect"
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
type Config struct {
	StandardOn     time.Time
	StandardOff    time.Time
	OverrideState  int
	UpdateInterval int
	Weekends       bool
	Timestamp      int64
	key            *datastore.Key `datastore:"-"`
}

var typ = reflect.TypeOf(Config{})

func (c *Config) Key() *datastore.Key {
	return c.key
}

func (c *Config) SetKey(k *datastore.Key) {
	c.key = k
}

func (c *Config) Kind() string {
	return "Config"
}

func (_ *Config) Ancestor() *datastore.Key {
	return nil
}

//SaveConfig saves data provided into the Config record.
func (c *Config) Save(ctx appengine.Context) (err error) {
	if c.Key() == nil {
		ctx.Infof("Creating new config key")
		key := datastore.NewIncompleteKey(ctx, "Config", nil)
		c.SetKey(key)
	} else if err != nil {
		return
	}

	c.StandardOn = util.NormalizeTime(c.StandardOn)
	c.StandardOff = util.NormalizeTime(c.StandardOff)

	c.Timestamp = time.Now().Unix()

	c.forceUTC()
	_, err = datastore.Put(ctx, c.Key(), c)
	c.forceLocal()
	if err != nil {
		return fmt.Errorf("Error when putting: %v", err)
	}

	action.DeleteFor(c, ctx)
	return
}

//Get fetches the Config record from Datastore and
//returns its data.
func Get(ctx appengine.Context) (c *Config, err error) {
	var value_c Config
	key, err := datastore.NewQuery("Config").Run(ctx).Next(&value_c)
	if err != nil {
		return
	}
	c = &value_c
	c.forceLocal()
	c.SetKey(key)
	return
}

func UpdateTimestamp(ctx appengine.Context) (err error) {
	c, err := Get(ctx)
	if err != nil {
		return
	}

	c.Timestamp = time.Now().Unix()
	//This also removes Actions for us
	err = c.Save(ctx)
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
