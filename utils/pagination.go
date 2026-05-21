package utils

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	DefaultPageLimit = 10
	MaxPageLimit     = 100
)

type PaginationQuery struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PaginatedData struct {
	Items        any   `json:"items"`
	TotalRecords int64 `json:"total_records"`
	Page         int   `json:"page"`
	Limit        int   `json:"limit"`
}

func (p PaginationQuery) Offset() int {
	return (p.Page - 1) * p.Limit
}

func ParsePaginationQuery(ctx *fiber.Ctx) (PaginationQuery, error) {
	page := 1
	if pageValue := ctx.Query("page"); pageValue != "" {
		parsedPage, err := strconv.Atoi(pageValue)
		if err != nil {
			return PaginationQuery{}, fmt.Errorf("invalid page query parameter")
		}

		page = parsedPage
	}

	if page < 1 {
		return PaginationQuery{}, fmt.Errorf("page must be greater than 0")
	}

	limit := DefaultPageLimit
	if limitValue := ctx.Query("limit"); limitValue != "" {
		parsedLimit, err := strconv.Atoi(limitValue)
		if err != nil {
			return PaginationQuery{}, fmt.Errorf("invalid limit query parameter")
		}

		limit = parsedLimit
	}

	if limit < 1 {
		return PaginationQuery{}, fmt.Errorf("limit must be greater than 0")
	}

	if limit > MaxPageLimit {
		limit = MaxPageLimit
	}

	return PaginationQuery{
		Page:  page,
		Limit: limit,
	}, nil
}

func NewPaginatedData(items any, totalRecords int64, pagination PaginationQuery) PaginatedData {
	return PaginatedData{
		Items:        items,
		TotalRecords: totalRecords,
		Page:         pagination.Page,
		Limit:        pagination.Limit,
	}
}
