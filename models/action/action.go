package action

import (
	"appengine"
	"appengine/datastore"
	"time"
)

const (
	UpdateNotification ActionType = iota //Client was notified it should update.
	DownloadStart      ActionType = iota //Client started download.
	DownloadFinish     ActionType = iota //Client finished download.
)

//ActionType defines the type of the Action performed
type ActionType int

//String returns the type of action in human-readable form
func (at ActionType) String() string {
	switch at {
	case UpdateNotification:
		return "Update notification"
	case DownloadStart:
		return "Download started"
	case DownloadFinish:
		return "Download is finished"
	}
	return "Unknown action"
}

//Action is a type used for recording the time when clients
//did an action specified by the ActionType.
type Action struct {
	Type   ActionType     //Type of action
	Client string         //Name of the client
	Time   time.Time      //Time of the action
	Model  *datastore.Key //What object is the action related to.
	Key    string         `datastore:"-"`
}

//Model is an interface specifying models - structs stored in Datastore.
type Model interface {
	GetKey(appengine.Context) (*datastore.Key, error)
}

//GetKey is used to obtain Model's key.
func (a Action) GetKey() (k *datastore.Key, err error) {
	k, err = datastore.DecodeKey(a.Key)
	return
}

//New returns a pointer to an Action with its fields set according to arguments.
func New(k *datastore.Key, at ActionType, client string) (a *Action) {
	a = new(Action)
	a.Model = k
	a.Type = at
	a.Client = client
	a.Time = time.Now()
	return a
}

//Make creates a new Action using New and then saves it to Datastore.
func Make(m Model, at ActionType, client string, c appengine.Context) (a *Action, err error) {
	k, err := m.GetKey(c)
	if err != nil {
		return
	}
	a = New(k, at, client)
	err = a.Save(c)
	return
}

//Save saves an Action to Datastore.
//If its Key field is set, it will replace existing record
//that has that key. If not, it will use datastore.NewIncompleteKey()
//to create a new key and set the field.
func (a *Action) Save(c appengine.Context) (err error) {
	if a.Key == "" {
		var key *datastore.Key
		key, err = datastore.Put(c, datastore.NewIncompleteKey(c, "Action", a.Model), a)
		a.Key = key.Encode()
		return
	} else {
		var key *datastore.Key
		key, err = datastore.DecodeKey(a.Key)
		if err != nil {
			return
		}
		_, err = datastore.Put(c, key, a)
	}
	return
}

//Log works like Make but logs errors instead of returning them.
func Log(m Model, client string, at ActionType, c appengine.Context) {
	if client == "" {
		c.Infof("%v called without client name.", at)
		return
	}

	_, err := Make(m, at, client, c)
	if err != nil {
		c.Infof("Can't log Action to Datastore: %v", err)
	}
}

//GetFor returns a slice containing all the Actions for a given Model.
func GetFor(m Model, c appengine.Context) (as []Action, err error) {
	key, err := m.GetKey(c)
	if err != nil {
		return
	}
	as = make([]Action, 12)
	keys, err := datastore.NewQuery("Action").Ancestor(key).GetAll(c, &as)
	if err != nil {
		return
	}
	for i := range keys {
		as[i].Key = keys[i].Encode()
	}
	return
}

//GetCountFor returns how many times given ActionType was performed on a Model.
func GetCountFor(at ActionType, m Model, c appengine.Context) (count int, err error) {
	key, err := m.GetKey(c)
	if err != nil {
		return
	}
	count, err = datastore.NewQuery("Action").Ancestor(key).Filter("Type =", at).Count(c)
	return
}

//WasPerformedOn returns whether given client performed given ActionType on gived Model.
func WasPerformedOn(at ActionType, m Model, client string, c appengine.Context) (bool, error) {
	key, err := m.GetKey(c)
	if err != nil {
		return false, err
	}
	i := datastore.NewQuery("Action").Ancestor(key).Filter("Client =", client).Filter("Type =", at).KeysOnly().Run(c)
	_, err = i.Next(nil)
	if err == datastore.Done {
		return false, nil
	}
	return true, err
}

//DeleteFor deletes all Actions for a specified Model.
func DeleteFor(m Model, c appengine.Context) (err error) {
	var keys []*datastore.Key
	key, err := m.GetKey(c)
	if err != nil {
		return
	}
	keys, err = datastore.NewQuery("Action").Ancestor(key).KeysOnly().GetAll(c, nil)
	if err != nil {
		return
	}
	err = datastore.DeleteMulti(c, keys)
	if err != nil {
		return
	}
	return
}
