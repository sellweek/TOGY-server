package admin

import (
	"appengine/datastore"
	"appengine/user"
	"fmt"
	"net/http"
	"util"
)

func Logout(c util.Context) (err error) {
	url, err := user.LogoutURL(c, "/")
	if err != nil {
		return
	}

	http.Redirect(c.W, c.R, url, 303)
	return
}

//Migrate migrates Datastore data from a previous version.
func Migrate(c util.Context) (err error) {
	keys, err := datastore.NewQuery("Action").KeysOnly().GetAll(c, nil)
	if err != nil {
		return err
	}
	err = datastore.DeleteMulti(c, keys)
	if err != nil {
		return
	}
	fmt.Fprintf(c.W, "Success, %d actions deleted", len(keys))
	return
}
