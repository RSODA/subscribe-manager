package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/RSODA/subscribe-manager/internal/domain"
	"github.com/google/uuid"
)

const monthYearLayout = "01-2006"

type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name"`
	UserID      uuid.UUID `json:"user_id"`
	Price       int       `json:"price"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string    `json:"service_name"`
	UserID      *uuid.UUID `json:"user_id"`
	Price       *int       `json:"price"`
	StartDate   *string    `json:"start_date"`
	EndDate     *string    `json:"end_date"`
}

type SubscriptionResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Price       int64     `json:"price"`
	ServiceName string    `json:"service_name"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date"`
}

type TotalCostResponse struct {
	TotalCost int `json:"total_cost"`
}

func (r CreateSubscriptionRequest) ToDomain() (*domain.Subscription, error) {
	return toDomainSubscription(r.ServiceName, r.UserID, r.Price, r.StartDate, r.EndDate, uuid.Nil)
}

func (r UpdateSubscriptionRequest) ToDomain(id uuid.UUID) (*domain.Subscription, error) {
	res := &domain.Subscription{ID: id}

	if r.UserID != nil {
		if *r.UserID == uuid.Nil {
			return nil, fmt.Errorf("user_id must be valid uuid")
		}
		res.UserID = *r.UserID
	}

	if r.ServiceName != nil {
		serviceName := strings.TrimSpace(*r.ServiceName)
		if serviceName == "" {
			return nil, fmt.Errorf("service_name cannot be empty")
		}
		res.ServiceName = serviceName
	}

	if r.Price != nil {
		if *r.Price <= 0 {
			return nil, fmt.Errorf("price must be greater than zero")
		}
		res.Price = int64(*r.Price)
	}

	if r.StartDate != nil {
		parsedStartDate, err := parseMonthYear(*r.StartDate)
		if err != nil {
			return nil, fmt.Errorf("start_date must use MM-YYYY format")
		}
		res.StartDate = parsedStartDate
	}

	if r.EndDate != nil {
		parsedEndDate, err := parseOptionalMonthYear(r.EndDate)
		if err != nil {
			return nil, fmt.Errorf("end_date must use MM-YYYY format")
		}
		res.EndDate = parsedEndDate
	}

	if res.UserID == uuid.Nil && res.ServiceName == "" && res.Price == 0 && res.StartDate.IsZero() && res.EndDate == nil {
		return nil, fmt.Errorf("at least one field is required for update")
	}

	return res, nil
}

func NewSubscriptionResponse(sub *domain.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:          sub.ID,
		UserID:      sub.UserID,
		Price:       sub.Price,
		ServiceName: sub.ServiceName,
		StartDate:   formatMonthYear(sub.StartDate),
		EndDate:     formatOptionalMonthYear(sub.EndDate),
	}
}

func NewSubscriptionResponses(subs []*domain.Subscription) []SubscriptionResponse {
	res := make([]SubscriptionResponse, 0, len(subs))
	for _, sub := range subs {
		res = append(res, NewSubscriptionResponse(sub))
	}
	return res
}

func NewTotalCostResponse(total int) TotalCostResponse {
	return TotalCostResponse{TotalCost: total}
}

func toDomainSubscription(serviceName string, userID uuid.UUID, price int, startDate string, endDate *string, id uuid.UUID) (*domain.Subscription, error) {
	serviceName = strings.TrimSpace(serviceName)
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user_id is required")
	}
	if serviceName == "" {
		return nil, fmt.Errorf("service_name is required")
	}
	if price <= 0 {
		return nil, fmt.Errorf("price must be greater than zero")
	}

	parsedStartDate, err := parseMonthYear(startDate)
	if err != nil {
		return nil, fmt.Errorf("start_date must use MM-YYYY format")
	}

	parsedEndDate, err := parseOptionalMonthYear(endDate)
	if err != nil {
		return nil, fmt.Errorf("end_date must use MM-YYYY format")
	}

	return &domain.Subscription{
		ID:          id,
		UserID:      userID,
		Price:       int64(price),
		ServiceName: serviceName,
		StartDate:   parsedStartDate,
		EndDate:     parsedEndDate,
	}, nil
}

func parseMonthYear(value string) (time.Time, error) {
	parsed, err := time.Parse(monthYearLayout, strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}

func parseOptionalMonthYear(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil, nil
	}

	parsed, err := parseMonthYear(trimmed)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func formatMonthYear(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return value.Format(monthYearLayout)
}

func formatOptionalMonthYear(value *time.Time) *string {
	if value == nil || value.IsZero() {
		return nil
	}

	formatted := formatMonthYear(*value)
	return &formatted
}
