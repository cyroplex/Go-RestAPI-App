package app

import "product-app/common/postgresql"

type ConfigrationManager struct {
	PostgreSqlConfig postgresql.Config
}

func NewConfigrationManager() *ConfigrationManager {
	postgreSqlConfig := getPostgreSqlConfig()
	return &ConfigrationManager{
		PostgreSqlConfig: postgreSqlConfig,
	}
}

func getPostgreSqlConfig() postgresql.Config {
	return postgresql.Config{
		Host:                   "localhost",
		Port:                   "6432",
		DbName:                 "productapp",
		UserName:               "postgres",
		Password:               "postgres",
		MaxConnections:         "10",
		MaxConnectionsIdleTime: "30s",
	}
}
