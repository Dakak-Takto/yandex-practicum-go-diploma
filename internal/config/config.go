package config

import (
	"flag"
	"fmt"
	"sync"

	"github.com/caarlos0/env/v6"
)

type cfg struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	CookieStoreKey       string `env:"COOKIE_STORE_KEY"`
}

var config cfg

func InitConfig() {
	var once sync.Once

	once.Do(func() {

		if err := env.Parse(&config); err != nil {
			fmt.Printf("%+v\n", err)
		}

		flag.StringVar(&config.RunAddress, "a", config.RunAddress, "host:port")
		flag.StringVar(&config.DatabaseURI, "d", config.DatabaseURI, "database dsn")
		flag.StringVar(&config.AccrualSystemAddress, "r", config.AccrualSystemAddress, "http://host:port")
		flag.StringVar(&config.CookieStoreKey, "k", config.CookieStoreKey, "secret key for session cookie")
		flag.Parse()
	})
}

func RunAddr() string {
	return config.RunAddress
}

func DatabaseURI() string {
	return config.DatabaseURI
}

func AccrualSystemAddress() string {
	return config.AccrualSystemAddress
}

func CookieStoreKey() string {
	return config.CookieStoreKey
}
