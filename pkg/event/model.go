package event

import (
	"time"
)

type Event struct {
	ID                 string    `json:"id"`
	CookieID           string    `json:"cid" query:"cid" validate:"required"`
	Country            string    `json:"c" query:"c" validate:"max=3, min=1"`
	Email              string    `json:"e" query:"e" validate:"email"`
	Hotel              string    `json:"h" query:"h"`
	ConfirmationNumber string    `json:"cf" query:"cf"`
	ExtraField         string    `json:"ex"`
	CreatedAt          time.Time `json:"createdAt"`
}
