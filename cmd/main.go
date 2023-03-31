package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	ory "github.com/ory/client-go"
	userService "github.com/pmoieni/kratos-test/service/user"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
)

type appConfig struct {
	ServerPort string `json:"serverPort"`
	Dev        bool   `json:"dev"`
}

func main() {

}

func serve(cfg *appConfig) error {
	sCtx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	// ? should this defined within the instantiation of a new service
	c := cors.Options{
		AllowedOrigins:   []string{"*"}, // ? band-aid, needs to change to a flag
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposedHeaders:   []string{"Location"},
		Debug:            cfg.Dev,
	}

	mux := chi.NewMux().Route("/v0", func(r chi.Router) {
		r.Handle("/users", newUserService(sCtx))
	})

	/* START SERVICES BLOCK */
	srv := http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: cors.New(c).Handler(mux),
		// max time to read request from the client
		ReadTimeout: 10 * time.Second,
		// max time to write response to the client
		WriteTimeout: 10 * time.Second,
		// max time for connections using TCP Keep-Alive
		IdleTimeout: 120 * time.Second,
		BaseContext: func(_ net.Listener) context.Context { return sCtx },
		ErrorLog:    log.Default(),
	}

	g, gCtx := errgroup.WithContext(sCtx)

	g.Go(func() error {
		// Run the server
		srv.ErrorLog.Printf("App server starting on %s", srv.Addr)
		return srv.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()
		return srv.Shutdown(context.Background())
	})

	return g.Wait()
}

func newUserService(ctx context.Context) *userService.Service {
	proxyPort := os.Getenv("PROXY_PORT")
	if proxyPort == "" {
		proxyPort = "5173"
	}

	// register a new Ory client with the URL set to the Ory CLI Proxy
	// we can also read the URL from the env or a config file
	c := ory.NewConfiguration()
	c.Servers = ory.ServerConfigurations{{URL: fmt.Sprintf("http://localhost:%s/.ory", proxyPort)}}

	service := userService.New(ctx, ory.NewAPIClient(c))
	return service
}
