package models

import "time"

// User represents a system user.
type User struct {
	ID          int64      `db:"id"            json:"id"`
	Username    string     `db:"username"      json:"username"`
	Email       string     `db:"email"         json:"email"`
	Password    string     `db:"password"      json:"-"`
	FirstName   string     `db:"first_name"    json:"firstName,omitempty"`
	LastName    string     `db:"last_name"     json:"lastName,omitempty"`
	Country     string     `db:"country"       json:"country,omitempty"`
	DateOfBirth *time.Time `db:"date_of_birth" json:"dateOfBirth,omitempty"`
	CreatedAt   time.Time  `db:"created_at"    json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at"    json:"updatedAt"`
}
