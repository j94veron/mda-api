package service

import (
	"context"
	"fmt"
	"go-api/internal/platform/mda"
	"go-api/internal/storage"
)

const (
	_errorStatus    = "error"
	_finishedStatus = "finished"
)

// OrderService - Process and send orders to MDA.
type OrderService interface {
	ProcessOrders(ctx context.Context) ([]mda.OrderResponse, error)
	SendOrder(ctx context.Context, request OrderRequest) (mda.OrderResponse, error)
}

// OrderServiceImp - Impl. from OrderService.
type orderServiceImp struct {
	mdaClient   mda.MDAClient
	mysqlClient storage.MysqlRepository
}

// ProcessOrders -
func (o *orderServiceImp) ProcessOrders(ctx context.Context) ([]mda.OrderResponse, error) {
	orders, err := o.mysqlClient.FindPendingOrders()
	if err != nil {
		return nil, err
	}

	var response []mda.OrderResponse
	var status string
	for _, value := range orders {
		res, err := o.mdaClient.SendOrder(ctx, mda.OrderRequest{
			OrderID:   value.Token,
			OrderDev:  value.Devolucion,
			OrderNro:  value.NroRemito,
			OrderDate: value.FechaEmision,
			OrderCuit: value.CuitDestino,
			Product:   transformProductToOrderProduct(value.Productos),
		})
		if err != nil {
			status = _errorStatus
		} else {
			status = _finishedStatus
		}
		fmt.Printf("Update (id: %d) %s %s", value.NroInterno, status, res.Message)

		// Update con el nro interno, status y message

		err = o.mysqlClient.UpdateOrder(value.NroInterno, status, res.Message)

		response = append(response, res)
		//Update order in DB.
	}

	return response, nil
}

// SendOrder -
func (o *orderServiceImp) SendOrder(ctx context.Context, request OrderRequest) (mda.OrderResponse, error) {
	// request for mda
	req := mda.OrderRequest{
		OrderID:   request.OrderID,
		OrderDev:  request.OrderDev,
		OrderNro:  request.OrderNro,
		OrderDate: request.OrderDate,
		OrderCuit: request.OrderCuit,
		Product:   transformRequestToOrderProduct(request.Products),
	}

	res, err := o.mdaClient.SendOrder(ctx, req)
	if err != nil {
		return mda.OrderResponse{}, err
	}

	return res, nil
}

func transformRequestToOrderProduct(products []ProductRequest) []mda.OrderProduct {
	var response []mda.OrderProduct
	for _, value := range products {
		response = append(response, mda.OrderProduct{
			Register: value.RegisterNumber,
			Quantity: value.QuantityNumber,
			Content:  value.Content,
		})
	}

	return response
}

func transformProductToOrderProduct(products []storage.Product) []mda.OrderProduct {
	var response []mda.OrderProduct
	for _, value := range products {
		response = append(response, mda.OrderProduct{
			Register: value.NroRegistro,
			Quantity: value.Cantidad,
			Content:  value.Contenido,
		})
	}

	return response
}

func NewOrderService(mdaClient mda.MDAClient, mysqlClient storage.MysqlRepository) OrderService {
	return &orderServiceImp{
		mdaClient:   mdaClient,
		mysqlClient: mysqlClient,
	}
}

type OrderRequest struct {
	OrderID    string
	OrderDev   int
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
