package controller

import (
	"klinik-api/models"
	"klinik-api/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RusunController struct {
	Service *service.RusunService
}

func NewRusunController(service *service.RusunService) *RusunController {
	return &RusunController{Service: service}
}

func (c *RusunController) GetAll(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	provinceID, _ := strconv.Atoi(ctx.Query("province_id", "0"))
	regencyID, _ := strconv.Atoi(ctx.Query("regency_id", "0"))
	status := ctx.Query("status", "")
	search := ctx.Query("search", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	rusun, total, err := c.Service.GetAll(page, limit, uint(provinceID), uint(regencyID), status, search)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data",
			"error":   err.Error(),
		})
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Data berhasil diambil",
		"data":    rusun,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (c *RusunController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	rusun, err := c.Service.GetByID(uint(id))
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "rusun tidak ditemukan" {
			status = fiber.StatusNotFound
		}
		return ctx.Status(status).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Data berhasil diambil",
		"data":    rusun,
	})
}

func (c *RusunController) GetByProvince(ctx *fiber.Ctx) error {
	provinceID, err := strconv.ParseUint(ctx.Params("province_id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Province ID tidak valid",
		})
	}

	rusun, err := c.Service.GetByProvince(uint(provinceID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Data berhasil diambil",
		"data":    rusun,
	})
}

func (c *RusunController) GetMapData(ctx *fiber.Ctx) error {
	provinceID, _ := strconv.Atoi(ctx.Query("province_id", "0"))

	mapData, err := c.Service.GetMapData(uint(provinceID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data peta",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Data peta berhasil diambil",
		"data":    mapData,
	})
}

func (c *RusunController) Create(ctx *fiber.Ctx) error {
	var rusun models.Rusun

	if err := ctx.BodyParser(&rusun); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak valid",
		})
	}

	if err := c.Service.Create(&rusun); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Data berhasil dibuat",
		"data":    rusun,
	})
}

func (c *RusunController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	var rusun models.Rusun
	if err := ctx.BodyParser(&rusun); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak valid",
		})
	}

	if err := c.Service.Update(uint(id), &rusun); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "rusun tidak ditemukan" {
			status = fiber.StatusNotFound
		}
		return ctx.Status(status).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Data berhasil diupdate",
		"data":    rusun,
	})
}

func (c *RusunController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	if err := c.Service.Delete(uint(id)); err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "rusun tidak ditemukan" {
			status = fiber.StatusNotFound
		}
		return ctx.Status(status).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Data berhasil dihapus",
	})
}

func (c *RusunController) GetStatistics(ctx *fiber.Ctx) error {
	provinceID, _ := strconv.Atoi(ctx.Query("province_id", "0"))

	stats, err := c.Service.GetStatistics(uint(provinceID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil statistik",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Statistik berhasil diambil",
		"data":    stats,
	})
}

