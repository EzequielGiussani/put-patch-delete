package main

import (
	"app/internal/application"
	"fmt"
)

func main() {
	app := application.NewServerChi(":8080")

	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
