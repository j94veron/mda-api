package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-api/config"
	"go-api/internal/handler"
	"go-api/internal/platform/mda"
	"go-api/internal/service"
	"go-api/internal/storage"
	"net/http"
)

func Run() error {
	var appConfig config.Configuration
	err := config.LoadConfig(&appConfig)
	if err != nil {
		fmt.Println(fmt.Printf("error loading config %v", err))
		return err
	}

	mysqlClient := storage.NewMysqlRepository(appConfig.MySQLConfig)
	mdaClient := mda.NewMDAClient(appConfig.MDAConfig)
	orderService := service.NewOrderService(mdaClient, mysqlClient)

	r := mux.NewRouter()
	r.HandleFunc("/order/send", handler.NewSendOrderHandler(orderService)).Methods("POST")
	r.HandleFunc("/order", handler.NewPendingOrdersHandler(orderService)).Methods("GET")
	http.Handle("/", r)

	return nil
}
