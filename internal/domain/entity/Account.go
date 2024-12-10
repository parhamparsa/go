package entity

import (
	"errors"
	"fmt"
)

type AccountType string

const (
	AccountTypeCredit AccountType = "credit"
	AccountTypeNormal AccountType = "normal"
)

type Account struct {
	ID      int64       `db:"id" json:"id"`
	UserID  int64       `db:"user_id" json:"user_id"`
	Type    AccountType `db:"type" json:"type"`
	Balance int64       `db:"balance" json:"balance"`
}

func (a *Account) Validate() error {
	if a.UserID <= 0 {
		return errors.New("user_id must be greater than zero")
	}

	if !a.IsValidType() {
		return fmt.Errorf("invalid account type: %s", a.Type)
	}

	if a.Balance < 0 {
		return errors.New("balance cannot be negative")
	}

	return nil
}

func (a *Account) IsValidType() bool {
	switch a.Type {
	case AccountTypeCredit, AccountTypeNormal:
		return true
	default:
		return false
	}
}
