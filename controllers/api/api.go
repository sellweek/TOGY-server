//Package api provides controllers for the API
package api

import (
	"appengine/blobstore"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"github.com/russross/blackfriday"
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

//Status returns information about active broadcasts and
//config. It returns a JSON like:
//	{
//		"Broadcasts": [
//			{
//				"Key": "aghkZXZ-Tm9uZXIZCxIMUHJlc2VudGF0aW9uGICAgICAgIAKDA",
//				"FileType": "pptx"
//			},
//			{
//				"Key": "aghkZXZ-Tm9uZXIZCxIMUHJlc2VudGF0aW9uGICAgICAwO8KDA",
//				"FileType": "ppt"
//			}
//		]
//		"Config": 1383212550
//	}
func Status(c util.Context) (err error) {
	type broadcastInfo struct {
		Key      string
		FileType string
	}
	type updateInfo struct {
		Broadcasts []broadcastInfo
		Config     int64
	}

	ui := updateInfo{}

	ps, err := presentation.GetActive(c.Ac)
	if err != nil {
		return
	}

	ui.Broadcasts = make([]broadcastInfo, len(ps))

	for i, p := range ps {
		ui.Broadcasts[i] = broadcastInfo{Key: p.Key().Encode(), FileType: p.FileType}
	}

	conf, err := config.Get(c.Ac)
	if err != nil {
		return
	}

	ui.Config = conf.Timestamp

	data, err := json.Marshal(ui)
	if err != nil {
		return
	}

	fmt.Fprint(c.W, string(data))

	return
}

//Download serves the broadcast from blobstore.
func Download(c util.Context) (err error) {
	p, err := getPresentation(c)
	if err != nil {
		return
	}

	blobstore.Send(c.W, p.BlobKey)
	return
}

//Activated is called by clients to announce that
//they have activated the broadcast.
func Activated(c util.Context) (err error) {
	p, err := getPresentation(c)
	if err != nil {
		return
	}

	action.Log(p, action.Activated, c.R.FormValue("client"), c.Ac)
	return
}

//Deactivated is called by clients to announce that
//they have deactivated and deleted a broadcast
func Deactivated(c util.Context) (err error) {
	p, err := getPresentation(c)
	if err != nil {
		return
	}

	action.Log(p, action.Deactivated, c.R.FormValue("client"), c.Ac)
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
	return
}

//GotConfig is called by clients to announce that
//they have downloaded the broadcast.
func GotConfig(c util.Context) (err error) {
	conf, err := config.Get(c.Ac)
	if err != nil {
		return
	}

	action.Log(conf, action.Activated, c.R.FormValue("client"), c.Ac)
	return
}

func ScheduleActivation(c util.Context) (err error) {
	p, err := getPresentation(c)
	if err != nil {
		return
	}

	defer c.R.Body.Close()
	data, err := ioutil.ReadAll(c.R.Body)
	if err != nil {
		return
	}

	params := make(map[string]string)
	err = json.Unmarshal(data, &params)
	if err != nil {
		return
	}

	t, err := time.ParseInLocation("2006-01-02 15:04", params["time"], util.C.Tz)
	if err != nil {
		return
	}

	var op activation.Operation
	if params["operation"] == "activate" {
		op = activation.Activate
	} else {
		op = activation.Deactivate
	}

	_, err = activation.Make(op, t, p.Key(), c.Ac)
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

	for _, a := range as {
		var p *presentation.Presentation
		p, err = presentation.GetByKey(a.Presentation, c.Ac)
		if err != nil {
			c.Ac.Errorf("Couldn't load presentation %s", a.Presentation.Encode())
			continue
		}

		if a.Op == activation.Activate {
			p.Active = true
		} else {
			p.Active = false
		}

		err = p.Save(c.Ac)
		if err != nil {
			c.Ac.Errorf("Couldn't save presentation: %v", p.Name)
			continue
		}

		err = a.Delete(c.Ac)
		if err != nil {
			c.Ac.Errorf("Couldn't remove activation: %s", a.Key)
			continue
		}
	}
	return
}

func DeleteActivation(c util.Context) (err error) {
	k, err := datastore.DecodeKey(c.Vars["key"])
	if err != nil {
		return
	}

	a, err := activation.GetByKey(k, c.Ac)
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
	key, err := datastore.DecodeKey(c.Vars["key"])
	if err != nil {
		return
	}

	p, err = presentation.GetByKey(key, c.Ac)
	return
}
