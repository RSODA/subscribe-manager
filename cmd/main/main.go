package main

import (
	"fmt"

	"github.com/RSODA/subscribe-manager/internal/app"
)

// @title           Subscription Manager API
// @version         1.0
// @description     REST сервис для агрегации данных об онлайн подписках
// @host            localhost:8080
// @BasePath        /

func main() {
	if err := app.Run(); err != nil {
		fmt.Printf("Application stopped with error: %v\n", err)
		panic(err)
	}
}
