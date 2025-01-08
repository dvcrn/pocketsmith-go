package pocketsmith

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Institution struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	CurrencyCode string `json:"currency_code"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func (c *Client) CreateInstitution(userID int, title string, currencyCode string) (*Institution, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d/institutions", userID)

	payload := struct {
		Title        string `json:"title"`
		CurrencyCode string `json:"currency_code"`
	}{
		Title:        title,
		CurrencyCode: currencyCode,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	var institution Institution
	err = c.doAndDecode(req, &institution)
	if err != nil {
		return nil, err
	}

	return &institution, nil
}

func (c *Client) ListInstitutions(userID int) ([]*Institution, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d/institutions", userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")

	var institutions []*Institution
	err = c.doAndDecode(req, &institutions)
	if err != nil {
		return nil, err
	}

	return institutions, nil
}

func (c *Client) FindInstitutionByName(userID int, name string) (*Institution, error) {
	institutions, err := c.ListInstitutions(userID)
	if err != nil {
		return nil, err
	}

	for _, institution := range institutions {
		if institution.Title == name {
			return institution, nil
		}
	}

	return nil, ErrNotFound
}

func (c *Client) FindInstitutionsByNameContains(userID int, name string) ([]*Institution, error) {
	institutions, err := c.ListInstitutions(userID)
	if err != nil {
		return nil, err
	}

	found := make([]*Institution, 0)
	for _, institution := range institutions {
		if strings.Contains(institution.Title, name) {
			found = append(found, institution)
		}
	}

	return found, nil
}

func (c *Client) DeleteInstitution(institutionID int, mergeIntoInstitutionID int) error {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/institutions/%d", institutionID)
	if mergeIntoInstitutionID > 0 {
		url = fmt.Sprintf("%s?merge_into_institution_id=%d", url, mergeIntoInstitutionID)
	}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")

	err = c.doAndDecode(req, nil)
	if err != nil {
		return err
	}

	return nil
}
