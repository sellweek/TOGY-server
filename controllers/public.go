//Package controllers provides functions that handle the requests and respond to them.
package controllers

import (
	"appengine/blobstore"
	"encoding/json"
	"fmt"
	"models"
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

	p, err := models.GetActive(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}

	bc, err := models.WasDownloadedBy(p, client, c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	ui.Broadcast = !bc

	ui.FileType = p.FileType

	conf, err := models.WasDownloadedBy(new(models.Config), client, c.Ac)
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
		models.LogQueryTime(*p, client, models.UpdateNotification, c.Ac)
	}
	if ui.Config {
		models.LogQueryTime(new(models.Config), client, models.UpdateNotification, c.Ac)
	}

}

//Serves the broadcast from blobstore.
func Download(c util.Context) {
	id := c.R.FormValue("id")
	var p *models.Presentation
	var err error
	if id == "" {
		p, err = models.GetActive(c.Ac)
		models.LogQueryTime(*p, c.R.FormValue("client"), models.DownloadStart, c.Ac)
	} else {
		p, err = models.GetByKey(id, c.Ac)
	}
	if err != nil {
		util.Log500(err, c)
		return
	}
	blobstore.Send(c.W, p.BlobKey)
}

//Clients call it with ther ID to inform the server
//that they have finished downloading the broadcast.
func DownloadFinish(c util.Context) {
	p, err := models.GetActive(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	models.LogQueryTime(*p, c.R.FormValue("client"), models.DownloadFinish, c.Ac)
}

//Serves the configuration.
func GetConfig(c util.Context) {
	json, err := models.ConfigJSON(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	fmt.Fprint(c.W, string(json))
	if client := c.R.FormValue("client"); client != "" {
		models.LogQueryTime(&models.Config{}, client, models.DownloadStart, c.Ac)
	}
}

//Used by clients in the same manner as DownloadFinish to inform 
//that they have downloaded the configuration file.
func GotConfig(c util.Context) {
	models.LogQueryTime(new(models.Config), c.R.FormValue("client"), models.DownloadFinish, c.Ac)
}
