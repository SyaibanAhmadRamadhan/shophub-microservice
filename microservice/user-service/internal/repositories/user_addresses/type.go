package useraddresses

import "time"

type Entity struct {
	ID          int64     `db:"id"`
	UserID      int64     `db:"user_id"`
	FullAddress int64     `db:"full_address"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
