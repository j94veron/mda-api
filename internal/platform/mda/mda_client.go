package mda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-api/config"
	"io/ioutil"
	"net/http"
)

type MDAClient interface {
	SendOrder(ctx context.Context, request OrderRequest) (OrderResponse, error)
}

type mdaClientImpl struct {
	endpoint string
}

func (c *mdaClientImpl) SendOrder(ctx context.Context, request OrderRequest) (OrderResponse, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return OrderResponse{Message: err.Error()}, err
	}

	res, err := http.Post(c.endpoint, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return OrderResponse{Message: err.Error()}, err
	}

	defer res.Body.Close()

	// Check the HTTP status code
	if res.StatusCode != http.StatusOK {
		fmt.Println("HTTP request failed with status:", res.Status)
		responseError, err := createResponse(res)
		if err != nil {
			return OrderResponse{Message: err.Error()}, err
		}

		return responseError, nil
	}

	response, err := createResponse(res)
	if res != nil {
		return OrderResponse{}, err
	}

	return response, nil
}

type OrderResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type OrderRequest struct {
	OrderID   string         `json:"token"`
	OrderDev  int            `json:"devolucion"`
	OrderNro  string         `json:"nro_aremito"`
	OrderDate string         `json:"fecha_emision"`
	OrderCuit string         `json:"cuit_destino"`
	Product   []OrderProduct `json:"productos"`
}

type OrderProduct struct {
	Register string `json:"nro_registro"`
	Quantity int    `json:"cantidad"`
	Content  int    `json:"contenido"`
}

func NewMDAClient(config config.MDAConfig) MDAClient {
	endpoint := fmt.Sprintf("%s%s", config.URLBase, config.Path)
	return &mdaClientImpl{
		endpoint: endpoint,
	}
}

func createResponse(res *http.Response) (OrderResponse, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("error sending order to processw %w", err)
	}

	// Unmarshal the JSON response into the struct
	var response OrderResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return OrderResponse{}, fmt.Errorf("error sending order to process %w", err)
	}

	return response, nil
}
