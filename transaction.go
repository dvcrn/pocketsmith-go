package pocketsmith

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Transaction struct {
	Payee        string     `json:"payee"`
	Amount       float64    `json:"amount"`
	Date         string     `json:"date"`
	IsTransfer   bool       `json:"is_transfer"`
	Labels       []string   `json:"labels,omitempty"`
	CategoryID   CategoryID `json:"category_id,omitempty"`
	Note         string     `json:"note,omitempty"`
	Memo         string     `json:"memo,omitempty"`
	ChequeNumber string     `json:"cheque_number,omitempty"`
	NeedsReview  bool       `json:"needs_review"`
}

type DetailedTransaction struct {
	ID                   int64               `json:"id"`
	Payee                string              `json:"payee"`
	OriginalPayee        string              `json:"original_payee"`
	Date                 string              `json:"date"`
	UploadSource         string              `json:"upload_source"`
	Category             *Category           `json:"category"`
	ClosingBalance       float64             `json:"closing_balance"`
	ChequeNumber         string              `json:"cheque_number"`
	Memo                 string              `json:"memo"`
	Amount               float64             `json:"amount"`
	AmountInBaseCurrency float64             `json:"amount_in_base_currency"`
	Type                 string              `json:"type"`
	IsTransfer           bool                `json:"is_transfer"`
	NeedsReview          bool                `json:"needs_review"`
	Status               string              `json:"status"`
	Note                 string              `json:"note"`
	Labels               []string            `json:"labels"`
	TransactionAccount   *TransactionAccount `json:"transaction_account"`
	CreatedAt            string              `json:"created_at"`
	UpdatedAt            string              `json:"updated_at"`
}

// AddTransaction creates a new transaction for the specified account.
// It takes an accountID and a CreateTransaction struct, and returns the created transaction and any error.
// The CreateTransaction struct contains the details of the new transaction to be created.
// The function makes a POST request to the PocketSmith API to create the new transaction.
func (c *Client) AddTransaction(transactionAccountID int, transaction *Transaction) (*Transaction, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transaction_accounts/%d/transactions", transactionAccountID)

	payload, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	var createdTransaction Transaction
	if err := c.doAndDecode(req, &createdTransaction); err != nil {
		return nil, err
	}

	return &createdTransaction, nil
}

