package internal

import "time"

type User struct {
	PK        int
	Name      string
	Email     string
	Address   string
	Cpf       string
	Phone     string
	Birthday  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
