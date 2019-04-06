package http

import (
	"net/http"

	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

func SetupHealthRoutes() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(g.Version))
	})
}
