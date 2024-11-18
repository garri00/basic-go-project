package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"basic-go-project/api/adapters/db/postgres"
	"basic-go-project/api/handlers"
	"basic-go-project/api/usecases"
)

func account(postgresClient *pgxpool.Pool, log *zerolog.Logger, r chi.Router) {
	accountCollection := postgres.NewAccountStoragePG(postgresClient, log)
	accountUseCase := usecases.NewAccountUseCase(accountCollection, log)
	accountHandler := handlers.NewAccountHandler(accountUseCase, log)

	r.Post("/", accountHandler.Create)
	r.Get("/", accountHandler.GetAll)
	r.Get("/{id}", accountHandler.GetByID)
	r.Patch("/{id}", accountHandler.Update)
	r.Delete("/{id}", accountHandler.Delete)
}
