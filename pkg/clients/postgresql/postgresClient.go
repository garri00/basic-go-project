package postgresql

import (
	"context"
	"fmt"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"basic-go-project/src/config"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

func NewClient(ctx context.Context, configs config.PostgresDBConf, l *zerolog.Logger) (db *pgxpool.Pool, err error) {
	//TODO: change to sslmode=requier
	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		configs.Username,
		configs.Password,
		configs.Host,
		configs.Port,
		configs.Database,
		configs.SSLMode)

	dbPool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		err = fmt.Errorf("pgxpool.New failed: %w", err)
		l.Err(err).Send()

		return nil, err
	}

	if err := dbPool.Ping(ctx); err != nil {
		err = fmt.Errorf("dbPool.Ping failed: %w", err)
		l.Err(err).Send()

		return nil, err
	}

	log.Info().Msg("successfully connected to PostgresDB")

	return dbPool, nil
}

const (
	defaultMaxConns          = int32(4)
	defaultMinConns          = int32(0)
	defaultMaxConnLifetime   = time.Hour
	defaultMaxConnIdleTime   = time.Minute * 30
	defaultHealthCheckPeriod = time.Minute
	defaultConnectTimeout    = time.Second * 5
)

// Config TODO: for implementing pgxpool.NewWithConfig
func Config(connectionURL string, l zerolog.Logger) *pgxpool.Config {
	dbConfig, err := pgxpool.ParseConfig(connectionURL)
	if err != nil {
		l.Log().Err(err).Msgf("Failed to create a config")

		return nil
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		l.Log().Msg("Before acquiring the connection pool to the database!")

		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		l.Log().Msg("After releasing the connection pool to the database!")

		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		l.Log().Msg("Closed the connection pool to the database!")
	}

	return dbConfig
}
