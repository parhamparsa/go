package entity

import "fmt"

type User struct {
	ID        int64  `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
	Active    bool   `db:"active" json:"active"`
}

func (u *User) Validate() error {
	if u.FirstName == "" {
		return fmt.Errorf("name is required")
	}
	if u.LastName == "" {
		return fmt.Errorf("name is required")
	}
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}
