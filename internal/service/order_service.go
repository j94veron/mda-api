package service

import (
	"context"
	"fmt"
	"go-api/internal/platform"
)

// OrderService - Process and send orders to MDA.
type OrderService interface {
	ProcessOrders(ctx context.Context) error
	SendOrder(ctx context.Context, request OrderRequest) error
}

// OrderServiceImp - Impl. from OrderService.
type orderServiceImp struct {
	mdaClient   platform.MDAClient
	mysqlClient string
}

// ProcessOrders -
func (o *orderServiceImp) ProcessOrders(ctx context.Context) error {
	//
	return nil
}

// SendOrder -
func (o *orderServiceImp) SendOrder(ctx context.Context, request OrderRequest) error {
	// request for mda
	req := platform.OrderRequest{
		OrderID: request.OrderID,
	}

	err := o.mdaClient.SendOrder(ctx, req)
	if err != nil {
		return err
	}

	fmt.Printf("Created Order: %+v", req)
	return nil
}

func NewOrderService(mdaClient platform.MDAClient, mysqlClient string) OrderService {
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
}
