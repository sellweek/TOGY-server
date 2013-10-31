package admin

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"github.com/russross/blackfriday"
	"html/template"
	"models/action"
	"models/activation"
	"models/presentation"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"util"
)

//UTC time zone
//Used in ZeroTime comparisions.
var utc, _ = time.LoadLocation("UTC")

//Admin redirects to presentation upload page
func Admin(c util.Context) (err error) {
	http.Redirect(c.W, c.R, "/admin/presentation/upload", 301)
	return
}

//Upload renders the new presentation upload page.
func Upload(c util.Context) (err error) {
	uploadURL, err := blobstore.UploadURL(c.Ac, "/admin/presentation/upload", nil)
	if err != nil {
		return
	}

	acts, err := activation.GetAfterTime(time.Now(), c.Ac)
	if err != nil {
		return
	}

	type actWithName struct {
		A *activation.Activation
		P *presentation.Presentation
	}

	ans := make([]actWithName, len(acts))

	for i, a := range acts {
		pk := a.Presentation.Encode()
		var p *presentation.Presentation
		p, err = presentation.GetByKey(pk, c.Ac)
		if err != nil {
			c.Ac.Errorf("Could not load presentation: %v", err)
			continue
		}
		ans[i] = actWithName{a, p}
	}

	util.RenderLayout("upload.html", "Nahrať prezentáciu", struct {
		ActivePresentation string
		UploadURL          *url.URL
		Ans                []actWithName
	}{activeName, uploadURL, ans}, c)
	return
}

//UploadHandler handles upload of a new presentation and saving its metadata
//to Datastore.
//
//Doesn't support filenames with non-ASCII characters. GAE encodes
//those into base-64 string with encoding prefixed and I don't want
//to include additional logic to differentiate between ASCII and
//non-ASCII filenames.
func UploadHandler(c util.Context) (err error) {
	blobs, formVal, err := blobstore.ParseUpload(c.R)
	if err != nil {
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

	p, err := presentation.Make(blob.BlobKey, fileType, name, []byte(formVal["description"][0]), active, c.Ac)
	if err != nil {
		return
	}

	http.Redirect(c.W, c.R, "/admin/presentation/"+p.Key, 303)
	return
}

//Archive handles showing listing of presentations.
func Archive(c util.Context) (err error) {
	page, err := strconv.Atoi(c.Vars["page"])
	if err != nil {
		return
	}

	type tmplData struct {
		P *presentation.Presentation
		C int
	}
	ps, err := presentation.GetListing(page, 10, c.Ac)
	if err != nil {
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

	maxPages, err := presentation.PageCount(10, c.Ac)
	if err != nil {
		return
	}

	util.RenderLayout("archive.html", "Archív prezentácií", struct {
		Data     []tmplData
		Page     int
		MaxPages int
		Tz       *time.Location
	}{downloads, page, maxPages, util.Tz}, c)
	return
}

//Presentation handles showing page with details about a presentation.
func Presentation(c util.Context) (err error) {
	p, err := presentation.GetByKey(c.Vars["id"], c.Ac)
	if err != nil {
		return
	}
	as, err := action.GetFor(p, c.Ac)
	if err != nil {
		return
	}

	a := prepareActions(as)

	desc := blackfriday.MarkdownCommon(p.Description)

	secs := make([]float64, 0)
	for _, t := range a {
		dur := t[2].Sub(t[1])
		secs = append(secs, dur.Seconds())
	}

	avgDL := util.Round(util.Average(secs...), 2)

	//We can safely ignore errors here since we already
	//got the presentation using the same key
	pk, _ := datastore.DecodeKey(p.Key)

	acts, err := activation.GetForPresentation(pk, c.Ac)
	if err != nil {
		return
	}

	util.RenderLayout("presentation.html", "Info o prezentácií", struct {
		P           *presentation.Presentation
		A           map[string][]time.Time
		Desc        template.HTML
		ZeroTime    time.Time
		Avg         float64
		Domain      string
		Activations []*activation.Activation
		Tz          *time.Location
	}{p, a, template.HTML(desc), time.Date(0001, 01, 01, 00, 00, 00, 00, utc), avgDL, appengine.DefaultVersionHostname(c.Ac), acts, util.Tz}, c, "/static/js/underscore-min.js", "/static/js/jquery-ui-1.9.2.custom.min.js", "/static/js/timepicker-min.js", "/static/js/presentation.js")
	return
}

//Activate handles activation of presentation.
func Activate(c util.Context) (err error) {
	key := c.R.FormValue("id")
	p, err := presentation.GetByKey(key, c.Ac)
	if err != nil {
		return
	}
	p.Active = true
	p.Save(c.Ac)
	http.Redirect(c.W, c.R, "/admin/presentation/archive/1", 303)
	return
}

//Delete handles deleting of presentation.
func Delete(c util.Context) (err error) {
	key := c.R.FormValue("id")
	p, err := presentation.GetByKey(key, c.Ac)
	if err != nil {
		return
	}
	err = p.Delete(c.Ac)
	if err != nil {
		return
	}
	http.Redirect(c.W, c.R, c.Referer(), 303)
	return
}

func Deactivate(c util.Context) (err error) {
	key := c.R.FormValue("id")
	p, err := presentation.GetByKey(key, c.Ac)
	if err != nil {
		return
	}

	p.Active = false
	err = p.Save(c.Ac)
	if err != nil {
		return
	}

	http.Redirect(c.W, c.R, c.R.Referer(), 303)
	return
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
