package mock

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestMockServer(t *testing.T) {
	// Start the mock server
	baseURL := StartMockServer()

	// Test GET request to /third-party/payments/payment123
	t.Run("GET /third-party/payments/payment123", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/third-party/payments/payment123")
		if err != nil {
			t.Fatalf("GET request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
		}
	})

	// Test POST request to /third-party/payments
	t.Run("POST /third-party/payments", func(t *testing.T) {
		// Define the request body
		requestBody := map[string]interface{}{
			"account_id": "1234567890",
			"reference":  "payment123",
			"amount":     100.50,
		}

		// Convert request data to JSON
		reqBody, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}

		// Create a POST request with the JSON request body
		req, err := http.NewRequest(http.MethodPost, baseURL+"/third-party/payments", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to create POST request: %v", err)
		}

		// Set content type header
		req.Header.Set("Content-Type", "application/json")

		// Perform the POST request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("POST request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
		}
	})
}
