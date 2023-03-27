package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pmoieni/auth-server/internal/service"
	"github.com/pmoieni/auth-server/internal/user"
)

type Service struct {
	mux  service.Service
	repo user.Repo
}

func New(ctx context.Context, r user.Repo) *Service {
	s := Service{
		mux:  service.New(),
		repo: r,
	}
	s.routes()
	return &s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Service) routes() {
	s.mux.Post("/v0/auth/login", s.handleLogin())
	s.mux.Post("/v0/auth/registry", s.handleRegistration())
	s.mux.Get("/v0/auth/refresh", s.handleRefreshToken())
}

func (s *Service) handleLogin() http.HandlerFunc {
	type req struct {
	}

	type res struct {
	}

	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *Service) handleRegistration() http.HandlerFunc {
	type req struct {
	}

	type res struct {
	}

	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *Service) handleRefreshToken() http.HandlerFunc {
	type req struct {
	}

	type res struct {
	}

	return func(w http.ResponseWriter, r *http.Request) {

	}
}

type AuthError struct {
	StatusCode int   `json:"status"`
	Err        error `json:"err"`
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("status %d: err %v", e.StatusCode, e.Err)
}
