package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/config"
	"github.com/Dakak-Takto/yandex-practicum-go-diploma/internal/entity"
)

func (s *service) ProcessNewOrders() error {

	orders, err := s.storage.SelectNewOrders()
	if err != nil {
		s.log.Errorf("Error get new orders: %s", err)
		s.log.Errorf("sleep 5 second")
		time.Sleep(time.Second * 5)
	}
	for _, order := range orders {

		url := fmt.Sprintf("%s/%s/%d", config.AccrualSystemAddress(), "/api/orders/", order.Number)

		response, err := http.Get(url)
		if err != nil {
			s.log.Errorf("Error get accrual info: %s", err)
			continue
		}
		if response.StatusCode == http.StatusTooManyRequests {
			var sleepTime int = 60

			sleepTime, err := strconv.Atoi(response.Header.Get("Retry-After"))
			if err != nil {
				s.log.Errorf("error parse retry-after time: %s", err)
			}

			s.log.Debug("accrual client sleep: %d second", sleepTime)
			time.Sleep(time.Second * time.Duration(sleepTime))
		}

		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			s.log.Errorf("error read body: %s", err)
			continue
		}

		var r struct {
			Order   int                `json:"order"`
			Status  entity.OrderStatus `json:"status"`
			Accrual float64            `json:"accrual"`
		}

		err = json.Unmarshal(body, &r)
		if err != nil {
			s.log.Errorf("error unmarshal accrual response json: %s", err)
			continue
		}

		order.Accrual = r.Accrual
		order.Status = r.Status

		err = s.UpdateOrder(order)
		if err != nil {
			s.log.Errorf("error update order: %s", err)
			continue
		}

		user, err := s.storage.GetUserByID(order.UserID)
		if err != nil {
			s.log.Errorf("error get user: %s", err)
			continue
		}
		user.Balance = user.Balance + r.Accrual
		err = s.UpdateUser(user)
		if err != nil {
			s.log.Errorf("error update user: %s", err)
			continue
		}
	}
	return nil
}
