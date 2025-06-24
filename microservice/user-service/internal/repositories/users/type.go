package users

import "time"

type Entity struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Email       string    `db:"email"`
	PhoneNumber string    `db:"phone_number"`
	Password    string    `db:"password"`
	IsVerified  bool      `db:"is_verified"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
