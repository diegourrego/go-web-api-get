package main

import (
	"first_api/internal/application"
	"fmt"
)

func main() {
	app := application.NewDefaultHTTP(":8080")
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
