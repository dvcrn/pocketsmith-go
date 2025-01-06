package pocketsmith

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CategoryID int

const CategoryIDNone CategoryID = -1

func (c CategoryID) MarshalJSON() ([]byte, error) {
	if c == CategoryIDNone {
		return []byte(`""`), nil
	}
	return json.Marshal(int(c))
}

type Category struct {
	ID              int         `json:"id"`
	Title           string      `json:"title"`
	Colour          string      `json:"colour"`
	IsTransfer      bool        `json:"is_transfer"`
	IsBill          bool        `json:"is_bill"`
	RefundBehaviour string      `json:"refund_behaviour"`
	Children        []*Category `json:"children"`
	ParentID        int         `json:"parent_id"`
	RollUp          bool        `json:"roll_up"`
	CreatedAt       string      `json:"created_at"`
	UpdatedAt       string      `json:"updated_at"`
}

type CategoryRule struct {
	ID           int64     `json:"id"`
	Category     *Category `json:"category"`
	PayeeMatches string    `json:"payee_matches"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

func (rule *CategoryRule) Matches(target string) bool {
	return strings.Contains(target, rule.PayeeMatches)
}

// GetCategoryRules retrieves all category rules for a given user
func (c *Client) ListCategoryRules(userID int) ([]*CategoryRule, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d/category_rules", userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var rules []*CategoryRule
	if err := c.doAndDecode(req, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}
