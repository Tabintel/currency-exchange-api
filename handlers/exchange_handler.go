// handlers/exchange_handler.go

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/tabintel/currency-exchange-api/services"
)

// ExchangeRequest represents the JSON request body for the exchange rate endpoint
type ExchangeRequest struct {
	CurrencyPair string `json:"currency-pair"`
}

// ExchangeRateHandler handles requests to retrieve exchange rates for a given currency pair
func ExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req ExchangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Validate the currency pair format
	if !isValidCurrencyPair(req.CurrencyPair) {
		http.Error(w, "Invalid currency pair format", http.StatusBadRequest)
		return
	}

	// Create an instance of ExchangeService
	exchangeService := services.NewExchangeService()

	// Fetch exchange rates for the given currency pair
	rate, err := exchangeService.FetchExchangeRates(req.CurrencyPair)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch exchange rates: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the exchange rate as JSON response
	response := map[string]float64{req.CurrencyPair: rate}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// isValidCurrencyPair checks if the currency pair has a valid format (e.g., "USD-EUR")
func isValidCurrencyPair(currencyPair string) bool {
	// Use regular expression to validate currency pair format
	pattern := `^[A-Z]{3}-[A-Z]{3}$`
	match, _ := regexp.MatchString(pattern, currencyPair)
	return match
}
