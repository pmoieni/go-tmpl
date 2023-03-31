package user

import (
	"context"
	"fmt"
	"net/http"

	ory "github.com/ory/client-go"
	"github.com/pmoieni/kratos-test/internal/service"
)

type Service struct {
	mux service.Service
	ory *ory.APIClient
}

func New(ctx context.Context, ory *ory.APIClient) *Service {
	s := Service{
		mux: service.New(),
		ory: ory,
	}
	s.routes()
	return &s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Service) routes() {
	s.mux.Get("/", s.ValidateSession(s.handleHelloWorld()))
}

func (s *Service) handleHelloWorld() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.mux.Respond(w, "Hello World!", http.StatusOK)
	}
}

type AuthError struct {
	StatusCode int   `json:"status"`
	Err        error `json:"err"`
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("status %d: err %v", e.StatusCode, e.Err)
}
