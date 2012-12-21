//Contains routing information.
package main

import (
	"code.google.com/p/gorilla/mux"
	"controllers"
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
	r.HandleFunc("/api/download", util.Handler(controllers.Download))
	r.HandleFunc("/api/downloadComplete", util.Handler(controllers.DownloadFinish))
	r.HandleFunc("/api/config", util.Handler(controllers.GetConfig))
	r.HandleFunc("/api/gotConfig", util.Handler(controllers.GotConfig))
	r.HandleFunc("/api/update", util.Handler(controllers.Update))
	r.StrictSlash(true)
	http.Handle("/", r)
}
