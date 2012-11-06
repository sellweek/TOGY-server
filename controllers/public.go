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
//Clients insert the timecode of their active presentation into the "d" querystring parameter.
//	http://togy.appspot.com?d=20121012150456
//Server responds either with:
//	false
//if theit presentation is active and doesn't need updating or with
//	true,20121013105512.ppt
//where the part after comma is the filename of the presentation that the client should download.
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

	conf, err := models.WasDownloadedBy(models.Config{}, client, c.Ac)
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
		models.LogQueryTime(models.Config{}, client, models.UpdateNotification, c.Ac)
	}

}

//Serves a presentation from blobstore.
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

func DownloadFinish(c util.Context) {
	p, err := models.GetActive(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	models.LogQueryTime(*p, c.R.FormValue("client"), models.DownloadFinish, c.Ac)
}

func GetConfig(c util.Context) {
	conf, err := models.GetConfig(c.Ac)
	if err != nil {
		util.Log500(err, c)
	}
	fmt.Fprint(c.W, string(conf))
	models.LogQueryTime(models.Config{}, c.R.FormValue("client"), models.DownloadStart, c.Ac)
}

func GotConfig(c util.Context) {
	models.LogQueryTime(models.Config{}, c.R.FormValue("client"), models.DownloadFinish, c.Ac)
}
