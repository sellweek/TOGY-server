//Contains routing information.
package main

import (
	"controllers/admin"
	"controllers/api"
	"controllers/public"
	"github.com/gorilla/mux"
	"net/http"
	"util"
)

//init sets up routes for Google App Engine
func init() {
	r := mux.NewRouter()

	//Public routes
	r.HandleFunc("/", util.Handler(public.Index))
	r.HandleFunc("/presentation", util.Handler(public.Presentations))

	//Administrative routes
	r.HandleFunc("/admin", util.Handler(admin.Admin))
	r.HandleFunc("/admin/presentation/upload", util.Handler(admin.Upload)).Methods("GET")
	r.HandleFunc("/admin/presentation/upload", util.Handler(admin.UploadHandler)).Methods("POST")
	r.HandleFunc("/admin/presentation/archive", util.Handler(admin.Archive))
	r.HandleFunc("/admin/presentation/activate", util.Handler(admin.Activate)).Methods("POST")
	r.HandleFunc("/admin/presentation/delete", util.Handler(admin.Delete)).Methods("POST")
	r.HandleFunc("/admin/presentation/{id}", util.Handler(admin.Presentation))
	r.HandleFunc("/admin/config", util.Handler(admin.ShowConfig)).Methods("GET")
	r.HandleFunc("/admin/config", util.Handler(admin.SetConfig)).Methods("POST")
	r.HandleFunc("/admin/config/timeOverride", util.Handler(admin.TimeOverride))
	r.HandleFunc("/admin/config/timeOverride/edit/{id}", util.Handler(admin.TimeOverrideEdit)).Methods("GET")
	r.HandleFunc("/admin/config/timeOverride/edit", util.Handler(admin.TimeOverrideEdit)).Methods("GET")
	r.HandleFunc("/admin/config/timeOverride/edit/{id}", util.Handler(admin.TimeOverrideSubmit)).Methods("POST")
	r.HandleFunc("/admin/config/timeOverride/edit/", util.Handler(admin.TimeOverrideSubmit)).Methods("POST")
	r.HandleFunc("/admin/config/timeOverride/delete", util.Handler(admin.TimeOverrideDelete)).Methods("POST")

	//Auxiliary admin routes
	r.HandleFunc("/admin/bootstrap", util.Handler(admin.Bootstrap))
	r.HandleFunc("/admin/migrate", util.Handler(admin.Migrate))

	//Client API routes
	//For download and downloadComplete actions of presentations
	//"active" can be used instead of key to select the currently active presentation
	r.HandleFunc("/api/presentation/{key}/download", util.Handler(api.Download))
	//This rute is the same as the one above, but is used in UI so that
	//file downloaded by users will not be called "download"
	r.HandleFunc("/api/presentation/{key}/download/{filename}", util.Handler(api.Download))
	r.HandleFunc("/api/presentation/{key}/downloadComplete", util.Handler(api.DownloadFinish))
	r.HandleFunc("/api/presentation/{key}/description", util.Handler(api.GetDescription)).Methods("GET")
	r.HandleFunc("/api/presentation/{key}/description", util.Handler(api.UpdateDescription)).Methods("POST")
	r.HandleFunc("/api/presentation/{key}/name", util.Handler(api.GetName)).Methods("GET")
	r.HandleFunc("/api/presentation/{key}/name", util.Handler(api.UpdateName)).Methods("POST")
	r.HandleFunc("/api/presentation/{key}/schedule", util.Handler(api.ScheduleActivation)).Methods("POST")
	r.HandleFunc("/api/presentation/activate", util.Handler(api.ActivateScheduled)).Methods("GET")
	r.HandleFunc("/api/config/download", util.Handler(api.GetConfig)).Methods("GET")
	r.HandleFunc("/api/config/downloadComplete", util.Handler(api.GotConfig)).Methods("GET")
	r.HandleFunc("/api/update", util.Handler(api.Update))

	r.StrictSlash(true)
	http.Handle("/", r)
}
