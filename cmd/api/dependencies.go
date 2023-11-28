package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go-api/config"
	"go-api/internal/handler"
	"go-api/internal/platform/mda"
	"go-api/internal/repository"
	"go-api/internal/service"
	"go-api/logger"
	"net/http"
)

func Run() error {
	log := logger.NewLoggerInstance()
	var appConfig config.Configuration
	err := config.LoadConfig(&appConfig)
	if err != nil {
		log.Error(fmt.Sprintf("error loading config %v", err))
		return err
	}

	dbConnection, err := newMySqlConnection(appConfig.MySQLConfig)
	if err != nil {
		log.Error(fmt.Sprintf("error connecting Mysql %v", err))
		return err
	}

	orderRepository := repository.NewOrderRepository(dbConnection)
	mdaClient := mda.NewMDAClient(appConfig.MDAConfig)
	orderService := service.NewOrderService(mdaClient, orderRepository)

	r := mux.NewRouter()
	r.HandleFunc("/order/send", handler.NewSendOrderHandler(orderService)).Methods("POST")
	r.HandleFunc("/order", handler.NewPendingOrdersHandler(orderService)).Methods("GET")
	http.Handle("/", r)

	log.Info(fmt.Sprintf("success run app"))
	return nil
}

func newMySqlConnection(config config.DBConfig) (*sql.DB, error) {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.DBUser,
		config.DBPass,
		config.DBServer,
		config.DBPort,
		config.DBName,
	)

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
