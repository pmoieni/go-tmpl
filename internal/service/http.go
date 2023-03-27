package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Service interface {
	chi.Router

	Log(...any)
	Logf(string, ...any)

	Decode(*http.Request, any) error
	Respond(http.ResponseWriter, *http.Request, any, int)
	Created(http.ResponseWriter, *http.Request, string)
	SetCookie(http.ResponseWriter, *http.Cookie)
}

type service struct {
	chi.Router
}

func (*service) Log(v ...any) { log.Println(v...) }

func (*service) Logf(format string, v ...any) { log.Printf(format, v...) }

func (*service) Decode(r *http.Request, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func (*service) Respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, "Could not encode in json", status)
		}
	}
}

func (s *service) Created(w http.ResponseWriter, r *http.Request, id string) {
	path := r.URL.Path
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	w.Header().Add("Location", "//"+r.Host+path+id)
	s.Respond(w, r, nil, http.StatusCreated)
}

func (*service) SetCookie(w http.ResponseWriter, c *http.Cookie) { http.SetCookie(w, c) }

func New(opt ...Option) Service {
	var s service
	for _, o := range opt {
		o(&s)
	}

	if s.Router == nil {
		s.Router = chi.NewRouter()
	}

	return &s
}

type Option func(*service)

func WithRouter(mux chi.Router) Option {
	return func(s *service) {
		s.Router = mux
	}
}
