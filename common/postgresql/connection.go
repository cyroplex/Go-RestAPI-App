package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
)

func GetConnectionPool(context context.Context, config Config) *pgxpool.Pool {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable pool_max_conns=%s pool_max_conn_idle_time=%s",
		config.Host,
		config.Port,
		config.UserName,
		config.Password,
		config.DbName,
		config.MaxConnections,
		config.MaxConnectionsIdleTime)

	connConfig, parseConfigErr := pgxpool.ParseConfig(connString)
	if parseConfigErr != nil {
		panic(parseConfigErr)
	}

	conn, err := pgxpool.NewWithConfig(context, connConfig)
	if err != nil {
		log.Error("Unable to connect to database %v\n", err)
		panic(err)
	}

	return conn
}
