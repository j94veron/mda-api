package mda

import "context"

type MDAClient interface {
	SendOrder(ctx context.Context, request OrderRequest) error
}

type mdaClientImpl struct {
	endpoint string
}

func (o *mdaClientImpl) SendOrder(ctx context.Context, request OrderRequest) error {

	//body := `"order_id": 50`
	//err := endpoint.post(body)

	return nil
}

type OrderRequest struct {
	OrderID   string `json:"token"`
	OrderDev  string `json:"devolucion"`
	OrderNro  string `json:"nro_aremito"`
	OrderDate string `json:"fecha_emision"`
	OrderCuit string `json:"cuit_destino"`
}

func NewMDAClient(url string) MDAClient {
	return &mdaClientImpl{
		endpoint: url,
	}
}
