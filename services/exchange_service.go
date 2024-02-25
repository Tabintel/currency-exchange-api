// services/exchange_service.go

package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// ExchangeService represents the service for fetching exchange rate data from external APIs
type ExchangeService struct {
	APIKeyA string // API key for Service A
	APIKeyB string // API key for Service B
}


// ExchangeRateResponse represents the response structure for exchange rate data
type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

// FetchExchangeRates fetches exchange rate data from two external services concurrently
func (es *ExchangeService) FetchExchangeRates(currencyPair string) (float64, error) {
	// Channel to receive responses from the external services
	respChan := make(chan float64, 2)
	errChan := make(chan error, 2)

	// Fetch exchange rate from Service A
	go func() {
		rate, err := es.fetchExchangeRateFromServiceA(currencyPair)
		if err != nil {
			errChan <- err
		} else {
			respChan <- rate
		}
	}()

	// Fetch exchange rate from Service B
	go func() {
		rate, err := es.fetchExchangeRateFromServiceB(currencyPair)
		if err != nil {
			errChan <- err
		} else {
			respChan <- rate
		}
	}()

	// Wait for the first response
	select {
	case rate := <-respChan:
		return rate, nil
	case err := <-errChan:
		return 0, err
	}
}

// fetchExchangeRateFromServiceA fetches exchange rate data from Service A (Open Exchange Rates API)
func (es *ExchangeService) fetchExchangeRateFromServiceA(currencyPair string) (float64, error) {
	url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s", es.APIKeyA)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch exchange rate from Service A: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body from Service A: %w", err)
	}

	var exchangeRateResponse ExchangeRateResponse
	if err := json.Unmarshal(body, &exchangeRateResponse); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response from Service A: %w", err)
	}

	rate, ok := exchangeRateResponse.Rates[currencyPair]
	if !ok {
		return 0, fmt.Errorf("exchange rate for currency pair %s not found in response from Service A", currencyPair)
	}

	return rate, nil
}

// fetchExchangeRateFromServiceB fetches exchange rate data from Service B (ExchangeRatesAPI.io)
func (es *ExchangeService) fetchExchangeRateFromServiceB(currencyPair string) (float64, error) {
	url := fmt.Sprintf("https://api.exchangeratesapi.io/latest?base=USD&access_key=%s", es.APIKeyB)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch exchange rate from Service B: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body from Service B: %w", err)
	}

	var exchangeRateResponse ExchangeRateResponse
	if err := json.Unmarshal(body, &exchangeRateResponse); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response from Service B: %w", err)
	}

	rate, ok := exchangeRateResponse.Rates[currencyPair]
	if !ok {
		return 0, fmt.Errorf("exchange rate for currency pair %s not found in response from Service B", currencyPair)
	}

	return rate, nil
}
