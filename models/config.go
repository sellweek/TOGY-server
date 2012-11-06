package models

import (
	"appengine"
	"appengine/datastore"
)

//Config stores central configuration in Datastore.
//There is always only one Config record in Datastore.
type Config struct {
	Data []byte //Stores JSON that is sent to clients as centralConfig
}

//GetKey returns key of the Datastore Config record.
func (c Config) GetKey(ctx appengine.Context) (k *datastore.Key, err error) {
	k, err = datastore.NewQuery("Config").KeysOnly().Run(ctx).Next(nil)
	return
}

//SaveConfig saves data provided into the Config record.
func SaveConfig(newData []byte, c appengine.Context) (err error) {
	var key *datastore.Key
	key, err = datastore.NewQuery("Config").KeysOnly().Run(c).Next(nil)
	if err == datastore.Done {
		key = datastore.NewIncompleteKey(c, "Config", nil)
	} else if err != nil {
		return
	}
	_, err = datastore.Put(c, key, &Config{newData})
	if err != nil {
		return
	}
	DeleteQueryTimesFor(Config{}, c)

	return
}

//GetConfig fetches the Config record from Datastore and
//returns its data.
func GetConfig(c appengine.Context) (data []byte, err error) {
	var conf Config
	_, err = datastore.NewQuery("Config").Run(c).Next(&conf)
	if err != nil {
		return
	}
	return conf.Data, nil
}
