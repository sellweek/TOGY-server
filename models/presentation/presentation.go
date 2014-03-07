package presentation

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"github.com/sellweek/gaemodel"
	"models/action"
	"reflect"
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

	key *datastore.Key `datastore:"-"`
}

var typ = reflect.TypeOf(Presentation{})

func (p *Presentation) Key() *datastore.Key {
	return p.key
}

func (p *Presentation) SetKey(k *datastore.Key) {
	p.key = k
}

func (p *Presentation) Kind() string {
	return "Presentation"
}

func (_ *Presentation) Ancestor() *datastore.Key {
	return nil
}

//GetActive gets the active presentations from the Datastore.
func GetActive(c appengine.Context) (ps []*Presentation, err error) {
	q := datastore.NewQuery("Presentation").Filter("Active =", true)
	is, err := gaemodel.MultiQuery(c, typ, "Presentation", q)
	if err != nil {
		return
	}
	ps = is.([]*Presentation)
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
	return gaemodel.Save(c, p)
}

//Delete deletes the presentation record from Datastore and
//its data file from Blobstore.
func (p *Presentation) Delete(c appengine.Context) (err error) {
	err = gaemodel.Delete(c, p)
	if err != nil {
		return
	}

	err = blobstore.Delete(c, p.BlobKey)
	if err != nil {
		return
	}
	err = action.DeleteFor(p, c)
	return
}

//GetListing gets paginated Presentations from Datastore.
func GetListing(page int, perPage int, c appengine.Context) (ps []*Presentation, err error) {
	is, err := gaemodel.GetAll(c, typ, "Presentation", page, perPage)
	if err != nil {
		return
	}
	ps = is.([]*Presentation)
	return
}

//PageCount returns how many pages the Presentations in Datastore would need
//if there wer perPage presentations listed on a single page
func PageCount(perPage int, c appengine.Context) (int, error) {
	return gaemodel.PageCount(c, "Presentation", perPage)
}

//GetByKey gets the presentation with given key from Datastore and returns a pointer to it.
func GetByKey(key *datastore.Key, c appengine.Context) (p *Presentation, err error) {
	i, err := gaemodel.GetByKey(c, typ, key)
	if err != nil {
		return
	}
	p = i.(*Presentation)
	return
}
