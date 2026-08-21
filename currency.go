package pocketsmith

import (
	"fmt"
	"net/http"
)

// CurrencySeparators are the formatting characters used for a currency.
type CurrencySeparators struct {
	Major string `json:"major"`
	Minor string `json:"minor"`
}

// Currency is a currency supported by PocketSmith.
type Currency struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Symbol     string             `json:"symbol"`
	MinorUnit  int                `json:"minor_unit"`
	Separators CurrencySeparators `json:"separators"`
}

// TimeZone is a time zone supported by PocketSmith.
type TimeZone struct {
	Name            string `json:"name"`
	UTCOffset       int    `json:"utc_offset"`
	FormattedName   string `json:"formatted_name"`
	FormattedOffset string `json:"formatted_offset"`
	Abbreviation    string `json:"abbreviation"`
	Identifier      string `json:"identifier"`
}

// ListCurrencies retrieves all currencies supported by PocketSmith.
func (c *Client) ListCurrencies() ([]*Currency, error) {
	url := "https://api.pocketsmith.com/v2/currencies"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")

	var currencies []*Currency
	if err := c.doAndDecode(req, &currencies); err != nil {
		return nil, err
	}

	return currencies, nil
}

// GetCurrency retrieves a single currency by its code, for example "nzd".
func (c *Client) GetCurrency(currencyID string) (*Currency, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/currencies/%s", currencyID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")

	var currency Currency
	if err := c.doAndDecode(req, &currency); err != nil {
		return nil, err
	}

	return &currency, nil
}

// ListTimeZones retrieves all time zones supported by PocketSmith.
func (c *Client) ListTimeZones() ([]*TimeZone, error) {
	url := "https://api.pocketsmith.com/v2/time_zones"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")

	var timeZones []*TimeZone
	if err := c.doAndDecode(req, &timeZones); err != nil {
		return nil, err
	}

	return timeZones, nil
}
