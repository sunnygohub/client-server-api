package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ExchangeRate struct {
	Bid        string    `json:"bid"`
	CreateDate time.Time `json:"create_date"`
}

func getExchangeRate(ctx context.Context) (ExchangeRate, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return ExchangeRate{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ExchangeRate{}, err
	}
	defer resp.Body.Close()

	var result map[string]ExchangeRate
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return ExchangeRate{}, err
	}

	exchangeRate, ok := result["USDBRL"]
	if !ok {
		return ExchangeRate{}, fmt.Errorf("key 'USDBRL' not found in response")
	}

	return exchangeRate, nil
}

func CotacaoHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	rate, err := getExchangeRate(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching exchange rate: %v", err), http.StatusInternalServerError)
		log.Println("Fetch Error: ", err)
		return
	}

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer dbCancel()

	err = InsertExchangeRate(dbCtx, dbConn, rate)
	if err != nil {
		log.Println("DB Error: ", err)
	}

	response := map[string]string{"bid": rate.Bid}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
