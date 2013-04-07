package activation

import (
	"appengine"
	"appengine/datastore"
	"models/presentation"
	"time"
)

type Activation struct {
	Time         time.Time
	Presentation presentation.Presentation

	Key string `datastore:"-"`
}

func New(t time.Time, p presentation.Presentation) (a *Activation) {
	a = new(Activation)
	a.Time = t
	a.Presentation = p
	return
}

func Make(t time.Time, p presentation.Presentation, c appengine.Context) (a *Activation, err error) {
	a = New(t, p)
	err = a.Save(c)
	return
}

func (a *Activation) Save(c appengine.Context) (err error) {
	var k *datastore.Key
	if a.Key == "" {
		var pKey *datastore.Key
		pKey, err = datastore.DecodeKey(a.Presentation.Key)
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
