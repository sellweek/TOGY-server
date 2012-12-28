//Contains routing information.
package main

import (
	"controllers"
	"github.com/gorilla/mux"
	"net/http"
	"util"
)

//init sets up routes for Google App Engine
func init() {
	r := mux.NewRouter()
	//Administrative routes
	r.HandleFunc("/admin", util.Handler(controllers.Admin))
	r.HandleFunc("/admin/presentation/upload", util.Handler(controllers.Upload)).Methods("GET")
	r.HandleFunc("/admin/presentation/upload", util.Handler(controllers.UploadHandler)).Methods("POST")
	r.HandleFunc("/admin/presentation/archive", util.Handler(controllers.Archive))
	r.HandleFunc("/admin/presentation/activate", util.Handler(controllers.Activate)).Methods("POST")
	r.HandleFunc("/admin/presentation/delete", util.Handler(controllers.Delete)).Methods("POST")
	r.HandleFunc("/admin/presentation/{id}", util.Handler(controllers.Presentation))
	r.HandleFunc("/admin/config", util.Handler(controllers.ShowConfig)).Methods("GET")
	r.HandleFunc("/admin/config", util.Handler(controllers.SetConfig)).Methods("POST")
	r.HandleFunc("/admin/config/timeOverride", util.Handler(controllers.TimeOverride))
	r.HandleFunc("/admin/config/timeOverride/edit/{id}", util.Handler(controllers.TimeOverrideEdit)).Methods("GET")
	r.HandleFunc("/admin/config/timeOverride/edit", util.Handler(controllers.TimeOverrideEdit)).Methods("GET")
	r.HandleFunc("/admin/config/timeOverride/edit/{id}", util.Handler(controllers.TimeOverrideSubmit)).Methods("POST")
	r.HandleFunc("/admin/config/timeOverride/edit/", util.Handler(controllers.TimeOverrideSubmit)).Methods("POST")
	r.HandleFunc("/admin/config/timeOverride/delete", util.Handler(controllers.TimeOverrideDelete)).Methods("POST")

	//Auxiliary admin routes
	r.HandleFunc("/admin/bootstrap", util.Handler(controllers.Bootstrap))
	r.HandleFunc("/admin/migrate", util.Handler(controllers.Migrate))

	//Client API routes
	//For download and downloadComplete actions of presentations
	//"active" can be used instead of key to select the currently active presentation
	r.HandleFunc("/api/presentation/{key}/download", util.Handler(controllers.Download))
	r.HandleFunc("/api/presentation/{key}/downloadComplete", util.Handler(controllers.DownloadFinish))
	r.HandleFunc("/api/presentation/{key}/description", util.Handler(controllers.GetDescription)).Methods("GET")
	r.HandleFunc("/api/presentation/{key}/description", util.Handler(controllers.UpdateDescription)).Methods("POST")
	r.HandleFunc("/api/presentation/{key}/name", util.Handler(controllers.GetName)).Methods("GET")
	r.HandleFunc("/api/presentation/{key}/name", util.Handler(controllers.UpdateName)).Methods("POST")
	r.HandleFunc("/api/config/download", util.Handler(controllers.GetConfig))
	r.HandleFunc("/api/cofig/downloadComplete", util.Handler(controllers.GotConfig))
	r.HandleFunc("/api/update", util.Handler(controllers.Update))
	r.StrictSlash(true)
	http.Handle("/", r)
}
