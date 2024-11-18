package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"basic-go-project/api"
	"basic-go-project/pkg/clients/postgresql"
	"basic-go-project/pkg/logger"
	"basic-go-project/pkg/utils"
	"basic-go-project/src/config"
	"basic-go-project/src/entities"
)

func main() {
	//Init project and env configs
	configs, err := config.GetConfig()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to load .env")
	}

	//Setup main logger with level
	if err := logger.SetLogger(configs); err != nil {
		logger.Log.Fatal().Err(err).Msg("logger.SetMainLogger failed")
	}

	//Setup postgres DB connection
	ctx := context.Background()
	postgresClient, err := postgresql.NewClient(ctx, configs.PostgresConf, &logger.Log)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to create postgres connection")
	}

	//Do all migrations
	if err := postgresql.MigrateUp(configs.PostgresConf, logger.Log); err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to start postgres migration")
	}

	//Attach routes
	routes := api.NewRouter(postgresClient, &logger.Log)

	logger.Log.Info().Msg("server started")

	//Set server attributes
	server := http.Server{
		Addr:              fmt.Sprintf("%s:%s", configs.Host, configs.Port),
		Handler:           routes,
		ReadHeaderTimeout: entities.ServiceRequestTimeout * time.Second,
	}

	//TODO:change to dynamically set service version
	serviceVersion := "v0.0.0"

	//Start server
	logger.Log.Info().Msgf("service version=%s", serviceVersion)
	logger.Log.Info().Msgf("host=%s, port=%s, tls=%v, mode=%s", configs.Host, configs.Port, configs.TLS, configs.Mode)

	if configs.TLS {
		//Server with TLS
		serverTLSCert, err := utils.LoadCertificate()
		if err != nil {
			logger.Log.Fatal().Err(fmt.Errorf("utils.LoadCertificate() failed: %w", err)).Msg("Error loading certificate and key file")
		}

		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{*serverTLSCert},
			MinVersion:   tls.VersionTLS12,
		}

		logger.Log.Debug().Msgf("https://%s:%s/health-check", configs.Host, configs.Port)

		if err := server.ListenAndServeTLS("", ""); err != nil {
			logger.Log.Fatal().Err(err).Msg("server crashed")
		}
	} else {
		//Regular http server
		logger.Log.Debug().Msgf("http://%s:%s/health-check", configs.Host, configs.Port)

		if err := server.ListenAndServe(); err != nil {
			logger.Log.Fatal().Err(err).Msg("server crashed")
		}
	}
}
