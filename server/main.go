package main

import (
	"log"
	"net/http"

	"github.com/sunnygohub/client-server-api/server/internal"
)

func main() {
	dbConn, err := internal.InitializeDatabase("./exchange_rates.db")
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}
	defer dbConn.Close()

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		internal.CotacaoHandler(w, r, dbConn)
	})

	log.Println("Server Running")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
