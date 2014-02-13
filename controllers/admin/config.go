package admin

import (
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
	conf, err := config.Get(c.Ac)
	if err != nil {
		return
	}

	as, err := action.GetFor(&config.Config{}, c.Ac)
	if err != nil {
		return
	}

	tcs, err := timeConfig.GetAll(c.Ac)
	if err != nil {
		return
	}

	a := prepareActions(as)

	util.RenderLayout("config.html", "Všeobecné nastavenia", struct {
		Conf     config.Config
		A        map[string][]time.Time
		ZeroTime time.Time
		Tz       *time.Location
		Tcs      []*timeConfig.TimeConfig
	}{conf, a, time.Date(0001, 01, 01, 00, 00, 00, 00, utc), util.Tz, tcs}, c, "/static/js/config.js")
	return
}

//SetConfig handles saving the new configuration to Datastore.
func SetConfig(c util.Context) (err error) {
	conf := new(config.Config)
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

	err = conf.Save(c.Ac)
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
		tc, err = timeConfig.GetByKey(key, c.Ac)
		if err != nil {
			return
		}
	}
	c.Ac.Infof("%+v", tc)
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
	tc.Key = c.Vars["id"]
	err = tc.Save(c.Ac)
	if err != nil {
		return
	}
	http.Redirect(c.W, c.R, "/admin/config", 303)
	return
}

//TimeOverrideDelete handles deleting of a time override.
func TimeOverrideDelete(c util.Context) (err error) {
	key := c.R.FormValue("key")
	tc, err := timeConfig.GetByKey(key, c.Ac)
	if err != nil {
		return
	}
	err = tc.Delete(c.Ac)
	if err != nil {
		return
	}
	http.Redirect(c.W, c.R, "/admin/config", 303)
	return
}
