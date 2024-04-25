package mock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

// StartMockServer starts a mock server with the provided router and returns the base URL
func StartMockServer() string {
	router := mux.NewRouter()

	// Mock handler for POST /third-party/payments
	router.HandleFunc("/third-party/payments", func(w http.ResponseWriter, r *http.Request) {
		var payment struct {
			AccountID string  `json:"account_id"`
			Reference string  `json:"reference"`
			Amount    float64 `json:"amount"`
		}
		err := json.NewDecoder(r.Body).Decode(&payment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Respond with the same payment details
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payment)
	}).Methods("POST")

	// Mock handler for GET /third-party/payments/:reference
	router.HandleFunc("/third-party/payments/{reference}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		reference := vars["reference"]

		// Mock response based on the reference ID
		payment := struct {
			AccountID string  `json:"account_id"`
			Reference string  `json:"reference"`
			Amount    float64 `json:"amount"`
		}{
			AccountID: "1234567890",
			Reference: reference,
			Amount:    100.50,
		}
		// Respond with the mock payment details
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payment)
	}).Methods("GET")

	// Start a test server using the router
	ts := httptest.NewServer(router)
	return ts.URL
}
