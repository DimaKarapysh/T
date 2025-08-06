package entity

import (
	"github.com/google/uuid"
	"time"
)

type AuthUser struct {
	ID           uuid.UUID `json:"id" goqu:"omitempty" db:"id" validate:"required"`
	RoleID       uuid.UUID `json:"role_id" goqu:"omitempty" db:"role_id" validate:"required"`
	PasswordHash string    `json:"password_hash" goqu:"omitempty" db:"password_hash" validate:"required,min=5"`
	CreatedAt    time.Time `json:"created_at" goqu:"omitempty" db:"created_at" validate:"required"`
	UpdatedAt    time.Time `json:"updated_at" goqu:"omitempty" db:"updated_at" validate:"required"`
}

type Subscription struct {
	ID          uuid.UUID  `json:"id" db:"id"`                       // Уникальный идентификатор записи
	ServiceName string     `json:"service_name" db:"service_name"`   // Название сервиса
	Price       int        `json:"price" db:"price"`                 // Цена в рублях
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`             // UUID пользователя
	StartDate   time.Time  `json:"start_date" db:"start_date"`       // Начало подписки (месяц и год)
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"` // Конец подписки (опционально)
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`       // Для информации
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`       // Для информации
}

type TotalFilter struct {
	UserID      uuid.UUID
	ServiceName string
	StartFrom   time.Time  // inclusive
	EndTo       *time.Time // inclusive
}
