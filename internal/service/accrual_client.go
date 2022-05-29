package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/config"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
)

func (s *service) ProcessNewOrders() error {
	orders, err := s.storage.SelectNewOrders()
	if err != nil {
		s.log.Errorf("Error get new orders: %s", err)
		return err
	}

	if len(orders) == 0 {
		s.log.Debug("no new orders to process")
	} else {
		s.log.Debugf("%d new orders. Process...", len(orders))
	}

	for _, order := range orders {

		url := fmt.Sprintf("%s/api/orders/%s", config.AccrualSystemAddress(), order.Number)

		s.log.Debugf("get %s", url)

		response, err := http.Get(url)
		if err != nil {
			s.log.Errorf("Error get accrual info: %s", err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusTooManyRequests {
			sleepTime := 60

			sleepTime, err := strconv.Atoi(response.Header.Get("Retry-After"))
			if err != nil {
				s.log.Errorf("error parse retry-after time: %s", err)
			}

			s.log.Debugf("accrual client sleep: %d second", sleepTime)
			time.Sleep(time.Second * time.Duration(sleepTime))
		}

		if response.StatusCode != http.StatusOK {
			s.log.Errorf("response status code: %d", response.StatusCode)
			continue
		}

		contentType := response.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			s.log.Errorf("content-type not json: %s", contentType)
			continue
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			s.log.Errorf("error read body: %s", err)
			continue
		}

		var r struct {
			Order   string             `json:"order"`
			Status  entity.OrderStatus `json:"status"`
			Accrual float64            `json:"accrual"`
		}

		err = json.Unmarshal(body, &r)
		if err != nil {
			s.log.Errorf("error unmarshal accrual response json: %s. Body: %s", err, body)
			continue
		}

		user, err := s.storage.GetUserByID(order.UserID)
		if err != nil {
			s.log.Errorf("error get user: %s", err)
			continue
		}

		user.Balance = user.Balance + r.Accrual

		s.log.Debugf("New user balance %f", user.Balance)

		err = s.UpdateUser(user)
		if err != nil {
			s.log.Errorf("error update user: %s", err)
			continue
		}

		order.Accrual = r.Accrual
		order.Status = r.Status

		err = s.UpdateOrder(order)
		if err != nil {
			s.log.Errorf("error update order: %s", err)
			continue
		}

		s.log.Debug("Order updated")

	}
	return nil
}
