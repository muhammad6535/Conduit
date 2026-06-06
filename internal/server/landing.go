package server

import (
	_ "embed"
	"net/http"
)

//go:embed landing.html
var landingHTML string

func LandingHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(landingHTML))
}
