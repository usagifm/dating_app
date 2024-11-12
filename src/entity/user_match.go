package entity

import "time"

type UserMatch struct {
	Id        int       `db:"id" json:"id"`
	UserId1   int       `db:"user_id_1" json:"user_id_1"`
	UserId2   int       `db:"user_id_2" json:"user_id_2"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
