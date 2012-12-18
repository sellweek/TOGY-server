package timeConfig

import (
	"appengine"
	"appengine/datastore"
	"models/action"
	"models/configuration/config"
	"time"
	"util"
)

//TimeConfig is a model used to store
//time override settings.
type TimeConfig struct {
	Date time.Time
	On   time.Time
	Off  time.Time
	Key  string `datastore:"-"`
}

//New returns a pointer to a TimeConfig with its fields set according to arguments.
func New(date, on, off time.Time) (tc *TimeConfig) {
	tc = new(TimeConfig)
	tc.Date = util.NormalizeDate(date)
	tc.Off = util.NormalizeTime(off)
	tc.On = util.NormalizeTime(on)
	return
}

//Make creates a new TimeConfig using New and saves it.
func Make(date, on, off time.Time, c appengine.Context) (tc *TimeConfig, err error) {
	tc = New(date, on, off)
	err = tc.Save(c)
	return
}

//Save saves a TimeConfig to Datastore.
//If its Key field is set, it will replace existing record
//that has that key. If not, it will use datastore.NewIncompleteKey()
//to create a new key and set the field.
func (tc *TimeConfig) Save(c appengine.Context) (err error) {
	tc.On = util.NormalizeTime(tc.On)
	tc.Off = util.NormalizeTime(tc.Off)
	tc.Date = util.NormalizeDate(tc.Date)

	var k *datastore.Key
	if tc.Key == "" {
		var confKey *datastore.Key
		confKey, err = new(config.Config).GetKey(c)
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
	action.DeleteFor(&config.Config{}, c)
	return
}

//Delete deletes a TimeConfig from Datastore, emptying its Key field.
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

//GetAll gets all TimeConfigs saved in Datastore and returns them in a slice.
func GetAll(c appengine.Context) (tcs []*TimeConfig, err error) {
	keys, err := datastore.NewQuery("TimeConfig").GetAll(c, &tcs)
	if err != nil {
		return
	}
	for i, k := range keys {
		tcs[i].Key = k.Encode()
	}
	return
}

//GetByKey gets the TimeConfig with given key from Datastore and returns a pointer to it.
func GetByKey(key string, c appengine.Context) (tc *TimeConfig, err error) {
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
