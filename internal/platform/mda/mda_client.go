package mda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-api/config"
	"net/http"
)

type MDAClient interface {
	SendOrder(ctx context.Context, request OrderRequest) error
}

type mdaClientImpl struct {
	endpoint string
}

func (c *mdaClientImpl) SendOrder(ctx context.Context, request OrderRequest) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	res, err := http.Post(c.endpoint, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	fmt.Println(res)
	return nil
}

type OrderRequest struct {
	OrderID   string `json:"token"`
	OrderDev  string `json:"devolucion"`
	OrderNro  string `json:"nro_aremito"`
	OrderDate string `json:"fecha_emision"`
	OrderCuit string `json:"cuit_destino"`
}

func NewMDAClient(config config.MDAConfig) MDAClient {
	endpoint := fmt.Sprintf("%s%s", config.URLBase, config.Path)
	return &mdaClientImpl{
		endpoint: endpoint,
	}
}
