package activation

import (
	"appengine"
	"appengine/datastore"
	"github.com/sellweek/gaemodel"
	"time"
)

const (
	Deactivate Operation = iota
	Activate   Operation = iota
)

type Operation int

type Activation struct {
	Op           Operation
	Time         time.Time
	Presentation *datastore.Key

	key *datastore.Key `datastore:"-"`
}

var typ = reflect.TypeOf(Activation{})

func (a *Activation) Key() *datastore.Key {
	return a.key
}

func (a *Activation) SetKey(k *datastore.Key) {
	a.key = k
}

func (a *Activation) Kind() string {
	return "Activation"
}

func (a *Activation) Ancestor() *datastore.Key {
	return a.Presentation
}

func New(op Operation, t time.Time, p *datastore.Key) (a *Activation) {
	a = new(Activation)
	a.Op = op
	a.Time = t
	a.Presentation = p
	return
}

func Make(op Operation, t time.Time, p *datastore.Key, c appengine.Context) (a *Activation, err error) {
	a = New(op, t, p)
	err = a.Save(c)
	return
}

func GetByKey(k *datastore.Key, c appengine.Context) (a *Activation, err error) {
	a = new(Activation)
	err = datastore.Get(c, dk, a)
	a.Key = k
	return
}

func GetForPresentation(p *datastore.Key, c appengine.Context) (as []*Activation, err error) {
	is, err := gaemodel.GetByAncestor(c, typ, "Activation", p)
	if err != nil {
		return
	}
	as = is.([]*Activation)
	return
}

func GetBeforeTime(t time.Time, c appengine.Context) ([]*Activation, error) {
	return timeQuery(t, "<", c)
}

func GetAfterTime(t time.Time, c appengine.Context) ([]*Activation, error) {
	return timeQuery(t, ">", c)
}

func (a *Activation) Save(c appengine.Context) error {
	return gaemodel.Save(c, a)
}

func (a *Activation) Delete(c appengine.Context) (err error) {
	gaemodel.Delete(c, a)
	if err != nil {
		return
	}

	a.Key = ""

	return
}

func timeQuery(t time.Time, sign string, c appengine.Context) (as []*Activation, err error) {
	q := datastore.NewQuery("Activation").Filter("Time "+sign, t).Order("Time")
	is, err := gaemodel.MultiQuery(c, typ, "Activation", q)
	if err != nil {
		return
	}
	as = is.([]*Activation)
	return
}
