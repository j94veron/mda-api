package main

import (
	"github.com/gorilla/mux"
	"go-api/internal/handler"
	"go-api/internal/platform"
	"go-api/internal/service"
	"net/http"
)

func Run() error {
	config := Configuration{
		MDAConfig: MDA{
			Path: "http://ministerio.com",
		},
	}
	mdaClient := platform.NewMDAClient(config.MDAConfig.Path)
	orderService := service.NewOrderService(mdaClient, "db")

	r := mux.NewRouter()
	r.HandleFunc("/order/send", handler.NewSendOrderHandler(orderService)).Methods("POST")
	http.Handle("/", r)

	return nil
}
