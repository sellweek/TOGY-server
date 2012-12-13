package action

import (
	"appengine"
	"appengine/datastore"
	"time"
)

const (
	UpdateNotification ActionType = iota //Client was notified it should update.
	DownloadStart      ActionType = iota //Client startde downloading.
	DownloadFinish     ActionType = iota //Client finished download.
)

type ActionType int

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

//QueryTime (better name would be Action, but it would be hard to change it now)
//is a type used for recording whe time when clients
//did an action specified by the ActionType.
type QueryTime struct {
	Action       ActionType     //The action client did
	Client       string         //Name of the client
	Time         time.Time      //Time of the action
	Presentation *datastore.Key //What object is the action related to.
	Key          string         `datastore:"-"`
}

//Model is an nterface specifying models - structs which Datastore keys can be obtained.
type Model interface {
	GetKey(appengine.Context) (*datastore.Key, error)
}

func (qt QueryTime) GetKey() (k *datastore.Key, err error) {
	k, err = datastore.DecodeKey(qt.Key)
	return
}

//NewQueryTime returns a pointer to a QueryTime with its fields set according
//to arguments.
func NewQueryTime(k *datastore.Key, at ActionType, client string) (qt *QueryTime) {
	qt = new(QueryTime)
	qt.Presentation = k
	qt.Action = at
	qt.Client = client
	qt.Time = time.Now()
	return qt
}

//MakeQueryTime creates a new QueryTime using NewQueryTime and then saves
//it to Datastore.
func MakeQueryTime(m Model, at ActionType, client string, c appengine.Context) (qt *QueryTime, err error) {
	k, err := m.GetKey(c)
	if err != nil {
		return
	}
	qt = NewQueryTime(k, at, client)
	err = qt.Save(c)
	return
}

//Save saves a QuryTime to Datastore.
//If its Key field is set, it will use it, replacing
//existing records. If not, it will use datastore.NewIncompleteKey()
//to create a new key and set the field.
func (qt *QueryTime) Save(c appengine.Context) (err error) {
	if qt.Key == "" {
		var key *datastore.Key
		key, err = datastore.Put(c, datastore.NewIncompleteKey(c, "QueryTime", qt.Presentation), qt)
		qt.Key = key.Encode()
		return
	} else {
		var key *datastore.Key
		key, err = datastore.DecodeKey(qt.Key)
		if err != nil {
			return
		}
		_, err = datastore.Put(c, key, qt)
	}
	return
}

//Functions like MakeQueryTime but logs errors instead of returning them.
func LogQueryTime(m Model, client string, at ActionType, c appengine.Context) {
	if client == "" {
		c.Infof("%v called without client name.", at)
		return
	}

	_, err := MakeQueryTime(m, at, client, c)
	if err != nil {
		c.Infof("Can't log QueryTime to Datastore: %v", err)
	}
}

//GetQueryTimes returns a slice containing all the QueryTimes for
//a given Model.
func GetQueryTimes(m Model, c appengine.Context) (qts []QueryTime, err error) {
	key, err := m.GetKey(c)
	if err != nil {
		return
	}
	qts = make([]QueryTime, 12)
	keys, err := datastore.NewQuery("QueryTime").Ancestor(key).GetAll(c, &qts)
	if err != nil {
		return
	}
	for i := range keys {
		qts[i].Key = keys[i].Encode()
	}
	return
}

//GetDownloadCount returns how many times given Model was downloaded.
func GetDownloadCount(m Model, c appengine.Context) (count int, err error) {
	key, err := m.GetKey(c)
	if err != nil {
		return
	}
	count, err = datastore.NewQuery("QueryTime").Ancestor(key).Filter("Action =", DownloadFinish).Count(c)
	return
}

//WasDownloadedBy returns whether given client downloaded file associated
//with given Model.
func WasDownloadedBy(m Model, client string, c appengine.Context) (bool, error) {
	key, err := m.GetKey(c)
	if err != nil {
		return false, err
	}
	i := datastore.NewQuery("QueryTime").Ancestor(key).Filter("Client =", client).Filter("Action =", DownloadFinish).KeysOnly().Run(c)
	_, err = i.Next(nil)
	if err == datastore.Done {
		return false, nil
	}
	return true, err
}

//DeleteQueryTimesFor deletes all the QueryTimes for a specified Model.
func DeleteQueryTimesFor(m Model, c appengine.Context) (err error) {
	var keys []*datastore.Key
	key, err := m.GetKey(c)
	if err != nil {
		return
	}
	keys, err = datastore.NewQuery("QueryTime").Ancestor(key).KeysOnly().GetAll(c, nil)
	if err != nil {
		return
	}
	err = datastore.DeleteMulti(c, keys)
	if err != nil {
		return
	}
	return
}
