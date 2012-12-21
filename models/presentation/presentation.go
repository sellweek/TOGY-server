//Package models provides models used in application.
package presentation

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"fmt"
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

//GetActive gets the active presentation from the Datastore.
//If more than one presentation has the Active field set to true,
//its results are unpredictable.
func GetActive(c appengine.Context) (*Presentation, error) {
	var p Presentation
	i := datastore.NewQuery("Presentation").Filter("Active =", true).Limit(1).Run(c)
	key, err := i.Next(&p)
	if err != nil {
		return &p, err
	}
	p.Key = key.Encode()
	return &p, err
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

//Sve saves presentation to Datastore. If its Active field is true,
//the currently active presentation is deactivated.
//TODO: Use transactions to deactivate and activate presentations.
func (p *Presentation) Save(c appengine.Context) (err error) {
	var k *datastore.Key
	if p.Key == "" {

		if p.Active {
			err = deactivateCurrent(c)
			if err != nil {
				return
			}
		}
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
			err = deactivateCurrent(c)
			if err != nil {
				return
			}
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

//deactivateCurrent is a helper function that deactivates current active presentation.
//It should only be used just before saving another active presentation.
func deactivateCurrent(c appengine.Context) (err error) {
	curr, err := GetActive(c)
	if err != nil {
		err = fmt.Errorf("Couldn't get active presentation: %v", err)
		return
	}
	curr.Active = false
	err = curr.Save(c)
	if err != nil {
		err = fmt.Errorf("Couldn't deactivate active presentation: %v", err)
		return
	}
	return
}
