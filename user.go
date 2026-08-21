package pocketsmith

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID                         int    `json:"id"`
	Login                      string `json:"login"`
	Name                       string `json:"name"`
	Email                      string `json:"email"`
	AvatarURL                  string `json:"avatar_url"`
	BetaUser                   bool   `json:"beta_user"`
	DemoUser                   bool   `json:"demo_user"`
	CountryCode                string `json:"country_code"`
	TimeZone                   string `json:"time_zone"`
	WeekStartDay               int    `json:"week_start_day"`
	MonthStartsOnDay           int    `json:"month_starts_on_day"`
	YearStartsOnMonth          int    `json:"year_starts_on_month"`
	IsReviewingTransactions    bool   `json:"is_reviewing_transactions"`
	BaseCurrencyCode           string `json:"base_currency_code"`
	AlwaysShowBaseCurrency     bool   `json:"always_show_base_currency"`
	UsingMultipleCurrencies    bool   `json:"using_multiple_currencies"`
	UsingNewTransactionsSearch bool   `json:"using_new_transactions_search"`
	AvailableAccounts          int    `json:"available_accounts"`
	AvailableBudgets           int    `json:"available_budgets"`
	AtDashboardLimit           bool   `json:"at_dashboard_limit"`
	AllowedDataFeeds           bool   `json:"allowed_data_feeds"`
	// TellAFriendAccess and TellAFriendCode are returned by the API but are
	// undocumented, and were null for every account inspected. They are kept
	// raw so their value is preserved without guessing at their type.
	TellAFriendAccess        json.RawMessage `json:"tell_a_friend_access"`
	TellAFriendCode          json.RawMessage `json:"tell_a_friend_code"`
	ForecastLastUpdatedAt    string          `json:"forecast_last_updated_at"`
	ForecastLastAccessedAt   string          `json:"forecast_last_accessed_at"`
	ForecastStartDate        string          `json:"forecast_start_date"`
	ForecastEndDate          string          `json:"forecast_end_date"`
	ForecastDeferRecalculate bool            `json:"forecast_defer_recalculate"`
	ForecastNeedsRecalculate bool            `json:"forecast_needs_recalculate"`
	FeedHistoryStartsFrom    string          `json:"feed_history_starts_from"`
	FeedHistoryTouched       bool            `json:"feed_history_touched"`
	LastLoggedInAt           string          `json:"last_logged_in_at"`
	LastActivityAt           string          `json:"last_activity_at"`
	CreatedAt                string          `json:"created_at"`
	UpdatedAt                string          `json:"updated_at"`
}

// Label is a transaction label belonging to a user.
type Label string

// SavedSearch is a saved transaction search belonging to a user.
type SavedSearch struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// UpdateUser holds the fields accepted by PUT /users/{id}. Fields left empty
// are omitted from the request and are left untouched by the API.
type UpdateUser struct {
	Name                   string `json:"name,omitempty"`
	TimeZone               string `json:"time_zone,omitempty"`
	WeekStartDay           *int   `json:"week_start_day,omitempty"`
	BetaUser               *bool  `json:"beta_user,omitempty"`
	BaseCurrencyCode       string `json:"base_currency_code,omitempty"`
	AlwaysShowBaseCurrency *bool  `json:"always_show_base_currency,omitempty"`
}

func (c *Client) GetCurrentUser() (*User, error) {
	url := "https://api.pocketsmith.com/v2/me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("X-Developer-Key", c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUser retrieves a user by their ID.
func (c *Client) GetUser(userID int) (*User, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d", userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")

	var user User
	if err := c.doAndDecode(req, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates a user's preferences.
func (c *Client) UpdateUser(userID int, update *UpdateUser) (*User, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d", userID)

	payload, err := json.Marshal(update)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	var user User
	if err := c.doAndDecode(req, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteForecastCache clears the cached forecast data for a user, causing it
// to be recalculated.
func (c *Client) DeleteForecastCache(userID int) error {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d/forecast_cache", userID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")

	return c.doAndDecode(req, nil)
}

// ListLabels retrieves all transaction labels for a user.
func (c *Client) ListLabels(userID int) ([]Label, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d/labels", userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")

	var labels []Label
	if err := c.doAndDecode(req, &labels); err != nil {
		return nil, err
	}

	return labels, nil
}

// ListSavedSearches retrieves all saved searches for a user.
func (c *Client) ListSavedSearches(userID int) ([]*SavedSearch, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d/saved_searches", userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")

	var savedSearches []*SavedSearch
	if err := c.doAndDecode(req, &savedSearches); err != nil {
		return nil, err
	}

	return savedSearches, nil
}
