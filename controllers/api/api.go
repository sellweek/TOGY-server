//Package api provides controllers for the API
package api

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"models/action"
	"models/activation"
	"models/configuration"
	"models/configuration/config"
	"models/presentation"
	"net/http"
	"time"
	"util"
)

//Update handles queries of clients about whether they should download a new
//presentation.
//Clients call address with their ID like
//	togpm5.appspot.com/?client=A1
//and the server responds with JSON:
//	{
//		"Broadcast": true,
//		"FileType": "ppt",
//		"Config": false
//	}
//Where Broadcast and Config fields signal whether client should download
//new presentation or configuration and FileType contains file type
//of the broadcast file.
func Update(c util.Context) (err error) {
	type updateInfo struct {
		Broadcast bool
		FileType  string
		Config    bool
	}

	ui := new(updateInfo)

	client := c.R.FormValue("client")

	p, err := presentation.GetActive(c.Ac)
	if err != nil {
		return
	}

	bc, err := action.WasPerformedOn(action.DownloadFinish, p, client, c.Ac)
	if err != nil {
		return
	}
	ui.Broadcast = !bc

	ui.FileType = p.FileType

	conf, err := action.WasPerformedOn(action.DownloadFinish, new(config.Config), client, c.Ac)
	if err != nil {
		return
	}
	ui.Config = !conf

	data, err := json.Marshal(ui)
	if err != nil {
		return
	}

	fmt.Fprint(c.W, string(data))

	if ui.Broadcast {
		action.Log(*p, client, action.UpdateNotification, c.Ac)
	}
	if ui.Config {
		action.Log(new(config.Config), client, action.UpdateNotification, c.Ac)
	}

	return
}

//Download serves the broadcast from blobstore.
func Download(c util.Context) (err error) {
	p, err := getPresentation(c)
	if err != nil {
		return
	}
	if client := c.R.FormValue("client"); client != "" {
		action.Log(*p, c.R.FormValue("client"), action.DownloadStart, c.Ac)
	}
	blobstore.Send(c.W, p.BlobKey)
	return
}

//DownloadFinish is called by clients to announce that
//they have downloaded the broadcast.
func DownloadFinish(c util.Context) (err error) {
	p, err := getPresentation(c)
	if err != nil {
		return
	}

	//Here, we're using Make instead of Log because the sole purpose of this controller
	//is to log the action, so we want to see the errors.
	_, err = action.Make(*p, action.DownloadFinish, c.R.FormValue("client"), c.Ac)
	if err != nil {
		return
	}
	return
}

//GetDescription responds with the description of a broadcast.
func GetDescription(c util.Context) (err error) {
	p, err := getPresentation(c)
	if err != nil {
		return
	}
	fmt.Fprint(c.W, string(p.Description))
	return
}

//UpdateDescription changes the description of a broadcast.
func UpdateDescription(c util.Context) (err error) {
	p, err := getPresentation(c)
	if err != nil {
		return
	}
	defer c.R.Body.Close()
	body, err := ioutil.ReadAll(c.R.Body)
	if err != nil {
		return
	}
	p.Description = body
	err = p.Save(c.Ac)
	if err != nil {
		return
	}
	fmt.Fprint(c.W, string(blackfriday.MarkdownCommon(body)))
	return
}

//GetName responds with the name of a broadcast.
func GetName(c util.Context) (err error) {
	p, err := getPresentation(c)
	if err != nil {
		return
	}
	fmt.Fprint(c.W, p.Name)

	return
}

//UpdateName changes the name of a broadcast.
func UpdateName(c util.Context) (err error) {
	p, err := getPresentation(c)
	if err != nil {
		return
	}
	defer c.R.Body.Close()
	body, err := ioutil.ReadAll(c.R.Body)
	if err != nil {
		return
	}
	p.Name = string(body)
	err = p.Save(c.Ac)
	if err != nil {
		return
	}
	return
}

//GetConfig serves the configuration.
func GetConfig(c util.Context) (err error) {
	json, err := configuration.JSON(c.Ac)
	if err != nil {
		return
	}
	fmt.Fprint(c.W, string(json))
	if client := c.R.FormValue("client"); client != "" {
		action.Log(&config.Config{}, client, action.DownloadStart, c.Ac)
	}
	return
}

