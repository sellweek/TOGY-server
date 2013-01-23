package admin

import (
	"appengine/datastore"
	"fmt"
	"models/configuration/config"
	"models/presentation"
	"util"
)

//Bootstrap inserts fake presentation and config into datastore.
//Used when the system doesn't have any presentation inserted
//and is in an inconsistent state because of that.
func Bootstrap(c util.Context) {
	p := presentation.New("test", "xxx", "DO NOT USE!", []byte("This is just a bootstrap presentation that can't be downloaded"), true)
	_, err := datastore.Put(c.Ac, datastore.NewIncompleteKey(c.Ac, "Presentation", nil), p)
	if err != nil {
		fmt.Fprintln(c.W, "Error with presentation: ", err)
	}

	//	zeroTime := time.Date(0001, 01, 01, 00, 00, 00, 00, utc)

	conf := new(config.Config)

	err = conf.Save(c.Ac)

	if err != nil {
		fmt.Fprintln(c.W, "Error with config:", err)
	}
	fmt.Fprint(c.W, "Do not start any clients until you have replaced this presentation.")
}

//Migrate migrates Datastore data from a previous version.
func Migrate(c util.Context) {
	fmt.Fprintf(c.W, "There is nothing to migrate in current version.")
}
