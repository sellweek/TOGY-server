package models

import (
	"appengine"
	"appengine/datastore"
)

type Config struct {
	Data []byte
}

func (c Config) GetKey(ctx appengine.Context) (k *datastore.Key, err error) {
	k, err = datastore.NewQuery("Config").KeysOnly().Run(ctx).Next(nil)
	return
}

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

func GetConfig(c appengine.Context) (data []byte, err error) {
	var conf Config
	_, err = datastore.NewQuery("Config").Run(c).Next(&conf)
	if err != nil {
		return
	}
	return conf.Data, nil
}
