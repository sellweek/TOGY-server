package public

import (
	"github.com/russross/blackfriday"
	"html/template"
	"models/presentation"
	"net/http"
	"strconv"
	"util"
)

const perPage = 5

//Index redirects user to the first page of presentation listing.
func Index(c util.Context) {
	http.Redirect(c.W, c.R, "/presentation?p=1", 301)
}

//Presentations shows listing of presentations in paginated form.
func Presentations(c util.Context) {
	page, err := strconv.Atoi(c.R.FormValue("p"))
	if err != nil {
		util.Log500(err, c)
		return
	}
	ps, err := presentation.GetListing(page, perPage, c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}

	type tplData struct {
		P presentation.Presentation
		D template.HTML
	}

	data := make([]tplData, len(ps), len(ps))

	for _, p := range ps {
		data = append(data, tplData{P: *p, D: template.HTML(blackfriday.MarkdownCommon(p.Description))})
	}

	maxPages, err := presentation.PageCount(c.Ac, perPage)
	if err != nil {
		util.Log500(err, c)
		return
	}

	util.RenderLayout("index.html", "Zoznam vysielaní", struct {
		Page     int
		MaxPages int
		Data     []tplData
	}{page, maxPages, data}, c)
}
