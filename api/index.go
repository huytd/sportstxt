package handler

import (
	"net/http"

	"huy.rocks/sports/sports"
)

var handler = sports.NewHandler()

func Handler(w http.ResponseWriter, r *http.Request) {
	handler.ServeHTTP(w, r)
}
