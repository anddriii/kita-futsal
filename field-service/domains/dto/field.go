package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type FieldRequest struct {
	Name         string                 `json:"name" validate:"required"`
	Code         string                 `json:"code" validate:"required"`
	PricePerHour int                    `json:"pricePerHour" validate:"required"`
	Images       []multipart.FileHeader `json:"images" validate:"required"`
}

type UpdateFieldRequest struct {
	Name         string                 `json:"name" validate:"required"`
	Code         string                 `json:"code" validate:"required"`
	PricePerHour int                    `json:"pricePerHour" validate:"required"`
	Images       []multipart.FileHeader `json:"images"`
}

type FieldResponse struct {
	UUID         uuid.UUID `json:"uuid"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	PricePerHour any       `json:"pricePerHour"`
	Images       []string  `json:"images"`
	CreatedAt    *time.Time
	UpdateAt     *time.Time
}

type FieldDetailReponse struct {
	Code         string   `json:"code"`
	Name         string   `json:"name"`
	PricePerHour int      `json:"pricePerHour"`
	Images       []string `json:"images"`
	CreatedAt    *time.Time
	UpdateAt     *time.Time
}

type FieldRequestParam struct {
	Page       int     `form:"page" validate:"required"`
	Limit      int     `form:"page" validate:"required"`
	SortColumn *string `form:"sortColumn"`
	SortOrder  *string `form:"sortOrder"`
}
