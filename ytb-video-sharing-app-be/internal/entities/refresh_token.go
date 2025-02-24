package entities

import "time"

type RefreshToken struct {
	ID        int64     `db:"id"`
	AccountID int64     `db:"account_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
