package server

import (
  "log"
   "net/http"
   "github.com/pkg/errors"
   "github.com/go-chi/chi/v5"
   "github.com/go-chi/chi/v5/middleware"
   "github.com/jtrw/go-rest"
)

type Server struct {
    Host           string
    Port           string
	PinSize        int
	MaxPinAttempts int
	WebRoot        string
	Version        string
}

func (s Server) Run() error {
    log.Printf("[INFO] Activate rest server")
    log.Printf("[INFO] Port: %s", s.Port)

	if err := http.ListenAndServe(":"+s.Port, s.routes()); err != http.ErrServerClosed {
		return errors.Wrap(err, "server failed")
	}

	return nil
}

func (s Server) routes() chi.Router {
	router := chi.NewRouter()

    router.Use(middleware.Logger)
    router.Use(rest.Ping)
    router.Route("/", func(r chi.Router) {
    })

	return router
}
