package timeConfig

import (
	"appengine"
	"appengine/datastore"
	"github.com/sellweek/gaemodel"
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
	key  *datastore.Key `datastore:"-"`
}

var typ = reflect.TypeOf(TimeConfig{})

func (tc *TimeConfig) Key() *datastore.Key {
	return tc.key
}

func (tc *TimeConfig) SetKey(k *datastore.Key) {
	tc.key = k
}

func (tc *TimeConfig) Kind() string {
	return "Config"
}

func (_ *TimeConfig) Ancestor() *datastore.Key {
	return nil
}

//New returns a pointer to a TimeConfig with its fields set according to arguments.
func New(date, on, off time.Time) (tc *TimeConfig) {
	tc = new(TimeConfig)
	tc.Date = date
	tc.Off = off
	tc.On = on
	return
}

//Make creates a new TimeConfig using New and saves it.
func Make(date, on, off time.Time, c appengine.Context) (tc *TimeConfig, err error) {
	tc = New(date, on, off)
	err = tc.Save(c)
	return
}

//Save saves a TimeConfig to Datastore.
//If its Key field is set, it will replace an existing record
//that has that key. If not, it will use datastore.NewIncompleteKey()
//to create a new key and set the field.
func (tc *TimeConfig) Save(c appengine.Context) (err error) {
	tc.forceUTC()
	defer tc.forceLocal()

	err = gaemodel.Save(c, tc)
	if err != nil {
		return
	}

	err = config.UpdateTimestamp(c)
	return
}

//Delete deletes a TimeConfig from Datastore, emptying its Key field.
func (tc *TimeConfig) Delete(c appengine.Context) (err error) {
	err = gaemodel.Delete(c, tc)
	if err != nil {
		return
	}

	err = config.UpdateTimestamp(c)

	return
}

//GetAll gets all TimeConfigs saved in Datastore and returns them in a slice.
func GetAll(c appengine.Context) (tcs []*TimeConfig, err error) {
	is, err = gaemodel.GetAll(c, typ, "TimeConfig", 0, 0)
	if err != nil {
		return
	}
	tcs = is.([]*TimeConfig)
	for _, tc := range tcs {
		tc.forceLocal()
	}
	return
}

//GetByKey gets the TimeConfig with given key from Datastore and returns a pointer to it.
func GetByKey(key *datastore.Key, c appengine.Context) (tc *TimeConfig, err error) {
	tc = new(TimeConfig)
	err = datastore.Get(c, k, tc)
	if err != nil {
		return
	}
	tc.Key = key
	tc.forceLocal()
	return
}

func (tc *TimeConfig) forceUTC() {
	tc.force(time.UTC)
}

func (tc *TimeConfig) forceLocal() {
	tc.force(util.C.Tz)
}

func (tc *TimeConfig) force(loc *time.Location) {
	tzDate := time.Date(tc.Date.Year(), tc.Date.Month(), tc.Date.Day(), 0, 0, 0, 0, loc)
	tzOff := time.Date(1, 1, 1, tc.Off.Hour(), tc.Off.Minute(), tc.Off.Second(), tc.Off.Nanosecond(), loc)
	tzOn := time.Date(1, 1, 1, tc.On.Hour(), tc.On.Minute(), tc.On.Second(), tc.On.Nanosecond(), loc)

	tc.Date = tzDate
	tc.On = tzOn
	tc.Off = tzOff
}
