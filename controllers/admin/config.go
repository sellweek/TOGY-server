package admin

import (
	"appengine/datastore"
	"models/action"
	"models/configuration/config"
	"models/configuration/timeConfig"
	"net/http"
	"strconv"
	"time"
	"util"
)

const (
	timeFormat = "15:04"
	dateFormat = "2.1.2006"
)

//ShowConfig handles showing the page where user can see and edit
//the central configuration for clients.
func ShowConfig(c util.Context) (err error) {
	conf, err := config.Get(c)
	if err != nil {
		if err == datastore.Done {
			conf, err = mkConfig(c)
			if err != nil {
				c.Errorf("Error when creating new config: %v", err)
			}
			http.Redirect(c.W, c.R, "/admin/config", 303)
			err = nil
			return
		} else {
			return
		}
	}

	as, err := action.GetFor(conf, c)
	if err != nil {
		return
	}

	tcs, err := timeConfig.GetAll(c)
	if err != nil {
		return
	}

	a := prepareActions(as)

	util.RenderLayout("config.html", "Všeobecné nastavenia", struct {
		Conf     *config.Config
		A        map[string][]time.Time
		ZeroTime time.Time
		Tz       *time.Location
		Tcs      []*timeConfig.TimeConfig
	}{conf, a, time.Date(0001, 01, 01, 00, 00, 00, 00, utc), util.Tz, tcs}, c, "/static/js/config.js")
	return
}

//Creates a default config
func mkConfig(c util.Context) (conf *config.Config, err error) {
	conf = &config.Config{
		StandardOn:     time.Date(0, 0, 0, 6, 0, 0, 0, time.UTC),
		StandardOff:    time.Date(0, 0, 0, 18, 0, 0, 0, time.UTC),
		UpdateInterval: 300,
		OverrideState:  0,
		Weekends:       false,
	}
	err = conf.Save(c)
	if err != nil {
		return
	}
	return
}

//SetConfig handles saving the new configuration to Datastore.
func SetConfig(c util.Context) (err error) {
	conf, err := config.Get(c)
	if err != nil {
		return
	}

	on, err := time.Parse(timeFormat, c.R.FormValue("standardOn"))
	if err != nil {
		return
	}
	off, err := time.Parse(timeFormat, c.R.FormValue("standardOff"))
	if err != nil {
		return
	}
	conf.OverrideState, err = strconv.Atoi(c.R.FormValue("overrideState"))
	if err != nil {
		return
	}

	conf.UpdateInterval, err = strconv.Atoi(c.R.FormValue("updateInterval"))
	if err != nil {
		return
	}

	if c.R.FormValue("weekends") == "true" {
		conf.Weekends = true
	}

	conf.StandardOn = util.NormalizeTime(on)
	conf.StandardOff = util.NormalizeTime(off)

	err = conf.Save(c)
	if err != nil {
		return
	}
	http.Redirect(c.W, c.R, "/admin/config", 303)
	return
}

//TimeOverrideEdit handles editing of existing time overrides.
//If it doesn't find id value in the path, it adds a new override.
func TimeOverrideEdit(c util.Context) (err error) {
	var tc *timeConfig.TimeConfig
	if key := c.Vars["id"]; key == "" {
		tc = nil
	} else {
		var k *datastore.Key
		k, err = datastore.DecodeKey(key)
		if err != nil {
			return
		}

		tc, err = timeConfig.GetByKey(k, c)
		if err != nil {
			return
		}
	}
	c.Infof("%+v", tc)
	util.RenderLayout("timeConfigEdit.html", "Úprava výnimky", tc, c, "/static/js/editTC.js")
	return
}

//TimeOverrideSubmit handles saving of time overrides into Datastore.
func TimeOverrideSubmit(c util.Context) (err error) {
	date, err := time.Parse(dateFormat, c.R.FormValue("date"))
	if err != nil {
		return
	}

	on, err := time.Parse(timeFormat, c.R.FormValue("on"))
	if err != nil {
		return
	}

	off, err := time.Parse(timeFormat, c.R.FormValue("off"))
	if err != nil {
		return
	}
	tc := timeConfig.New(util.NormalizeDate(date), util.NormalizeTime(on), util.NormalizeTime(off))
	var k *datastore.Key
	if _, ok := c.Vars["id"]; ok {
		k, err = datastore.DecodeKey(c.Vars["id"])
		if err != nil {
			return
		}
	} else {
		k = nil
	}
	tc.SetKey(k)
	err = tc.Save(c)
	if err != nil {
		return
	}
	http.Redirect(c.W, c.R, "/admin/config", 303)
	return
}

//TimeOverrideDelete handles deleting of a time override.
func TimeOverrideDelete(c util.Context) (err error) {
	key, err := datastore.DecodeKey(c.R.FormValue("key"))
	if err != nil {
		return
	}

	tc, err := timeConfig.GetByKey(key, c)
	if err != nil {
		return
	}
	err = tc.Delete(c)
	if err != nil {
		return
	}
	http.Redirect(c.W, c.R, "/admin/config", 303)
	return
}
