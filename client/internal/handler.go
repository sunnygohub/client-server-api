package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Response struct {
	Bid string `json:"bid"`
}

func FetchExchangeRate(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusInternalServerError {
		return "", fmt.Errorf("internal server error")
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.Bid, nil
}

func SaveTxtFile(value string) error {
	content := fmt.Sprintf("Dólar: %s", value)
	return os.WriteFile("cotacao.txt", []byte(content), 0644)
}
