package config

import (
	"flag"
	"fmt"
	"sync"

	"github.com/caarlos0/env/v6"
)

// адрес и порт запуска сервиса: переменная окружения ОС RUN_ADDRESS или флаг -a
// адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -d
// адрес системы расчёта начислений: переменная окружения ОС ACCRUAL_SYSTEM_ADDRESS или флаг -r

type cfg struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

var config cfg

func GetConfig() cfg {
	var once sync.Once

	once.Do(func() {

		if err := env.Parse(&config); err != nil {
			fmt.Printf("%+v\n", err)
		}

		flag.StringVar(&config.RunAddress, "a", "", "host:port")
		flag.StringVar(&config.DatabaseURI, "d", "", "database dsn")
		flag.StringVar(&config.AccrualSystemAddress, "r", "", "http://host:port")
		flag.Parse()
	})
	return config
}
