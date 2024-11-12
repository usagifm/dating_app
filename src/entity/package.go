package entity

import "time"

type Package struct {
	Id          int       `json:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	Periode     int       `json:"periode" db:"periode"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
