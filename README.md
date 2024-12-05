# pocketsmith-go 

CLI wrapper in Golang for pocketsmith API https://developers.pocketsmith.com/reference/get_me-1

## Supported API methods (so far) 

### User
- Get current user information (`GetCurrentUser`)

### Institution
- Create a new institution (`CreateInstitution`)
- List all institutions (`ListInstitutions`)
- Find institution by name (`FindInstitutionByName`)

### Account
- List all accounts (`ListAccounts`)
- List all transaction accounts (`ListTransactionAccounts`)
- Create a new account (`CreateAccount`)
- Find account by name (`FindAccountByName`)
- Update transaction account (`UpdateTransactionAccount`)

### Transaction
- Add a new transaction (`AddTransaction`)
- Search transactions (`SearchTransactions`)
- List transactions with filters (`ListTransactions`)

## Examples


// Get current user
user, err := client.GetCurrentUser()

// Create an institution
institution, err := client.CreateInstitution(userID, "Bank Name", "usd")

// Create an account
account, err := client.CreateAccount(userID, institutionID, "Savings", "usd", AccountTypeBank)

// Add a transaction
transaction := &CreateTransaction{
    Payee:      "Store Name",
    Amount:     -50.00,
    Date:       "2024-01-01",
    IsTransfer: false,
}
result, err := client.AddTransaction(accountID, transaction)


## License

MIT