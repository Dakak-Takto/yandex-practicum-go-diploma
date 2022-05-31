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
		log.Errorf("Error get new orders: %s", err)
		return err
	}

	if len(orders) == 0 {
		log.Debug("no new orders to process")
	} else {
		log.Debugf("%d new orders. Process...", len(orders))
	}

	for _, order := range orders {
		func() {
			url := fmt.Sprintf("%s/api/orders/%s", config.AccrualSystemAddress(), order.Number)

			log.Debugf("get %s", url)

			response, err := http.Get(url)
			if err != nil {
				log.Errorf("Error get accrual info: %s", err)
				return
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
				}

				log.Debugf("accrual client sleep: %d second", sleepTime)
				time.Sleep(time.Second * time.Duration(sleepTime))
			}

			if response.StatusCode != http.StatusOK {
				log.Errorf("response status code: %d", response.StatusCode)
				return
			}

			contentType := response.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				log.Errorf("content-type not json: %s", contentType)
				return
			}

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Errorf("error read body: %s", err)
				return
			}

			var r struct {
				Order   string             `json:"order"`
				Status  entity.OrderStatus `json:"status"`
				Accrual float64            `json:"accrual"`
			}

			err = json.Unmarshal(body, &r)
			if err != nil {
				log.Errorf("error unmarshal accrual response json: %s. Body: %s", err, body)
				return
			}

			user, err := s.storage.GetUserByID(order.UserID)
			if err != nil {
				log.Errorf("error get user: %s", err)
				return
			}

			user.Balance = user.Balance + r.Accrual

			log.Debugf("New user balance %f", user.Balance)

			err = s.UpdateUser(user)
			if err != nil {
				log.Errorf("error update user: %s", err)
				return
			}

			order.Accrual = r.Accrual
			order.Status = r.Status

			err = s.UpdateOrder(order)
			if err != nil {
				log.Errorf("error update order: %s", err)
				return
			}

			log.Debug("Order updated")
		}()
	}
	return nil
}
