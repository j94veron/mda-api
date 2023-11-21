package handler

import (
	"encoding/json"
	"fmt"
	"go-api/internal/service"
	"net/http"
)

func NewSendOrderHandler(orderService service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqDTO order
		err := json.NewDecoder(r.Body).Decode(&reqDTO)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		orderReq := service.OrderRequest{
			OrderID:   reqDTO.Token,
			OrderDev:  reqDTO.Devolucion,
			OrderNro:  reqDTO.NroRemito,
			OrderDate: reqDTO.FechaEmision,
			OrderCuit: reqDTO.CuitDestino,
			Products:  getProducts(reqDTO.Productos),
		}

		res, err := orderService.SendOrder(r.Context(), orderReq)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonResponse, err := json.Marshal(res)
		if err != nil {
			http.Error(w, fmt.Sprintf("error sending orders %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

func NewPendingOrdersHandler(orderService service.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := orderService.ProcessOrders(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting pending orders %s", err.Error()), http.StatusInternalServerError)
			return
		}

		jsonResponse, err := json.Marshal(res)
		if err != nil {
			http.Error(w, "error parsing response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

func getProducts(productos []product) []service.ProductRequest {
	var products []service.ProductRequest
	for _, value := range productos {
		products = append(products, service.ProductRequest{
			RegisterNumber: value.NroRegistro,
			QuantityNumber: value.Cantidad,
			Content:        value.Contenido,
		})
	}
	return products
}

type order struct {
	Token        string    `json:"token" validate:"required"`
	Devolucion   int       `json:"devolucion"`
	NroRemito    string    `json:"nro_aremito"`
	CuitDestino  string    `json:"cuit_destino"`
	FechaEmision string    `json:"fecha_emision"`
	Productos    []product `json:"productos"`
}

type product struct {
	NroRegistro string `json:"nro_registro"`
	Cantidad    int    `json:"cantidad"`
	Contenido   int    `json:"contenido"`
}
