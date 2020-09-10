package rest

import (
	"net/http"

	"github.com/gorilla/pat"

	"github.com/urfave/negroni"
)

func RunAPI(address string) error {
	mux := pat.New()
	h, _ := NewHandler()
	mux.Get("/", h.GetMainPage)
	mux.Get("/signin", h.SignIn)

	nh := negroni.Classic()
	nh.UseHandler(mux)
	return http.ListenAndServe(address, nh)
}
