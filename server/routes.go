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
	r.Handle("/", util.Handler(public.Index))
	r.Handle("/presentation", util.Handler(public.Presentations))

	//Administrative routes
	r.Handle("/admin", util.Handler(admin.Admin))
	r.Handle("/admin/presentation/upload", util.Handler(admin.Upload)).Methods("GET")
	r.Handle("/admin/presentation/upload", util.Handler(admin.UploadHandler)).Methods("POST")
	r.Handle("/admin/presentation/archive/{page}", util.Handler(admin.Archive))
	r.Handle("/admin/presentation/activate", util.Handler(admin.Activate)).Methods("POST")
	r.Handle("/admin/presentation/deactivate", util.Handler(admin.Deactivate)).Methods("POST")
	r.Handle("/admin/presentation/delete", util.Handler(admin.Delete)).Methods("POST")
	r.Handle("/admin/presentation/{id}", util.Handler(admin.Presentation))
	r.Handle("/admin/config", util.Handler(admin.ShowConfig)).Methods("GET")
	r.Handle("/admin/config", util.Handler(admin.SetConfig)).Methods("POST")
	r.Handle("/admin/config/timeOverride", util.Handler(admin.TimeOverride))
	r.Handle("/admin/config/timeOverride/edit/{id}", util.Handler(admin.TimeOverrideEdit)).Methods("GET")
	r.Handle("/admin/config/timeOverride/edit", util.Handler(admin.TimeOverrideEdit)).Methods("GET")
	r.Handle("/admin/config/timeOverride/edit/{id}", util.Handler(admin.TimeOverrideSubmit)).Methods("POST")
	r.Handle("/admin/config/timeOverride/edit/", util.Handler(admin.TimeOverrideSubmit)).Methods("POST")
	r.Handle("/admin/config/timeOverride/delete", util.Handler(admin.TimeOverrideDelete)).Methods("POST")

	//Auxiliary admin routes
	r.Handle("/admin/bootstrap", util.Handler(admin.Bootstrap))
	r.Handle("/admin/migrate", util.Handler(admin.Migrate))

	//Client API routes
	//For download and downloadComplete actions of presentations
	//"active" can be used instead of key to select the currently active presentation
	r.Handle("/api/presentation/{key}/download", util.Handler(api.Download))
	//This rute is the same as the one above, but is used in UI so that
	//file downloaded by users will not be called "download"
	r.Handle("/api/presentation/{key}/download/{filename}", util.Handler(api.Download))
	r.Handle("/api/presentation/{key}/downloadComplete", util.Handler(api.DownloadFinish))
	r.Handle("/api/presentation/{key}/deactivated", util.Handler(api.Deactivated))
	r.Handle("/api/presentation/{key}/description", util.Handler(api.GetDescription)).Methods("GET")
	r.Handle("/api/presentation/{key}/description", util.Handler(api.UpdateDescription)).Methods("POST")
	r.Handle("/api/presentation/{key}/name", util.Handler(api.GetName)).Methods("GET")
	r.Handle("/api/presentation/{key}/name", util.Handler(api.UpdateName)).Methods("POST")
	r.Handle("/api/presentation/{key}/schedule", util.Handler(api.ScheduleActivation)).Methods("POST")
	r.Handle("/api/activation/{key}/delete", util.Handler(api.DeleteActivation)).Methods("POST")
	r.Handle("/api/presentation/activate", util.Handler(api.ActivateScheduled)).Methods("GET")
	r.Handle("/api/config/download", util.Handler(api.GetConfig)).Methods("GET")
	r.Handle("/api/config/downloadComplete", util.Handler(api.GotConfig)).Methods("GET")
	r.Handle("/api/status", util.Handler(api.Status))

	r.StrictSlash(true)
	http.Handle("/", r)
}