//GotConfig is called by clients to announce that
//they have downloaded the broadcast.
func GotConfig(c util.Context) (err error) {
	action.Log(new(config.Config), c.R.FormValue("client"), action.DownloadFinish, c.Ac)
	return
}

func ScheduleActivation(c util.Context) (err error) {
	p, err := getPresentation(c)
	if err != nil {
		return
	}

	defer c.R.Body.Close()
	timeString, err := ioutil.ReadAll(c.R.Body)
	if err != nil {
		return
	}

	t, err := time.Parse("Mon Jan 2 2006 15:04:05 GMT-0700 (MST)", string(timeString))
	if err != nil {
		return
	}

	pk, err := datastore.DecodeKey(p.Key)
	if err != nil {
		return
	}

	_, err = activation.Make(t, pk, c.Ac)
	if err != nil {
		return
	}
	return
}

func ActivateScheduled(c util.Context) (err error) {
	t := time.Now()
	as, err := activation.GetBeforeTime(t, c.Ac)
	if err != nil {
		return
	}

	l := len(as)
	if l == 0 {
		return
	}

	for i, a := range as {
		if i == l-1 {
			break
		}
		err = a.Delete(c.Ac)
		if err != nil {
			c.Ac.Errorf("Error when deleting expired Activation: %v", err)
		}
	}

	pk := as[l-1].Presentation
	p, err := presentation.GetByKey(pk.Encode(), c.Ac)
	if err != nil {
		c.Ac.Errorf("Error when loading Presentation: %v", err)
		return
	}

	p.Active = true
	err = p.Save(c.Ac)
	if err != nil {
		c.Ac.Errorf("Error when activating Presentation: %v", err)
		return
	}

	err = as[l-1].Delete(c.Ac)
	if err != nil {
		c.Ac.Errorf("Error when deleting used Activation: %v", err)
		return
	}
	return
}

func DeleteActivation(c util.Context) (err error) {
	a, err := activation.GetByKey(c.Vars["key"], c.Ac)
	if err != nil {
		return
	}

	err = a.Delete(c.Ac)
	if err != nil {
		return
	}

	http.Redirect(c.W, c.R, c.R.FormValue("redirect"), 303)
	return
}

func getPresentation(c util.Context) (p *presentation.Presentation, err error) {
	key := c.Vars["key"]
	if key == "active" {
		p, err = presentation.GetActive(c.Ac)
	} else {
		p, err = presentation.GetByKey(key, c.Ac)
	}
	return
}

func ZipAll(c util.Context) (err error) {
	c.Ac.Infof("Getting presentations")
	ps, err := presentation.GetAll(c.Ac)
	if err != nil {
		return
	}

	fileNo := 1
	pNo := 0

	var (
		blob *blobstore.Writer
		r    blobstore.Reader
		z    *zip.Writer
	)

	for i, p := range ps {
		pNo++
		if i%12 == 0 {
			c.Ac.Infof("Creating blob %d", fileNo)
			blob, err = blobstore.Create(c.Ac, "application/zip")
			if err != nil {
				return
			}

			c.Ac.Infof("Creating zip")
			z = zip.NewWriter(blob)

			c.Ac.Infof("Reading presentation: " + p.Name)
			r = blobstore.NewReader(c.Ac, p.BlobKey)
			if err != nil {
				return
			}
		}

		c.Ac.Infof("Creating a file inside zip")
		var pw io.Writer
		pw, err = z.Create(fmt.Sprint(p.Name, ".", p.FileType))
		if err != nil {
			return
		}

		c.Ac.Infof("Writing to zip")
		_, err = io.Copy(pw, r)
		if err != nil {
			return
		}

		if i%11 == 0 && i != 0 {
			c.Ac.Infof("Closing zip file")
			err = z.Close()
			if err != nil {
				return
			}

			c.Ac.Infof("Closing blob")
			err = blob.Close()
			if err != nil {
				return
			}

			var key appengine.BlobKey
			key, err = blob.Key()
			if err != nil {
				return
			}

			c.Ac.Infof("%d/%d presentations zipped. The key is: %v", pNo, len(ps), key)

			fileNo++
		}
	}
	return
}
