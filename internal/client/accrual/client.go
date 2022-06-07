package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/logger"
	_service "github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/service"
)

type Accrual interface {
	Run(context.Context)
}

type accrual struct {
	service _service.Service
	baseURL string
}

var log = logger.GetLoggerInstance()

func New(service _service.Service, accrualSystemAddress string) Accrual {
	return &accrual{
		service: service,
		baseURL: accrualSystemAddress,
	}
}

func (a *accrual) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Second)
			a.processOrders()
		}
	}
}

func (a *accrual) processOrders() {
	orders, err := a.service.GetNewOrders(context.Background())
	if err != nil {
		log.Errorf("Error get new orders: %s", err)
		return
	}

	if len(orders) == 0 {
		log.Debug("no new orders to process")
	} else {
		log.Debugf("%d new orders. Process...", len(orders))
	}

	for _, order := range orders {
		result, err := a.getOrderInfo(order)
		if err != nil {
			log.Error(err)
			continue
		}

		order.Accrual = result.Accrual

		switch result.Status {
		case entity.OrderStatusNEW.String():
			order.Status = entity.OrderStatusNEW
		case entity.OrderStatusREGISTERED.String():
			order.Status = entity.OrderStatusREGISTERED
		case entity.OrderStatusINVALID.String():
			order.Status = entity.OrderStatusPROCESSED
		case entity.OrderStatusPROCESSING.String():
			order.Status = entity.OrderStatusPROCESSING
		case entity.OrderStatusPROCESSED.String():
			order.Status = entity.OrderStatusPROCESSED
		default:
			order.Status = entity.OrderStatusUNKNOWN
		}

		if err := a.service.UserBalanceChange(context.Background(), order.UserID, +order.Accrual); err != nil {
			log.Errorf("error update user balance: userID: %d, delta %f", order.UserID, order.Accrual)
		}

		log.Debugf("Change user balance: userID %d, delta %f", order.UserID, order.Accrual)

		err = a.service.UpdateOrder(context.Background(), order)
		if err != nil {
			log.Errorf("error update order: %s", err)
			return
		}

		log.Debug("Order updated")
	}
}

func (a *accrual) getOrderInfo(order *entity.Order) (*orderAccrualResponseDTO, error) {

	url := fmt.Sprintf("%s/api/orders/%s", a.baseURL, order.Number)

	log.Debugf("get %s", url)

	response, err := http.Get(url)
	if err != nil {
		log.Errorf("Error get accrual info: %s", err)
		return nil, err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Errorf("error close body: %s", err)
		}
	}()

	if response.StatusCode == http.StatusTooManyRequests {
		sleepTime := 60

		sleepTime, err := strconv.Atoi(response.Header.Get("Retry-After"))
		if err != nil {
			log.Errorf("error parse retry-after time: %s", err)
			return nil, err
		}

		log.Debugf("accrual client sleep: %d second", sleepTime)
		time.Sleep(time.Second * time.Duration(sleepTime))
	}

	if response.StatusCode != http.StatusOK {
		log.Errorf("response status code: %d", response.StatusCode)
		return nil, fmt.Errorf("status code: %d", response.StatusCode)
	}

	contentType := response.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		log.Errorf("content-type not json: %s", contentType)
		return nil, err
	}

	var r orderAccrualResponseDTO

	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		log.Errorf("error unmarshal accrual response json: %s", err)
		return nil, err
	}

	return &r, nil
}
