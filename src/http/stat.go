package http

import (
	"../g"
	"net/http"
)

func configStatRoutes() {

	http.HandleFunc("/stat/general", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, g.GetGStat())
	})

	http.HandleFunc("/stat/apps", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, g.GetAStat())
	})

	http.HandleFunc("/stat/crons", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, g.GetCStat())
	})

	http.HandleFunc("/stat/tasks", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, g.GetTStat())
	})

	http.HandleFunc("/stat/errors", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, g.GetEStat())
	})

}
