package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	m *Model
}

func NewHandler(s Store) *Handler {
	return &Handler{
		m: NewModel(s),
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Auth service is up and running")
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// Find finds the user by id.
// returns status code 204 if user does not exist.
// returns status code 500 on internal errors.
// returns valid user details on success.
func (h *Handler) Find(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.intercept(w, r)
	email := ps.ByName("name")
	var user *User
	if err := h.m.GetUser(email, user); err != nil {
		if user == nil {
			h.respond204(w)
		}
		h.respond500(w, err)
	}

	if err := h.respond200(w, user); err != nil {
		h.respond500(w, err)
	}
}

// Put upserts a user
func (h *Handler) Put(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h.intercept(w, r)

	var user *User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(user); err != nil {
		h.respond500(w, err)
	}
}

func (h *Handler) intercept(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func (h *Handler) respond200(w http.ResponseWriter, res interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		return err
	}

	return nil
}

func (h *Handler) respond500(w http.ResponseWriter, err error) {
	log.Println(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}

func (h *Handler) respond204(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
