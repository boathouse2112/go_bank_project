package bank

import (
	"errors"
	"fmt"
)

type Customer struct {
	Name    string
	Address string
	Phone   string
}

type Account struct {
	Customer
	Number  int64
	Balance float64
}

type Bank interface {
	Statement() string
}

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("the amount to deposit should be greater than 0")
	}

	a.Balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("the amount to withdraw should be greater than 0")
	}

	if a.Balance < amount {
		return errors.New("the amount to withdraw should not be greater than the balance of the account")
	}

	a.Balance -= amount
	return nil
}

func (a *Account) Transfer(receiver *Account, amount float64) error {

	if amount > a.Balance {
		return errors.New("the amount to transfer should not be greater than the balance of the account")
	}

	if err := a.Withdraw(amount); err != nil {
		return err
	}
	if err := receiver.Deposit(amount); err != nil {
		return err
	}
	return nil
}

func (a *Account) Statement() string {
	return fmt.Sprintf("%v - %v - %v", a.Number, a.Customer.Name, a.Balance)
}

func Statement(b Bank) string {
	return b.Statement()
}
