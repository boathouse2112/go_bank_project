package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	bank "github.com/boathouse2112/bank/bankcore"
)

var accounts = map[int]*JsonAccount{}

type JsonAccount struct {
	*bank.Account
}

func (a *JsonAccount) Statement() string {
	json, err := json.Marshal(a)
	if err != nil {
		return err.Error()
	}
	return string(json)
}

func deposit(w http.ResponseWriter, req *http.Request) {
	numberqs := req.URL.Query().Get("number")
	amountqs := req.URL.Query().Get("amount")

	if numberqs == "" {
		fmt.Fprintf(w, "Account number is missing!")
		return
	}

	if number, err := strconv.Atoi(numberqs); err != nil {
		fmt.Fprintf(w, "Invalid account number!")
	} else if amount, err := strconv.ParseFloat(amountqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid amount number!")
	} else {
		account, ok := accounts[number]
		if !ok {
			fmt.Fprintf(w, "Account with number %v can't be found!", number)
		} else {
			err := account.Deposit(amount)
			if err != nil {
				fmt.Fprintf(w, "%v", err)
			} else {
				fmt.Fprint(w, account.Statement())
			}
		}
	}
}

func withdraw(w http.ResponseWriter, req *http.Request) {
	numberqs := req.URL.Query().Get("number")
	amountqs := req.URL.Query().Get("amount")

	if numberqs == "" {
		fmt.Fprintf(w, "Account number is missing!")
		return
	}

	if number, err := strconv.Atoi(numberqs); err != nil {
		fmt.Fprintf(w, "Invalid account number!")
	} else if amount, err := strconv.ParseFloat(amountqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid amount number!")
	} else {
		account, ok := accounts[number]
		if !ok {
			fmt.Fprintf(w, "Account with number %v can't be found!", number)
		} else {
			err := account.Withdraw(amount)
			if err != nil {
				fmt.Fprintf(w, "%v", err)
			} else {
				fmt.Fprint(w, account.Statement())
			}
		}
	}
}

func statement(w http.ResponseWriter, req *http.Request) {
	numberqs := req.URL.Query().Get("number")

	if numberqs == "" {
		fmt.Fprintf(w, "Account number is missing!")
		return
	}

	number, err := strconv.Atoi(numberqs)
	if err != nil {
		fmt.Printf("Invalid account number!")
	} else {
		account, ok := accounts[number]
		if !ok {
			fmt.Fprintf(w, "Account with number %v can't be found.", numberqs)
		} else {
			json.NewEncoder(w).Encode(bank.Statement(account))
		}
	}
}

func Transfer(w http.ResponseWriter, req *http.Request) {
	fromnumberqs := req.URL.Query().Get("fromnumber")
	tonumberqs := req.URL.Query().Get("fromnumber")
	amountqs := req.URL.Query().Get("amountqs")

	if fromnumberqs == "" {
		fmt.Fprintf(w, "From-account number is missing!")
		return
	}
	if tonumberqs == "" {
		fmt.Fprintf(w, "To-account number is missing!")
		return
	}
	if amountqs == "" {
		fmt.Fprintf(w, "Amount number is missing!")
		return
	}

	if fromNumber, err := strconv.Atoi(fromnumberqs); err != nil {
		fmt.Fprintf(w, "Invalid from-account number!")
	} else if toNumber, err := strconv.Atoi(tonumberqs); err != nil {
		fmt.Fprintf(w, "Invalid to-account number!")
	} else if amount, err := strconv.ParseFloat(amountqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid amount number!")
	} else {
		fromAccount, ok := accounts[fromNumber]
		if !ok {
			fmt.Fprintf(w, "Account with number %v can't be found!", fromNumber)
			return
		}
		toAccount, ok := accounts[toNumber]
		if !ok {
			fmt.Fprintf(w, "Account with number %v can't be found!", toNumber)
			return
		}

		if err := fromAccount.Transfer(toAccount.Account, amount); err != nil {
			fmt.Fprintf(w, "%v", err)
		} else {
			fmt.Fprint(w, fromAccount.Statement())
			fmt.Fprint(w, toAccount.Statement())
		}
	}
}

func main() {
	accounts[1001] = &JsonAccount{
		Account: &bank.Account{
			Customer: bank.Customer{
				Name:    "Wendy",
				Address: "Los Angeles, California",
				Phone:   "(213) 555 0000",
			},
			Number: 1001,
		},
	}
	accounts[1001] = &JsonAccount{
		Account: &bank.Account{
			Customer: bank.Customer{
				Name:    "John",
				Address: "Los Angeles, California",
				Phone:   "(213) 555 0147",
			},
			Number: 1002,
		},
	}

	http.HandleFunc("/deposit", deposit)
	http.HandleFunc("/withdraw", withdraw)
	http.HandleFunc("/statement", statement)
	http.ListenAndServe("localhost:8000", nil)
}
