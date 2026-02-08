package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type FieldRequest struct {
	Name         string                 `form:"name" validate:"required"`
	Code         string                 `form:"code" validate:"required"`
	Latitude     float64                `form:"latitude"`
	Lonitude     float64                `form:"lonitude"`
	PricePerHour int                    `form:"pricePerHour" validate:"required"`
	Images       []multipart.FileHeader `form:"images" validate:"required"`
}

type UpdateFieldRequest struct {
	Name         string                 `form:"name" validate:"required"`
	Code         string                 `form:"code" validate:"required"`
	Latitude     float64                `form:"latitude"`
	Lonitude     float64                `form:"lonitude"`
	PricePerHour int                    `form:"pricePerHour" validate:"required"`
	Images       []multipart.FileHeader `form:"images"`
}

type FieldResponse struct {
	UUID         uuid.UUID `json:"uuid"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Latitude     float64   `form:"latitude"`
	Lonitude     float64   `form:"lonitude"`
	PricePerHour any       `json:"pricePerHour"`
	Images       []string  `json:"images"`
	Distance     float64   `json:"distance"`
	CreatedAt    *time.Time
	UpdateAt     *time.Time
}

type FieldDetailReponse struct {
	Code         string   `json:"code"`
	Name         string   `json:"name"`
	PricePerHour int      `json:"pricePerHour"`
	Images       []string `json:"images"`
	Latitude     float64  `form:"latitude"`
	Lonitude     float64  `form:"lonitude"`
	CreatedAt    *time.Time
	UpdateAt     *time.Time
}

type FieldRequestParam struct {
	Page       int     `form:"page" validate:"required"`
	Limit      int     `form:"limit" validate:"required"`
	SortColumn *string `form:"sortColumn"`
	SortOrder  *string `form:"sortOrder"`
}

type NearbyFields struct {
	Latitude float64 `form:"lat" validate:"required"`
	Lonitude float64 `form:"lon" validate:"required"`
}
