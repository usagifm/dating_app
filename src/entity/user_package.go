package entity

import "time"

type UserPackage struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id" db:"user_id"`
	PackageId int       `json:"package_id" db:"package_id"`
	ValidDate time.Time `json:"valid_date" db:"valid_date"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
