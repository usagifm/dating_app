package entity

import "time"

type UserSwipe struct {
	Id        int       `db:"id" json:"id"`
	SwiperId  int       `db:"swiper_id" json:"swiper_id"`
	SwipedId  int       `db:"swiped_id" json:"swiped_id"`
	SwipeType string    `db:"swipe_type" json:"swipe_type"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
