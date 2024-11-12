package entity

import "time"

type UserPreference struct {
	Id              int       `db:"id" json:"id"`
	UserId          int       `db:"user_id" json:"user_id"`
	PreferredGender string    `db:"preferred_gender" json:"preferred_gender"`
	MinAge          int       `db:"min_age" json:"min_age"`
	MaxAge          int       `db:"max_age" json:"max_age"`
	CreatedAt       time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
