//Package models provides models used in application.
package presentation

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"fmt"
	"math"
	"models/action"
	"time"
)

//Presentation stores data about uploaded presentations in Datastore.
type Presentation struct {
	BlobKey     appengine.BlobKey //Key linking the metadata with file stored in Blobstore
	Created     time.Time         //Time of upload
	Name        string            //Name used to identify presentation in administration
	Description []byte            //Textual description used in UI
	FileType    string            //Name that clients should store the presentation under
	Active      bool              //If true, presentation is distributed to clients. Only one presentation can be active.

	//Used for storage of Datastore key when passing the struct around. Doesn't get saved to DS.
	//The key is stored as an encoded string
	Key string `datastore:"-"`
}

func (p Presentation) GetKey(ctx appengine.Context) (k *datastore.Key, err error) {
	k, err = datastore.DecodeKey(p.Key)
	return
}

//GetActive gets the active presentations from the Datastore.
func GetActive(c appengine.Context) (ps []*Presentation, err error) {
	ps = make([]*Presentation, 0)
	keys, err := datastore.NewQuery("Presentation").Filter("Active =", true).GetAll(c, &ps)
	if err != nil {
		return
	}

	for i, p := range ps {
		p.Key = keys[i].Encode()
	}
	return
}

//New returns pointer to a presentation with fields set to given values.
func New(k appengine.BlobKey, fileType, name string, desc []byte, active bool) *Presentation {
	p := new(Presentation)
	p.BlobKey = k
	p.Created = time.Now()
	p.FileType = fileType
	p.Name = name
	p.Description = desc
	p.Active = active
	return p
}

//Make creates a new presentation with New, saves it to Datastore and returns a pointer to it.
func Make(k appengine.BlobKey, fileType, name string, desc []byte, active bool, c appengine.Context) (p *Presentation, err error) {
	p = New(k, fileType, name, desc, active)
	err = p.Save(c)
	return
}

//Save saves presentation to Datastore.
func (p *Presentation) Save(c appengine.Context) (err error) {
	var k *datastore.Key
	if p.Key == "" {

		k, err = datastore.Put(c, datastore.NewIncompleteKey(c, "Presentation", nil), p)
		if err != nil {
			return
		}

		p.Key = k.Encode()
	} else {
		k, err = datastore.DecodeKey(p.Key)
		if err != nil {
			return
		}

		if p.Active {
			action.DeleteFor(p, c)
		}
		_, err = datastore.Put(c, k, p)
		if err != nil {
			return
		}

	}
	return
}

//Delete deletes the presentation record from Datastore and
//its data file from Blobstore.
func (p *Presentation) Delete(c appengine.Context) (err error) {
	if p.Active {
		err = fmt.Errorf("Active presentation can't be deleted.")
		return
	}
	k, err := datastore.DecodeKey(p.Key)
	if err != nil {
		return
	}

	err = datastore.Delete(c, k)
	if err != nil {
		return
	}
	p.Key = ""

	err = blobstore.Delete(c, p.BlobKey)

	return
}

//GetAll returns a slice of all the presentations in Datastore.
func GetAll(c appengine.Context) (ps []*Presentation, err error) {
	ps = make([]*Presentation, 0)
	keys, err := datastore.NewQuery("Presentation").Order("-Created").GetAll(c, &ps)
	for i := 0; i < len(ps); i++ {
		ps[i].Key = keys[i].Encode()
	}
	return
}

//GetListing gets paginated Presentations from Datastore.
func GetListing(page int, perPage int, c appengine.Context) (ps []*Presentation, err error) {
	var q *datastore.Query
	if page == 1 {
		q = datastore.NewQuery("Presentation").Limit(perPage).Order("-Active").Order("-Created")
	} else {
		q = datastore.NewQuery("Presentation").Limit(perPage).Offset(perPage * (page - 1)).Order("-Active").Order("-Created")
	}

	keys, err := q.GetAll(c, &ps)
	if err != nil {
		return
	}
	for i, p := range ps {
		p.Key = keys[i].Encode()
	}
	return
}

//PageCount returns how many pages the Presentations in Datastore would need
//if there wer perPage presentations listed on a single page
func PageCount(perPage int, c appengine.Context) (pgs int, err error) {
	ps, err := datastore.NewQuery("Presentation").Count(c)
	if err != nil {
		return
	}
	pgs = int(math.Ceil(float64(ps) / float64(perPage)))
	return
}

//GetByKey gets the presentation with given key from Datastore and returns a pointer to it.
func GetByKey(key string, c appengine.Context) (p *Presentation, err error) {
	k, err := datastore.DecodeKey(key)
	if err != nil {
		return
	}
	p = new(Presentation)
	err = datastore.Get(c, k, p)
	if err != nil {
		return
	}
	p.Key = key
	return
}
