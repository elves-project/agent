package http

import (
	"github.com/elves-project/agent/src/g"
	"github.com/gy-games-libs/file"
	"github.com/gy-games-libs/seelog"
	"net/http"
	"path/filepath"
	"strings"
)

func configPageRoutes() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			seelog.Debug(filepath.Join(g.Root, "/public", r.URL.Path, "index.html"))
			if !file.IsExist(filepath.Join(g.Root, "/public", r.URL.Path, "index.html")) {
				http.NotFound(w, r)
				return
			}
		}
		http.FileServer(http.Dir(filepath.Join(g.Root, "/public"))).ServeHTTP(w, r)
	})

}
