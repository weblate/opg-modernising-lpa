package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/env"
)

func main() {
	port := env.Get("PORT", "8080")

	http.HandleFunc("/v2/notifications/email", func(w http.ResponseWriter, r *http.Request) {
		var v map[string]interface{}
		json.NewDecoder(r.Body).Decode(&v)
		log.Println("email:", v)
		json.NewEncoder(w).Encode(map[string]string{"id": "an-email-id"})
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
