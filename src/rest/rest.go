package rest

import (
	"net/http"

	"github.com/gorilla/pat"

	"github.com/urfave/negroni"
)

func RunAPI(address string) error {
	mux := pat.New()
	h, _ := NewHandler()
	nh := negroni.Classic()

	nh.UseHandler(mux)

	mux.Get("/", h.GetMainPage)
	mux.Post("/signup", h.SignUp)
	mux.Get("/signin", h.SignIn)
	mux.Get("/signout", h.SignOut)
	mux.Get("/standard", h.GetStandard)
	return http.ListenAndServe(address, nh)
}
