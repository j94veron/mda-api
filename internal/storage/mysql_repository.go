package storage

import (
	"database/sql"
	"fmt"
	"go-api/config"
)

func Connect() (*sql.DB, error) {
	config.LoadConfig()

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.AppConfig.MySQL.DBUser,
		config.AppConfig.MySQL.DBPass,
		config.AppConfig.MySQL.DBServer,
		config.AppConfig.MySQL.DBPort,
		config.AppConfig.MySQL.DBName,
	)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return db, nil
}
