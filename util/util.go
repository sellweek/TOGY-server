//Package util provides utility functions used in my this web project.
package util

import (
	"appengine"
	"appengine/user"
	"github.com/gorilla/mux"
	"html/template"
	"math"
	"net/http"
	"time"
)

const (
	TimeFormat = "20060102150405"                                                   //Time format used in update queries and reponses and in filenames.
	t          = "templates/"                                                       //Directory with templates
	JqCDN      = "https://ajax.googleapis.com/ajax/libs/jquery/1.9.0/jquery.min.js" //Address of jQuery
)

var templates *template.Template
var Tz, _ = time.LoadLocation("UTC")

//In init we inject few utility functions into templates we're using
func init() {
	temp := template.New("").Funcs(template.FuncMap{
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
	templates = template.Must(temp.ParseFiles(t+"upload.html", t+"layout/header.html", t+"layout/footer.html", t+"archive.html", t+"presentation.html", t+"config.html", t+"layout/configMenu.html", t+"timeConfig.html", t+"timeConfigEdit.html", t+"index.html"))
}

//Context is the type used for passing data to handlers
type Context struct {
	Ac   appengine.Context
	W    http.ResponseWriter
	R    *http.Request
	Vars map[string]string
}

//Handler maps standard net/http handlers to handlers accepting Context
func Handler(hand func(Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ac := appengine.NewContext(r)
		vars := mux.Vars(r)
		hand(Context{Ac: ac, W: w, R: r, Vars: vars})
	}
}

//Log500 sends an Internal Server Error to user with error message from the error.
func Log500(err error, c Context) {
	c.Ac.Warningf("Error 500. %v", err)
	http.Error(c.W, err.Error(), http.StatusInternalServerError)
}

//Log404 sends a Not Found Error to user with error message from the error.
func Log404(err error, c Context) {
	c.Ac.Infof("Error 404. %v", err)
	http.Error(c.W, err.Error(), http.StatusNotFound)
}

//RenderLayout inserts template with given name into the layout and sets the title and pipeline.
//The template should be loaded inside templates variable
//If any arguments are provided after the context, they will be treated like links
//to JavaScript scripts to load in the header of the template.
func RenderLayout(tmpl string, title string, data interface{}, c Context, jsIncludes ...string) {
	renderTemplate("header.html", struct {
		Title      string
		JsIncludes []string
		Admin      bool
	}{title, jsIncludes, user.IsAdmin(c.Ac)}, c)
	renderTemplate(tmpl, data, c)
	renderTemplate("footer.html", nil, c)
}

//renderTemplate renders a single template
func renderTemplate(tmpl string, data interface{}, c Context) {
	if err := templates.ExecuteTemplate(c.W, tmpl, data); err != nil {
		Log500(err, c)
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

// Round return rounded version of x with prec precision.
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
func NormalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, Tz)
}

//NormalizeTime strips the date part from time.Date leaving only
//hours, minutes, seconds and nanoseconds.
func NormalizeTime(t time.Time) time.Time {
	return time.Date(1, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), Tz)
}
