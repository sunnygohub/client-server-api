package main

import (
	"context"
	"log"
	"time"

	"github.com/sunnygohub/client-server-api/client/internal"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	rate, err := internal.FetchExchangeRate(ctx)
	if err != nil {
		log.Fatal("Error fetching exchange rate: ", err)
	}

	err = internal.SaveTxtFile(rate)
	if err != nil {
		log.Fatal("Error saving file: ", err)
	}

	log.Println("Exchange rate saved to cotacao.txt")
}
