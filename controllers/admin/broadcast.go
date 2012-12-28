package admin

import (
	"appengine/blobstore"
	"github.com/russross/blackfriday"
	"html/template"
	"models/action"
	"models/presentation"
	"net/http"
	"net/url"
	"strings"
	"time"
	"util"
)

//UTC time zone
var utc, _ = time.LoadLocation("UTC")

//Redirects to presentation upload page
func Admin(c util.Context) {
	http.Redirect(c.W, c.R, "/admin/presentation/upload", 301)
}

//Handles the new presentation upload page.
func Upload(c util.Context) {
	p, err := presentation.GetActive(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	uploadURL, err := blobstore.UploadURL(c.Ac, "/admin/presentation/upload", nil)
	if err != nil {
		util.Log500(err, c)
		return
	}
	util.RenderLayout("upload.html", "Nahrať prezentáciu", struct {
		ActivePresentation string
		UploadURL          *url.URL
	}{p.Name, uploadURL}, c)
}

//Handles upload of a new presentation and saving its metadata
//to Datastore.
//
//Doesn't support filenames with non-ASCII characters. GAE encodes
//those into base-64 string with encoding prefixed and I don't want
//to include additional logic to differentiate between ASCII and
//non-ASCII filenames.
func UploadHandler(c util.Context) {
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

	p, err := presentation.Make(blob.BlobKey, fileType, name, []byte(formVal["description"][0]), active, c.Ac)
	if err != nil {
		util.Log500(err, c)
	}

	http.Redirect(c.W, c.R, "/admin/presentation/"+p.Key, 303)
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

	desc := blackfriday.MarkdownCommon(p.Description)

	secs := make([]float64, 0)
	for _, t := range a {
		dur := t[2].Sub(t[1])
		secs = append(secs, dur.Seconds())
	}

	avgDL := util.Round(util.Average(secs...), 2)

	util.RenderLayout("presentation.html", "Info o prezentácií", struct {
		P        *presentation.Presentation
		A        map[string][]time.Time
		Desc     template.HTML
		ZeroTime time.Time
		Avg      float64
	}{p, a, template.HTML(desc), time.Date(0001, 01, 01, 00, 00, 00, 00, utc), avgDL}, c, "/static/js/jquery-1.8.3.js", "/static/js/underscore-min.js", "/static/js/presentation.js")
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
	http.Redirect(c.W, c.R, "/admin/presentation/archive", 303)
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
	http.Redirect(c.W, c.R, "/admin/presentation/archive", 303)
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
