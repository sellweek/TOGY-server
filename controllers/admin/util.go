package admin

import (
	"appengine/datastore"
	"fmt"
	"models/configuration/config"
	"models/presentation"
	"util"
)

//Bootstrap inserts fake presentation and config into datastore.
//Used when the system doesn't have any presentation inserted.
func Bootstrap(c util.Context) (err error) {
	p := presentation.New("test", "xxx", "DO NOT USE!", []byte("This is just a bootstrap presentation that can't be downloaded"), true)
	_, err = datastore.Put(c.Ac, datastore.NewIncompleteKey(c.Ac, "Presentation", nil), p)
	if err != nil {
		fmt.Fprintln(c.W, "Error with presentation: ", err)
		return nil
	}

	//	zeroTime := time.Date(0001, 01, 01, 00, 00, 00, 00, utc)

	conf := new(config.Config)

	err = conf.Save(c.Ac)

	if err != nil {
		fmt.Fprintln(c.W, "Error with config:", err)
		return nil
	}
	fmt.Fprint(c.W, "Do not start any clients until you have replaced this presentation.")
	return
}

//Migrate migrates Datastore data from a previous version.
func Migrate(c util.Context) (err error) {
	keys, err := datastore.NewQuery("Action").KeysOnly().GetAll(c.Ac, nil)
	if err != nil {
		return err
	}
	err = datastore.DeleteMulti(c.Ac, keys)
	if err != nil {
		return
	}
	fmt.Fprintf(c.W, "Success, %d actions deleted", len(keys))
	return
}
