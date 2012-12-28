//Package controllers provides functions that handle the requests and respond to them.
package controllers

import (
	"appengine/blobstore"
	"encoding/json"
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"models/action"
	"models/configuration"
	"models/configuration/config"
	"models/presentation"
	"util"
)

//Handles queries of clients about whether they should download a new
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
func Update(c util.Context) {
	type updateInfo struct {
		Broadcast bool
		FileType  string
		Config    bool
	}

	ui := new(updateInfo)

	client := c.R.FormValue("client")

	p, err := presentation.GetActive(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}

	bc, err := action.WasPerformedOn(action.DownloadFinish, p, client, c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	ui.Broadcast = !bc

	ui.FileType = p.FileType

	conf, err := action.WasPerformedOn(action.DownloadFinish, new(config.Config), client, c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	ui.Config = !conf

	data, err := json.Marshal(ui)
	if err != nil {
		util.Log500(err, c)
		return
	}

	fmt.Fprint(c.W, string(data))

	if ui.Broadcast {
		action.Log(*p, client, action.UpdateNotification, c.Ac)
	}
	if ui.Config {
		action.Log(new(config.Config), client, action.UpdateNotification, c.Ac)
	}

}

//Serves the broadcast from blobstore.
func Download(c util.Context) {
	p, err := getPresentation(c)
	if err != nil {
		util.Log500(err, c)
		return
	}
	if client := c.R.FormValue("client"); client != "" {
		action.Log(*p, c.R.FormValue("client"), action.DownloadStart, c.Ac)
	}
	blobstore.Send(c.W, p.BlobKey)
}

//Clients call it with ther ID to inform the server
//that they have finished downloading the broadcast.
func DownloadFinish(c util.Context) {
	p, err := getPresentation(c)
	if err != nil {
		util.Log500(err, c)
		return
	}

	//Here, we're using Make instead of Log because the sole purpose of this controller
	//is to log the action, so we want to see the errors.
	_, err = action.Make(*p, action.DownloadFinish, c.R.FormValue("client"), c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
}

func GetDescription(c util.Context) {
	p, err := getPresentation(c)
	if err != nil {
		util.Log500(err, c)
		return
	}
	fmt.Fprint(c.W, string(p.Description))
}

func UpdateDescription(c util.Context) {
	p, err := getPresentation(c)
	if err != nil {
		util.Log500(err, c)
		return
	}
	defer c.R.Body.Close()
	body, err := ioutil.ReadAll(c.R.Body)
	if err != nil {
		util.Log500(err, c)
		return
	}
	p.Description = body
	err = p.Save(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	fmt.Fprint(c.W, string(blackfriday.MarkdownCommon(body)))
}

func GetName(c util.Context) {
	p, err := getPresentation(c)
	if err != nil {
		util.Log500(err, c)
		return
	}
	fmt.Fprint(c.W, p.Name)
}

func UpdateName(c util.Context) {
	p, err := getPresentation(c)
	if err != nil {
		util.Log500(err, c)
		return
	}
	defer c.R.Body.Close()
	body, err := ioutil.ReadAll(c.R.Body)
	if err != nil {
		util.Log500(err, c)
		return
	}
	p.Name = string(body)
	err = p.Save(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
}

//Serves the configuration.
func GetConfig(c util.Context) {
	json, err := configuration.JSON(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	fmt.Fprint(c.W, string(json))
	if client := c.R.FormValue("client"); client != "" {
		action.Log(&config.Config{}, client, action.DownloadStart, c.Ac)
	}
}

//Used by clients in the same manner as DownloadFinish to inform 
//that they have downloaded the configuration file.
func GotConfig(c util.Context) {
	action.Log(new(config.Config), c.R.FormValue("client"), action.DownloadFinish, c.Ac)
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
