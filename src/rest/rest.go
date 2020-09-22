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
	mux.Post("/signin", h.SignIn)
	mux.Post("/signout", h.SignOut)
	mux.Get("/profile", h.GetProfile)
	mux.Get("/standard", h.GetStandard)
	return http.ListenAndServe(address, nh)
}
