package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func API(postgresClient *pgxpool.Pool, log *zerolog.Logger) chi.Router {
	r := chi.NewRouter()

	r.Route("/account", func(r chi.Router) {
		account(postgresClient, log, r)
	})

	//r.Route("/endpoint", func(r chi.Router) {
	//	endpoint(postgresClient, log, r)
	//})

	return r
}
