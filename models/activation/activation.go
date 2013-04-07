package activation

import (
	"appengine"
	"appengine/datastore"
	"time"
)

type Activation struct {
	Time         time.Time
	Presentation *datastore.Key

	Key string `datastore:"-"`
}

func New(t time.Time, p *datastore.Key) (a *Activation) {
	a = new(Activation)
	a.Time = t
	a.Presentation = p
	return
}

func Make(t time.Time, p *datastore.Key, c appengine.Context) (a *Activation, err error) {
	a = New(t, p)
	err = a.Save(c)
	return
}

func GetByKey(k string, c appengine.Context) (a *Activation, err error) {
	dk, err := datastore.DecodeKey(k)
	if err != nil {
		return
	}

	a = new(Activation)
	err = datastore.Get(c, dk, a)
	a.Key = k
	return

}

func GetForPresentation(p *datastore.Key, c appengine.Context) (as []*Activation, err error) {
	as = make([]*Activation, 0)

	keys, err := datastore.NewQuery("Activation").Ancestor(p).GetAll(c, &as)
	if err != nil {
		return
	}

	for i, k := range keys {
		as[i].Key = k.Encode()
	}
	return
}

func GetBeforeTime(t time.Time, c appengine.Context) ([]*Activation, error) {
	return timeQuery(t, "<", c)
}

func GetAfterTime(t time.Time, c appengine.Context) ([]*Activation, error) {
	return timeQuery(t, ">", c)
}

func (a *Activation) Save(c appengine.Context) (err error) {
	var k *datastore.Key
	if a.Key == "" {
		var pKey *datastore.Key
		pKey, err = datastore.DecodeKey(a.Presentation.Encode())
		if err != nil {
			return
		}

		k, err = datastore.Put(c, datastore.NewIncompleteKey(c, "Activation", pKey), a)
		if err != nil {
			return
		}
		a.Key = k.Encode()
	} else {
		k, err = datastore.DecodeKey(a.Key)
		if err != nil {
			return
		}

		_, err = datastore.Put(c, k, a)
		return
	}
	return
}

func (a *Activation) Delete(c appengine.Context) (err error) {
	k, err := datastore.DecodeKey(a.Key)
	if err != nil {
		return
	}

	err = datastore.Delete(c, k)
	if err != nil {
		return
	}
	a.Key = ""

	return
}

func timeQuery(t time.Time, sign string, c appengine.Context) (as []*Activation, err error) {
	as = make([]*Activation, 0)
	keys, err := datastore.NewQuery("Activation").Filter("Time "+sign, t).Order("Time").GetAll(c, &as)
	if err != nil {
		return
	}

	for i, k := range keys {
		as[i].Key = k.Encode()
	}
	return
}
