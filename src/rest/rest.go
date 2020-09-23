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
	mux.Post("/score", h.ReceiveScore)
	mux.Get("/score", h.GetScore)
	return http.ListenAndServe(address, nh)
}
