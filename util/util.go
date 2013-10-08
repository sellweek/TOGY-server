//Package util provides utility functions.
package util

import (
	"appengine"
	"appengine/user"
	"github.com/gorilla/mux"
	"github.com/mjibson/appstats"
	"html/template"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	t = "templates" //Directory with templates
)

var templates *template.Template
var Tz, _ = time.LoadLocation("Europe/Bratislava")

//init injects few utility functions into templates we're using
func init() {
	templates = template.New("").Funcs(template.FuncMap{
		"equal": func(x, y int) bool {
			return x == y
		},
		"subtract": func(x, y int) int {
			return x - y
		},
		"add": func(x, y int) int {
			return x + y
		}})
	// List of template files. When creating new template, add it here.
	templates = template.Must(parseFiles(templates, t))
	templates = template.Must(parseFiles(templates, t+string(os.PathSeparator)+"layout"))
}

//ParseFiles goes trough a folder (non-recursively), parsing and
//adding all HTML files into a template.
func parseFiles(t *template.Template, dir string) (temp *template.Template, err error) {
	f, err := os.Open(dir)
	if err != nil {
		return
	}

	fis, err := f.Readdir(0)
	if err != nil {
		return
	}

	filenames := make([]string, 0)

	for _, fi := range fis {
		if fi.IsDir() || getFileType(fi.Name()) != "html" {
			continue
		}

		filenames = append(filenames, dir+string(os.PathSeparator)+fi.Name())
	}

	temp, err = t.ParseFiles(filenames...)
	return
}

//Context is the type used for passing data to handlers
type Context struct {
	Ac   appengine.Context
	W    http.ResponseWriter
	R    *http.Request
	Vars map[string]string
}

//Handler maps standard net/http handlers to handlers accepting Context
func Handler(hand func(Context) error) http.Handler {
	return appstats.NewHandler(func(c appengine.Context, w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		err := hand(Context{Ac: c, W: w, R: r, Vars: vars})
		if err != nil {
			c.Errorf("Error 500. %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

//RenderLayout inserts template with given name into the layout and sets the title and pipeline.
//The template should be loaded inside templates variable
//If any arguments are provided after the context, they will be treated as links
//to JavaScript scripts to load in the header of the template.
func RenderLayout(tmpl string, title string, data interface{}, c Context, jsIncludes ...string) {
	RenderTemplate("header.html", struct {
		Title      string
		JsIncludes []string
		Admin      bool
		AppName    string
	}{title, jsIncludes, user.IsAdmin(c.Ac), C.Title}, c)
	RenderTemplate(tmpl, data, c)
	RenderTemplate("footer.html", template.HTML(C.Footer), c)
}

//renderTemplate renders a single template
func RenderTemplate(tmpl string, data interface{}, c Context) {
	if err := templates.ExecuteTemplate(c.W, tmpl, data); err != nil {
		c.Ac.Errorf("Couldn't render template. %v", err)
		http.Error(c.W, err.Error(), http.StatusInternalServerError)
	}
}

//Average returns an Average of its arguments
func Average(nums ...float64) float64 {
	var total float64
	for _, x := range nums {
		total += x
	}
	return total / float64(len(nums))
}

// Round returns a rounded version of x with prec precision.
//
// Special cases are:
//	Round(±0) = ±0
//	Round(±Inf) = ±Inf
//	Round(NaN) = NaN
func Round(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow

	if intermed < 0.0 {
		intermed -= 0.5
	} else {
		intermed += 0.5
	}
	rounder = float64(int64(intermed))

	return rounder / float64(pow)
}

//NormalizeDate strips the time part from time.Date leaving only
//year, month and day.
//If forceTZ is true, its location will be set to util.Tz,
//if false, it will be left as is.
func NormalizeDate(t time.Time, forceTZ bool) time.Time {
	var tz *time.Location
	if forceTZ {
		tz = Tz
	} else {
		tz = t.Location()
	}

	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, tz)
}

//NormalizeTime strips the date part from time.Date leaving only
//hours, minutes, seconds and nanoseconds.
//If forceTZ is true, its location will be set to util.Tz,
//if false, it will be left as is.
func NormalizeTime(t time.Time, forceTZ bool) time.Time {
	var tz *time.Location
	if forceTZ {
		tz = Tz
	} else {
		tz = t.Location()
	}

	return time.Date(1, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), tz)
}

//getFileType returns file extension.
//In case of filenames with multiple extensions, only the last one is returned.
//For example:
//	getFileType("data.tar.gz")
//returns "gz"
func getFileType(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return ""
	}

	return parts[len(parts)-1]
}
