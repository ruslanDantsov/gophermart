package client

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
	"github.com/ruslanDantsov/gophermart/internal/dto/view"
	"github.com/ruslanDantsov/gophermart/internal/errs"
	"net/http"
)

const GetOrderStatusURL = "%s/api/orders/%s"

type OrderStatusClient struct {
	httpClient *resty.Client
	baseURL    string
}

func NewOrderStatusClient(baseURL string) *OrderStatusClient {
	httpClient := resty.New()

	return &OrderStatusClient{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *OrderStatusClient) GetAccrualData(ctx context.Context, orderID string) (*view.AccrualResponse, error) {
	url := fmt.Sprintf(GetOrderStatusURL, c.baseURL, orderID)

	resp, err := c.httpClient.R().
		SetHeader("Content-Type", "text/plain").
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		errMessage := fmt.Sprintf("bad response from Accrual service %s: %v", orderID, resp.StatusCode())
		return nil, errs.New(errs.OrderStatusClient, errMessage, nil)
	}

	var responseBody view.AccrualResponse

	if err := easyjson.Unmarshal(resp.Body(), &responseBody); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &responseBody, err
}
