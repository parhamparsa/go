package v1

import (
	"app/config"
	"app/internal/application"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

type RestApiV1 struct {
	config *config.Config
	app    application.App
}

func NewRestApiV1(config *config.Config, app application.App) *RestApiV1 {
	return &RestApiV1{config: config, app: app}
}

func (restApiV1 *RestApiV1) initHandlers() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", restApiV1.healthCheckHandler)

		r.Route("/user", func(r chi.Router) {
			r.Post("/create", restApiV1.createUserHandler)
			r.Get("/find", restApiV1.findUserHandler)
		})

		r.Route("/account", func(r chi.Router) {
			r.Post("/create", restApiV1.createAccountHandler)
			r.Post("/update", restApiV1.updateAccountHandler)
			r.Get("/find", restApiV1.findAccountHandler)
		})

		r.Route("/transfer", func(r chi.Router) {
			r.Post("/money", restApiV1.transferHandler)
		})
	})

	return r
}

func (restApiV1 *RestApiV1) Run() error {
	srv := &http.Server{
		Addr:         restApiV1.config.Addr,
		Handler:      restApiV1.initHandlers(),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Minute,
	}
	log.Printf("listening on %s", restApiV1.config.Addr)
	return srv.ListenAndServe()
}
