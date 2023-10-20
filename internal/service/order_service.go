package service

import (
	"context"
	"go-api/internal/platform/mda"
	"go-api/internal/storage"
)

// OrderService - Process and send orders to MDA.
type OrderService interface {
	ProcessOrders(ctx context.Context) error
	SendOrder(ctx context.Context, request OrderRequest) (OrderResponse, error)
}

// OrderServiceImp - Impl. from OrderService.
type orderServiceImp struct {
	mdaClient   mda.MDAClient
	mysqlClient storage.MysqlRepository
}

// ProcessOrders -
func (o *orderServiceImp) ProcessOrders(ctx context.Context) error {
	//
	return nil
}

// SendOrder -
func (o *orderServiceImp) SendOrder(ctx context.Context, request OrderRequest) (OrderResponse, error) {
	// request for mda
	req := mda.OrderRequest{
		OrderID:   request.OrderID,
		OrderDev:  request.OrderDev,
		OrderNro:  request.OrderNro,
		OrderDate: request.OrderDate,
		OrderCuit: request.OrderCuit,
	}

	err := o.mdaClient.SendOrder(ctx, req)
	if err != nil {
		return OrderResponse{}, err
	}

	res := OrderResponse{
		Token:     request.OrderID,
		Productos: getProductMock(request.Products),
	}

	return res, nil
}

func getProductMock(productsReq []ProductRequest) []ProductResponse {
	var productsRes []ProductResponse
	for _, value := range productsReq {
		productsRes = append(productsRes, ProductResponse{
			NroRegistro: value.RegisterNumber,
			Cantidad:    value.QuantityNumber,
			Contenido:   value.Content,
		})
	}
	return productsRes
}

func NewOrderService(mdaClient mda.MDAClient, mysqlClient storage.MysqlRepository) OrderService {
	return &orderServiceImp{
		mdaClient:   mdaClient,
		mysqlClient: mysqlClient,
	}
}

type OrderRequest struct {
	OrderID    string
	OrderDev   string
	OrderNro   string
	OrderIdent string
	OrderDate  string
	OrderCuit  string
	Products   []ProductRequest
}

type ProductRequest struct {
	RegisterNumber string
	QuantityNumber int
	Content        int
}

type OrderResponse struct {
	Token        string            `json:"token"`
	Devolucion   string            `json:"devolucion"`
	NroRemito    string            `json:"nro_aremito"`
	CuitDestino  string            `json:"cuit_destino"`
	FechaEmision string            `json:"fecha_emision"`
	Productos    []ProductResponse `json:"productos"`
}

type ProductResponse struct {
	NroRegistro string `json:"nro_registro"`
	Cantidad    int    `json:"cantidad"`
	Contenido   int    `json:"contenido"`
}