// SearchTransactions retrieves a list of transactions for the specified account, with optional filtering by start date, end date, and search query.
func (c *Client) SearchTransactions(accountID int, startDate, endDate, search string) ([]*DetailedTransaction, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transaction_accounts/%d/transactions", accountID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if startDate != "" {
		q.Add("start_date", startDate)
	}
	if endDate != "" {
		q.Add("end_date", endDate)
	}
	if search != "" {
		q.Add("search", search)
	}
	req.URL.RawQuery = q.Encode()

	var transactions []*DetailedTransaction
	if err := c.doAndDecode(req, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

// ListTransactions retrieves a list of transactions for the specified account, with optional filtering by date range, update time, categorization, transaction type, review status, and search query. The results are paginated, with the page number specified as a parameter.
type ListTransactionsOption func(*listTransactionsOptions)

type listTransactionsOptions struct {
	startDate       string
	endDate         string
	updatedSince    string
	uncategorised   int
	transactionType string
	needsReview     int
	search          string
	page            int
}

func WithStartDate(date string) ListTransactionsOption {
	return func(o *listTransactionsOptions) {
		o.startDate = date
	}
}

func WithEndDate(date string) ListTransactionsOption {
	return func(o *listTransactionsOptions) {
		o.endDate = date
	}
}

func WithUpdatedSince(date string) ListTransactionsOption {
	return func(o *listTransactionsOptions) {
		o.updatedSince = date
	}
}

func WithUncategorised(uncategorised int) ListTransactionsOption {
	return func(o *listTransactionsOptions) {
		o.uncategorised = uncategorised
	}
}

func WithTransactionType(transactionType string) ListTransactionsOption {
	return func(o *listTransactionsOptions) {
		o.transactionType = transactionType
	}
}

func WithNeedsReview(needsReview int) ListTransactionsOption {
	return func(o *listTransactionsOptions) {
		o.needsReview = needsReview
	}
}

func WithSearch(search string) ListTransactionsOption {
	return func(o *listTransactionsOptions) {
		o.search = search
	}
}

func WithPage(page int) ListTransactionsOption {
	return func(o *listTransactionsOptions) {
		o.page = page
	}
}

func (c *Client) ListTransactions(accountID int, opts ...ListTransactionsOption) ([]*DetailedTransaction, error) {
	options := &listTransactionsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transaction_accounts/%d/transactions", accountID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	applyListTransactionsOptions(req, options)

	var transactions []*DetailedTransaction
	if err := c.doAndDecode(req, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

// UpdateTransaction updates an existing transaction with the provided transaction data.
// Setting CategoryIDNone will remove the transaction's category.
func (c *Client) UpdateTransaction(transactionID int64, transaction *Transaction) (*DetailedTransaction, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transactions/%d", transactionID)

	payload, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")

	var tx *DetailedTransaction
	if err := c.doAndDecode(req, &tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// SearchTransactionsByMemo searches for transactions by the memo field within a given date range.
// It takes an accountID, a referenceNo string to search for in the memo field, and a transactionDate time.Time.
// It returns a slice of matching Transaction pointers, or an error if the search fails.
func (c *Client) SearchTransactionsByMemo(accountID int, transactionDate time.Time, search string) ([]*DetailedTransaction, error) {
	startDate := transactionDate.Add(-1 * 24 * time.Hour).Format("2006-01-02")
	endDate := transactionDate.Add(1 * 24 * time.Hour).Format("2006-01-02")

	transactions, err := c.SearchTransactions(accountID, startDate, endDate, "")
	if err != nil {
		return nil, fmt.Errorf("error searching for transactions: %v", err)
	}

	var matchingTransactions []*DetailedTransaction
	for _, tx := range transactions {
		if tx.Memo == search {
			matchingTransactions = append(matchingTransactions, tx)
		}
	}

	return matchingTransactions, nil
}

// SearchTransactionsByMemoContains searches for transactions by the memo field within a given date range,
// where the memo contains the specified search string.
// It takes an accountID, a search string to look for in the memo field, and a transactionDate time.Time.
// It returns a slice of matching Transaction pointers, or an error if the search fails.
func (c *Client) SearchTransactionsByMemoContains(accountID int, transactionDate time.Time, search string) ([]*DetailedTransaction, error) {
	startDate := transactionDate.Add(-1 * 24 * time.Hour).Format("2006-01-02")
	endDate := transactionDate.Add(1 * 24 * time.Hour).Format("2006-01-02")

	transactions, err := c.SearchTransactions(accountID, startDate, endDate, "")
	if err != nil {
		return nil, fmt.Errorf("error searching for transactions: %v", err)
	}

	var matchingTransactions []*DetailedTransaction
	for _, tx := range transactions {
		if strings.Contains(tx.Memo, search) {
			matchingTransactions = append(matchingTransactions, tx)
		}
	}

	return matchingTransactions, nil
}

// SearchTransactionsByChequeNumber searches for transactions by the cheque number within a given date range.
// It takes an accountID, a transactionDate time.Time, and a chequeNum string to search for.
// It returns a slice of matching Transaction pointers, or an error if the search fails.
func (c *Client) SearchTransactionsByChequeNumber(accountID int, transactionDate time.Time, chequeNum string) ([]*DetailedTransaction, error) {
	startDate := transactionDate.Add(-1 * 24 * time.Hour).Format("2006-01-02")
	endDate := transactionDate.Add(1 * 24 * time.Hour).Format("2006-01-02")

	transactions, err := c.SearchTransactions(accountID, startDate, endDate, "")
	if err != nil {
		return nil, fmt.Errorf("error searching for transactions: %v", err)
	}

	var matchingTransactions []*DetailedTransaction
	for _, tx := range transactions {
		if tx.ChequeNumber == chequeNum {
			matchingTransactions = append(matchingTransactions, tx)
		}
	}

	return matchingTransactions, nil
}

func applyListTransactionsOptions(req *http.Request, options *listTransactionsOptions) {
	q := req.URL.Query()
	if options.startDate != "" {
		q.Add("start_date", options.startDate)
	}
	if options.endDate != "" {
		q.Add("end_date", options.endDate)
	}
	if options.updatedSince != "" {
		q.Add("updated_since", options.updatedSince)
	}
	if options.uncategorised > 0 {
		q.Add("uncategorised", fmt.Sprintf("%d", options.uncategorised))
	}
	if options.transactionType != "" {
		q.Add("type", options.transactionType)
	}
	if options.needsReview > 0 {
		q.Add("needs_review", fmt.Sprintf("%d", options.needsReview))
	}
	if options.search != "" {
		q.Add("search", options.search)
	}
	if options.page > 0 {
		q.Add("page", fmt.Sprintf("%d", options.page))
	}
	req.URL.RawQuery = q.Encode()
}

func (c *Client) listTransactions(url string, opts ...ListTransactionsOption) ([]*DetailedTransaction, error) {
	options := &listTransactionsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	applyListTransactionsOptions(req, options)

	var transactions []*DetailedTransaction
	if err := c.doAndDecode(req, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

// ListTransactionsInUser retrieves a list of transactions across all of the
// user's accounts.
func (c *Client) ListTransactionsInUser(userID int, opts ...ListTransactionsOption) ([]*DetailedTransaction, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/users/%d/transactions", userID)
	return c.listTransactions(url, opts...)
}

// ListTransactionsInAccount retrieves a list of transactions in an account.
func (c *Client) ListTransactionsInAccount(accountID int, opts ...ListTransactionsOption) ([]*DetailedTransaction, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/accounts/%d/transactions", accountID)
	return c.listTransactions(url, opts...)
}

// ListTransactionsInTransactionAccount retrieves a list of transactions in a
// transaction account. It is the same as ListTransactions, but named to make
// the resource it operates on unambiguous.
func (c *Client) ListTransactionsInTransactionAccount(transactionAccountID int, opts ...ListTransactionsOption) ([]*DetailedTransaction, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transaction_accounts/%d/transactions", transactionAccountID)
	return c.listTransactions(url, opts...)
}

// ListTransactionsInCategories retrieves a list of transactions in one or more
// categories.
func (c *Client) ListTransactionsInCategories(categoryIDs []int, opts ...ListTransactionsOption) ([]*DetailedTransaction, error) {
	ids := make([]string, 0, len(categoryIDs))
	for _, id := range categoryIDs {
		ids = append(ids, fmt.Sprintf("%d", id))
	}

	url := fmt.Sprintf("https://api.pocketsmith.com/v2/categories/%s/transactions", strings.Join(ids, ","))
	return c.listTransactions(url, opts...)
}

// GetTransaction retrieves a single transaction by its ID.
func (c *Client) GetTransaction(transactionID int64) (*DetailedTransaction, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transactions/%d", transactionID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")

	var transaction DetailedTransaction
	if err := c.doAndDecode(req, &transaction); err != nil {
		return nil, err
	}

	return &transaction, nil
}

// DeleteTransaction deletes a transaction by its ID.
func (c *Client) DeleteTransaction(transactionID int64) error {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transactions/%d", transactionID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")

	return c.doAndDecode(req, nil)
}
