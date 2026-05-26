package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `json:"id" example:"14ef01ac-0a66-4f3f-bf59-8f75bc50b6f8"`
	UserID      uuid.UUID  `json:"user_id" example:"0e7a1e5f-94f0-4968-838a-e0329b0d556e"`
	Price       int64      `json:"price" example:"399"`
	ServiceName string     `json:"service_name" example:"Telegram Premium"`
	StartDate   time.Time  `json:"start_date" example:"2026-07-01T00:00:00Z"`
	EndDate     *time.Time `json:"end_date" example:"2026-12-01T00:00:00Z"`
}
