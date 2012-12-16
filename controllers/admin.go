package controllers

import (
	"appengine/blobstore"
	"appengine/datastore"
	"fmt"
	"models/action"
	"models/configuration/config"
	"models/configuration/timeConfig"
	"models/presentation"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"util"
)

//UTC time zone
var utc, _ = time.LoadLocation("UTC")

//Handles the new presentation upload page.
func Admin(c util.Context) {
	p, err := presentation.GetActive(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	uploadURL, err := blobstore.UploadURL(c.Ac, "/admin/upload", nil)
	if err != nil {
		util.Log500(err, c)
		return
	}
	util.RenderLayout("admin.html", "Nahrať prezentáciu", struct {
		ActivePresentation string
		UploadURL          *url.URL
	}{p.Name, uploadURL}, c)
}

//Handles upload of a new presentation and saving its metadata
//to Datastore.
//
//Doesn't support filenames with non-ASCII characters. GAE endocdes
//those into base-64 string with encoding prefixed and I don't want
//to include additional logic to differentiate between ASCII and
//non-ASCII filenames.
func Upload(c util.Context) {
	blobs, formVal, err := blobstore.ParseUpload(c.R)
	if err != nil {
		util.Log500(err, c)
		return
	}
	blob := blobs["file"][0]
	fn := strings.Split(blob.Filename, ".")
	fileType := fn[len(fn)-1]

	var active bool
	if len(formVal["activate"]) == 0 {
		active = false
	} else {
		active = true
	}

	name := formVal["name"][0]
	if name == "" {
		name = "Neznáma prezentácia z " + time.Now().Format("2.1.2006")
	}

	_, err = presentation.Make(blob.BlobKey, fileType, name, active, c.Ac)
	if err != nil {
		util.Log500(err, c)
	}

	http.Redirect(c.W, c.R, "/admin", 301)
}

//Handles showing listing of presentations.
func Archive(c util.Context) {
	type tmplData struct {
		P *presentation.Presentation
		C int
	}
	ps, err := presentation.GetAll(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}

	downloads := make([]tmplData, 0)
	for _, p := range ps {
		count, err := action.GetCountFor(action.DownloadFinish, p, c.Ac)
		if err != nil {
			c.Ac.Infof("Error when getting download count: %v", err)
			count = -1
		}
		downloads = append(downloads, tmplData{p, count})
	}
	util.RenderLayout("archive.html", "Archív prezentácií", downloads, c)
}

//Handles showing page with details about a presentation.
func Presentation(c util.Context) {
	p, err := presentation.GetByKey(c.Vars["id"], c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	as, err := action.GetFor(p, c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}

	a := prepareActions(as)

	secs := make([]float64, 0)
	for _, t := range a {
		dur := t[2].Sub(t[1])
		secs = append(secs, dur.Seconds())
	}

	avgDL := util.Round(util.Average(secs...), 2)

	util.RenderLayout("presentation.html", "Info o prezentácií", struct {
		P        *presentation.Presentation
		A        map[string][]time.Time
		ZeroTime time.Time
		Avg      float64
	}{p, a, time.Date(0001, 01, 01, 00, 00, 00, 00, utc), avgDL}, c)
}

//Handles activation of presentation.
func Activate(c util.Context) {
	key := c.R.FormValue("id")
	p, err := presentation.GetByKey(key, c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	p.Active = true
	p.Save(c.Ac)
	http.Redirect(c.W, c.R, "/admin/archive", 301)
}

//Handles deleting of presentation.
func Delete(c util.Context) {
	key := c.R.FormValue("id")
	p, err := presentation.GetByKey(key, c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	err = p.Delete(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	http.Redirect(c.W, c.R, "/admin/archive", 301)
}

//Handles showing the page in which user can see and edit
//the central configuration for clients.
func ShowConfig(c util.Context) {
	conf, err := config.Get(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}

	as, err := action.GetFor(&config.Config{}, c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}

	a := prepareActions(as)

	util.RenderLayout("config.html", "Konfigurácia obrazoviek", struct {
		Conf     config.Config
		A        map[string][]time.Time
		ZeroTime time.Time
	}{conf, a, time.Date(0001, 01, 01, 00, 00, 00, 00, utc)}, c, "/static/js/jquery-1.8.3.js", "/static/js/jquery-ui-1.9.2.custom.min.js", "/static/js/timepicker.js", "/static/js/config.js")
}

//Handles saving the new configuration to Datastore.
func SetConfig(c util.Context) {
	var err error
	conf := new(config.Config)
	conf.StandardOn, err = time.Parse(config.ConfTimeFormat, c.R.FormValue("standardOn"))
	if err != nil {
		util.Log500(err, c)
		return
	}
	conf.StandardOff, err = time.Parse(config.ConfTimeFormat, c.R.FormValue("standardOff"))
	if err != nil {
		util.Log500(err, c)
		return
	}
	conf.OverrideState, err = strconv.Atoi(c.R.FormValue("overrideState"))
	if err != nil {
		util.Log500(err, c)
		return
	}

	conf.UpdateInterval, err = strconv.Atoi(c.R.FormValue("updateInterval"))
	if err != nil {
		util.Log500(err, c)
		return
	}
	err = conf.Save(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	http.Redirect(c.W, c.R, "/admin/config", 301)
}

func TimeOverride(c util.Context) {
	tcs, err := timeConfig.GetAll(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	util.RenderLayout("timeConfig.html", "Zoznam časových výnimiek", tcs, c)
}

func TimeOverrideEdit(c util.Context) {
	var tc *timeConfig.TimeConfig
	var err error
	if key := c.Vars["id"]; key == "" {
		tc = nil
	} else {
		tc, err = timeConfig.GetByKey(key, c.Ac)
		if err != nil {
			util.Log500(err, c)
			return
		}
	}
	util.RenderLayout("timeConfigEdit.html", "Úprava výnimky", tc, c, "/static/js/jquery-1.8.3.js", "/static/js/jquery-ui-1.9.2.custom.min.js", "/static/js/timepicker.js", "/static/js/editTC.js")

}

func TimeOverrideSubmit(c util.Context) {
	date, err := time.Parse(config.ConfDateFormat, c.R.FormValue("date"))
	if err != nil {
		util.Log500(err, c)
		return
	}

	on, err := time.Parse(config.ConfTimeFormat, c.R.FormValue("on"))
	if err != nil {
		util.Log500(err, c)
		return
	}

	off, err := time.Parse(config.ConfTimeFormat, c.R.FormValue("off"))
	if err != nil {
		util.Log500(err, c)
		return
	}
	tc := timeConfig.New(date, on, off)
	tc.Key = c.Vars["id"]
	err = tc.Save(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	http.Redirect(c.W, c.R, "/admin/config/timeOverride", 301)
}

//Handles deleting of a time override.
func TimeOverrideDelete(c util.Context) {
	key := c.R.FormValue("key")
	tc, err := timeConfig.GetByKey(key, c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	err = tc.Delete(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	http.Redirect(c.W, c.R, "/admin/config/timeOverride", 301)
}

//Inserts fake presentation and config into datastore.
//Used when the system doesn't have any presentation inserted
//and is in an inconsistent state because of that.
func Bootstrap(c util.Context) {
	p := presentation.New("test", "xxx", "DO NOT USE! DO NOT ACTIVATE IF YOU DON'T KNOW WHAT YOU'RE DOING", true)
	_, err := datastore.Put(c.Ac, datastore.NewIncompleteKey(c.Ac, "Presentation", nil), p)
	if err != nil {
		fmt.Fprintln(c.W, "Error with presentation: ", err)
	}

	//	zeroTime := time.Date(0001, 01, 01, 00, 00, 00, 00, utc)

	conf := new(config.Config)

	err = conf.Save(c.Ac)

	if err != nil {
		fmt.Fprintln(c.W, "Error with config:", err)
	}
	fmt.Fprint(c.W, "Do not start any clients until you have replaced this presentation.")
}

func Migrate(c util.Context) {
	type queryTime struct {
		Action       action.ActionType
		Presentation *datastore.Key
		Time         time.Time
		Client       string
	}
	ps, err := presentation.GetAll(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	for _, p := range ps {
		pKey, _ := datastore.DecodeKey(p.Key)
		qts := make([]queryTime, 15)
		_, err := datastore.NewQuery("QueryTime").Ancestor(pKey).GetAll(c.Ac, &qts)
		if err != nil {
			util.Log500(err, c)
			return
		}

		for _, q := range qts {
			a := action.New(q.Presentation, q.Action, q.Client)
			a.Time = q.Time
			err = a.Save(c.Ac)
			if err != nil {
				util.Log500(fmt.Errorf("Couldn't put: %v", err), c)
				return
			}
		}
	}

}

func prepareActions(as []action.Action) map[string][]time.Time {
	a := make(map[string][]time.Time)

	for _, v := range as {
		if a[v.Client] == nil {
			a[v.Client] = make([]time.Time, 3, 3)
		}
		a[v.Client][int(v.Type)] = v.Time
	}
	return a
}
