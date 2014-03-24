package public

import (
	"appengine"
	"github.com/russross/blackfriday"
	"html/template"
	"models/presentation"
	"net/http"
	"strconv"
	"util"
)

const perPage = 5

//Index redirects user to the first page of presentation listing.
func Index(c util.Context) (err error) {
	http.Redirect(c.W, c.R, "/presentation?p=1", 301)
	return
}

//Presentations shows listing of presentations in paginated form.
func Presentations(c util.Context) (err error) {
	page, err := strconv.Atoi(c.R.FormValue("p"))
	if err != nil {
		return
	}
	ps, err := presentation.GetListing(page, perPage, c)
	if err != nil {
		return
	}

	type templateData struct {
		P presentation.Presentation
		D template.HTML
	}

	data := make([]templateData, len(ps), len(ps))

	for _, p := range ps {
		data = append(data, templateData{P: *p, D: template.HTML(blackfriday.MarkdownCommon(p.Description))})
	}

	maxPages, err := presentation.PageCount(perPage, c)
	if err != nil {
		return
	}

	c.Infof("Hostname: %v", appengine.DefaultVersionHostname(c))

	util.RenderLayout("index.html", "Zoznam vysielaní", struct {
		Page     int
		MaxPages int
		Data     []templateData
		Domain   string
	}{Page: page, MaxPages: maxPages, Data: data, Domain: appengine.DefaultVersionHostname(c)}, c, "/static/js/index.js")
	return
}
