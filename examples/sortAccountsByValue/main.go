package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/dvcrn/pocketsmith-go"
)

func main() {
	token := os.Getenv("POCKETSMITH_TOKEN")
	if token == "" {
		log.Fatal("POCKETSMITH_TOKEN environment variable is required")
	}

	client := pocketsmith.NewClient(token)

	currentUser, err := client.GetCurrentUser()
	if currentUser == nil || err != nil {
		log.Fatal("Failed to get current user")
	}

	accounts, err := client.ListAccounts(currentUser.ID)
	if err != nil {
		log.Fatal(err)
	}

	// Sort accounts by current balance in descending order
	sort.Slice(accounts, func(i, j int) bool {
		// First check if either account is a credit card with non-zero balance
		isCreditCardI := accounts[i].Type == pocketsmith.AccountTypeCredits && accounts[i].CurrentBalance != 0
		isCreditCardJ := accounts[j].Type == pocketsmith.AccountTypeCredits && accounts[j].CurrentBalance != 0

		if isCreditCardI != isCreditCardJ {
			return isCreditCardI
		}

		// If neither or both are credit cards, sort by currency and balance
		if accounts[i].CurrencyCode == accounts[j].CurrencyCode {
			return accounts[i].CurrentBalance > accounts[j].CurrentBalance
		}

		return accounts[i].CurrentBalanceInBaseCurrency > accounts[j].CurrentBalanceInBaseCurrency
	})

	// print out all accounts in their new order, together with their current balance
	for _, account := range accounts {
		title := account.Title
		if len(title) > 30 {
			title = title[:27] + "..."
		}
		fmt.Printf("| %-30s | %15.2f | %5s |\n", title, account.CurrentBalance, account.CurrencyCode)
	}

	client.UpdateAccountsDisplayOrder(currentUser.ID, accounts)
}
