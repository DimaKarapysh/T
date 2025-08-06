package v01

import (
	"T/internal/entity"
	"fmt"
	"github.com/google/uuid"
	"time"
)

//	type LoginRequest struct {
//		UserID string `json:"user_id" validate:"required,uuid" example:"2c9c9c9c-9c9c-9c9c-9c9c-9c9c9c9c9c9c"`
//	}
//
//	type LoginResponse struct {
//		AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZDJjYmY0Yy04ZDg3LTRhMjktOGZhYy1iYTgyMGE3Njk0OGMiLCJleHAiOjE3MDA4ODAwMDB9.XYZ"`
//		RefreshToken string `json:"refresh_token" example:"dGhpc19pcyBhIHJlZnJlc2ggdG9rZW4gdGVzdA=="`
//	}
//
//	type RefreshTokenRequest struct {
//		AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9..."`
//		RefreshToken string `json:"refresh_token" example:"c29tZS1iYXNlNjQtZW5jb2RlZC10b2tlbg=="`
//	}
//
//	type RefreshTokenResponse struct {
//		AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9..."`
//		RefreshToken string `json:"refresh_token" example:"YW5vdGhlci1uZXctYmFzZTY0LXRva2Vu"`
//	}
type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" example:"Yandex Plus"`
	Price       int    `json:"price" example:"400"`
	UserID      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `json:"start_date" example:"2025-07-01"`
	EndDate     string `json:"end_date,omitempty" example:"2025-12-01"`
}

type SubscriptionResponse struct {
	ID          string `json:"id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

type TotalCostRequest struct {
	UserID      string `query:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string `query:"service_name" example:"Netflix"`
	StartFrom   string `query:"start_from" example:"2024-06-01"`
	EndTo       string `query:"end_to" example:"2025-08-01"`
}

type TotalCostResponse struct {
	Total int `json:"total" example:"1200"`
}

type ListQueryParams struct {
	Limit  *int `query:"limit" example:"20"`
	Offset *int `query:"offset" example:"40"`
}

func (r CreateSubscriptionRequest) ToEntity() (*entity.Subscription, error) {
	const layout = "2006-01-02" // MM-YYYY

	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	startDate, err := time.Parse(layout, r.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}

	var endDate time.Time
	if r.EndDate != "" {
		endDate, err = time.Parse(layout, r.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date: %w", err)
		}
	}

	return &entity.Subscription{
		ID:          uuid.New(),
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     &endDate,
	}, nil
}

type UpdateSubscriptionRequest struct {
	ServiceName string `json:"service_name" example:"Yandex Plus"`
	Price       int    `json:"price" example:"400"`
	UserID      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `json:"start_date" example:"2025-07-01"`
	EndDate     string `json:"end_date,omitempty" example:"2025-12-01"`
}

func (r UpdateSubscriptionRequest) ToEntity() (*entity.Subscription, error) {
	const layout = "2006-01-02"

	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	startDate, err := time.Parse(layout, r.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}

	var endDate time.Time
	if r.EndDate != "" {
		endDate, err = time.Parse(layout, r.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date: %w", err)
		}
	}

	return &entity.Subscription{
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     &endDate,
	}, nil
}

type TotalRequest struct {
	UserID      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string `json:"service_name" example:"Netflix"`
	StartFrom   string `json:"start_from" example:"2024-06-01"`
	EndTo       string `json:"end_to,omitempty" example:"2025-08-01"`
}

func (r TotalRequest) ToFilter() (*entity.TotalFilter, error) {
	const layout = "2006-01-02"

	var (
		userID uuid.UUID
		err    error
		from   time.Time
		toPtr  *time.Time
	)

	if r.UserID != "" {
		userID, err = uuid.Parse(r.UserID)
		if err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
	}

	if r.StartFrom != "" {
		from, err = time.Parse(layout, r.StartFrom)
		if err != nil {
			return nil, fmt.Errorf("invalid from date: %w", err)
		}
	}

	if r.EndTo != "" {
		to, err := time.Parse(layout, r.EndTo)
		if err != nil {
			return nil, fmt.Errorf("invalid to date: %w", err)
		}
		toPtr = &to
	}

	return &entity.TotalFilter{
		UserID:      userID,
		ServiceName: r.ServiceName,
		StartFrom:   from,
		EndTo:       toPtr,
	}, nil
}
