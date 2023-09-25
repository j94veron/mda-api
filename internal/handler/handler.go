package handler

import (
	"encoding/json"
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
			OrderID: reqDTO.Token,
		}

		err = orderService.SendOrder(r.Context(), orderReq)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

type order struct {
	Token        string    `json:"token" validate:"required"`
	Devolucion   string    `json:"devolucion"`
	NroRemito    string    `json:"nro_aremito"`
	CuitDestino  string    `json:"cuit_destino"`
	FechaEmision string    `json:"fecha_emision"`
	Productos    []product `json:"productos"`
}

type product struct {
	NroRegistro string  `json:"nro_registro"`
	Cantidad    int     `json:"cantidad"`
	Contenido   float64 `json:"contenido"`
}
