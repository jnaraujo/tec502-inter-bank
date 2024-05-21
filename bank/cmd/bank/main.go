package main

import (
	"fmt"

	"github.com/jnaraujo/tec502-inter-bank/bank/internal/config"
	"github.com/jnaraujo/tec502-inter-bank/bank/internal/http"
)

func main() {
	fmt.Println("Bank App")

	err := config.InitEnv()
	if err != nil {
		panic(err)
	}

	err = http.NewServer()
	if err != nil {
		panic(err)
	}
}
