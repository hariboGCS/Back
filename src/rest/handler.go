package rest

import (
	"fmt"
	"log"
	"net/http"
)

type HandlerInterface interface {
	GetMainPage(w http.ResponseWriter, r *http.Request)
	GetScore(w http.ResponseWriter, r *http.Request)
	GetRankPage(w http.ResponseWriter, r *http.Request)
	AddUser(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	SignOut(w http.ResponseWriter, r *http.Request)
	GetStandard(w http.ResponseWriter, r *http.Request)
	GetNotice(w http.ResponseWriter, r *http.Request)
	GetComplaints(w http.ResponseWriter, r *http.Request)
}
type Handler struct {
	handler HandlerInterface
}

func NewHandler() (*Handler, error) {
	return new(Handler), nil
}
func (h *Handler) GetMainPage(w http.ResponseWriter, r *http.Request) {
	log.Println("Main page....")
	fmt.Fprint(w, "Main page for secure API!!")
}
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Sign In Page")
}
func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {

}
