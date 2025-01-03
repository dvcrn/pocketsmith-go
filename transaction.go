package pocketsmith

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type CreateTransaction struct {
	Payee        string   `json:"payee"`
	Amount       float64  `json:"amount"`
	Date         string   `json:"date"`
	IsTransfer   bool     `json:"is_transfer"`
	Labels       []string `json:"labels,omitempty"`
	CategoryID   int      `json:"category_id,omitempty"`
	Note         string   `json:"note,omitempty"`
	Memo         string   `json:"memo,omitempty"`
	ChequeNumber string   `json:"cheque_number,omitempty"`
	NeedsReview  bool     `json:"needs_review"`
}

type Category struct {
	ID              int        `json:"id"`
	Title           string     `json:"title"`
	Colour          *string    `json:"colour"`
	IsTransfer      bool       `json:"is_transfer"`
	IsBill          bool       `json:"is_bill"`
	RefundBehaviour *string    `json:"refund_behaviour"`
	Children        []Category `json:"children"`
	ParentID        *int       `json:"parent_id"`
	RollUp          bool       `json:"roll_up"`
	CreatedAt       string     `json:"created_at"`
	UpdatedAt       string     `json:"updated_at"`
}

type ApiError struct {
	Err string `json:"error"`
}

func (a ApiError) Error() string {
	return fmt.Sprintf("Pocketsmith API Error: %s", a.Err)
}

type Transaction struct {
	ID                   int64              `json:"id"`
	Payee                string             `json:"payee"`
	OriginalPayee        string             `json:"original_payee"`
	Date                 string             `json:"date"`
	UploadSource         string             `json:"upload_source"`
	Category             Category           `json:"category"`
	ClosingBalance       float64            `json:"closing_balance"`
	ChequeNumber         string             `json:"cheque_number"`
	Memo                 string             `json:"memo"`
	Amount               float64            `json:"amount"`
	AmountInBaseCurrency float64            `json:"amount_in_base_currency"`
	Type                 string             `json:"type"`
	IsTransfer           bool               `json:"is_transfer"`
	NeedsReview          bool               `json:"needs_review"`
	Status               string             `json:"status"`
	Note                 string             `json:"note"`
	Labels               []string           `json:"labels"`
	TransactionAccount   TransactionAccount `json:"transaction_account"`
	CreatedAt            string             `json:"created_at"`
	UpdatedAt            string             `json:"updated_at"`
}

// AddTransaction creates a new transaction for the specified account.
// It takes an accountID and a CreateTransaction struct, and returns the created transaction and any error.
// The CreateTransaction struct contains the details of the new transaction to be created.
// The function makes a POST request to the PocketSmith API to create the new transaction.
func (c *Client) AddTransaction(accountID int, transaction *CreateTransaction) (*CreateTransaction, error) {
	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transaction_accounts/%d/transactions", accountID)

	payload, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	var createdTransaction CreateTransaction
	if err := c.doAndDecode(req, &createdTransaction); err != nil {
		return nil, err
	}

	return &createdTransaction, nil
}

func (c *Client) doAndDecode(req *http.Request, responseType any) error {
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-Developer-Key", c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBuf := new(bytes.Buffer)
	bodyBuf.ReadFrom(resp.Body)

	reader := bytes.NewReader(bodyBuf.Bytes())

	//b := bytes.Buffer{}
	//b.ReadFrom(reader)
	//fmt.Println("Body: ", string(b.Bytes()))
	//reader.Seek(0, 0)

	err = json.NewDecoder(reader).Decode(responseType)
	if err != nil {
		// reset reading position of the buffer
		reader.Seek(0, 0)

		// try to decode into ApiError instead
		var apiError *ApiError
		if err := json.NewDecoder(reader).Decode(&apiError); err != nil {
			return err
		}

		if apiError != nil {
			return apiError
		}
	}

	return nil
}

// SearchTransactions retrieves a list of transactions for the specified account, with optional filtering by start date, end date, and search query.
func (c *Client) SearchTransactions(accountID int, startDate, endDate, search string) ([]*Transaction, error) {
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

	var transactions []*Transaction
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

func (c *Client) ListTransactions(accountID int, opts ...ListTransactionsOption) ([]*Transaction, error) {
	options := &listTransactionsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	url := fmt.Sprintf("https://api.pocketsmith.com/v2/transaction_accounts/%d/transactions", accountID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

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

	var transactions []*Transaction
	if err := c.doAndDecode(req, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

// UpdateTransaction updates an existing transaction with the provided transaction data.
// It takes the ID of the transaction to update and a pointer to a CreateTransaction struct
// containing the updated transaction data. It returns an error if the update fails.
func (c *Client) UpdateTransaction(transactionID int64, transaction *CreateTransaction) (*Transaction, error) {
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

	var tx *Transaction
	if err := c.doAndDecode(req, &tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// SearchTransactionsByMemo searches for transactions by the memo field within a given date range.
// It takes an accountID, a referenceNo string to search for in the memo field, and a transactionDate time.Time.
// It returns a slice of matching Transaction pointers, or an error if the search fails.
func (c *Client) SearchTransactionsByMemo(accountID int, transactionDate time.Time, search string) ([]*Transaction, error) {
	startDate := transactionDate.Add(-1 * 24 * time.Hour).Format("2006-01-02")
	endDate := transactionDate.Add(1 * 24 * time.Hour).Format("2006-01-02")

	transactions, err := c.SearchTransactions(accountID, startDate, endDate, "")
	if err != nil {
		return nil, fmt.Errorf("error searching for transactions: %v", err)
	}

	var matchingTransactions []*Transaction
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
func (c *Client) SearchTransactionsByMemoContains(accountID int, transactionDate time.Time, search string) ([]*Transaction, error) {
	startDate := transactionDate.Add(-1 * 24 * time.Hour).Format("2006-01-02")
	endDate := transactionDate.Add(1 * 24 * time.Hour).Format("2006-01-02")

	transactions, err := c.SearchTransactions(accountID, startDate, endDate, "")
	if err != nil {
		return nil, fmt.Errorf("error searching for transactions: %v", err)
	}

	var matchingTransactions []*Transaction
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
func (c *Client) SearchTransactionsByChequeNumber(accountID int, transactionDate time.Time, chequeNum string) ([]*Transaction, error) {
	startDate := transactionDate.Add(-1 * 24 * time.Hour).Format("2006-01-02")
	endDate := transactionDate.Add(1 * 24 * time.Hour).Format("2006-01-02")

	transactions, err := c.SearchTransactions(accountID, startDate, endDate, "")
	if err != nil {
		return nil, fmt.Errorf("error searching for transactions: %v", err)
	}

	var matchingTransactions []*Transaction
	for _, tx := range transactions {
		if tx.ChequeNumber == chequeNum {
			matchingTransactions = append(matchingTransactions, tx)
		}
	}

	return matchingTransactions, nil
}
