package entity

import "time"

type User struct {
	Id         int       `json:"id"`
	IsVerified bool      `json:"is_verified" db:"is_verified"`
	Name       string    `json:"name"`
	Gender     string    `json:"gender"`
	Email      string    `json:"email,omitempty"`
	Password   string    `json:"password,omitempty" db:"password"`
	Age        int       `json:"age" db:"age"`
	Bio        string    `json:"bio" db:"bio"`
	PhotoUrl   string    `json:"photo_url" db:"photo_url"`
	CreatedAt  time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
