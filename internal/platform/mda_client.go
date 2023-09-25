package platform

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
	OrderID string `json:"order_id"`
}

func NewMDAClient(url string) MDAClient {
	return &mdaClientImpl{
		endpoint: "http://ministerio.com.ar",
	}
}
