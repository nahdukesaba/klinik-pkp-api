package controller

import (
	"klinik-api/models"
	"klinik-api/service"
	"klinik-api/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ImageController struct {
	Service *service.ImageService
}

func NewImageController(service *service.ImageService) *ImageController {
	return &ImageController{Service: service}
}

// UploadImage handle upload single image
func (c *ImageController) UploadImage(ctx *fiber.Ctx) error {
	// Ambil file dari form
	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("File tidak ditemukan"))
	}

	// Ambil parameter opsional
	entityType := ctx.FormValue("entity_type", "")
	entityIDStr := ctx.FormValue("entity_id", "0")
	entityID, _ := strconv.ParseUint(entityIDStr, 10, 32)

	// Ambil user ID dari context (jika ada auth middleware)
	var uploadedBy uint
	if userID := ctx.Locals("user_id"); userID != nil {
		uploadedBy = userID.(uint)
	}

	// Upload image (simpan metadata ke database)
	image, err := c.Service.UploadImage(file, entityType, uint(entityID), uploadedBy)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(err.Error()))
	}

	// Simpan file ke disk
	if err := c.Service.SaveFile(file, image.Filename); err != nil {
		// Rollback: hapus record dari database jika gagal simpan file
		c.Service.DeleteImage(image.ID)
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Gagal menyimpan file"))
	}

	// Extract dan update dimensi image
	width, height, err := c.Service.ExtractImageDimensions(file)
	if err == nil {
		c.Service.UpdateImageDimensions(image.ID, width, height)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Image berhasil diupload",
		"data":    image,
	})
}

// UploadMultipleImages handle upload multiple images
func (c *ImageController) UploadMultipleImages(ctx *fiber.Ctx) error {
	// Ambil multiple files
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid form data"))
	}

	files := form.File["images"]
	if len(files) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Tidak ada file yang diupload"))
	}

	// Ambil parameter
	entityType := ctx.FormValue("entity_type", "")
	entityIDStr := ctx.FormValue("entity_id", "0")
	entityID, _ := strconv.ParseUint(entityIDStr, 10, 32)

	var uploadedBy uint
	if userID := ctx.Locals("user_id"); userID != nil {
		uploadedBy = userID.(uint)
	}

	var uploadedImages []models.Image
	var errors []string

	// Upload setiap file
	for _, file := range files {
		image, err := c.Service.UploadImage(file, entityType, uint(entityID), uploadedBy)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}

		if err := c.Service.SaveFile(file, image.Filename); err != nil {
			c.Service.DeleteImage(image.ID)
			errors = append(errors, "Gagal menyimpan "+file.Filename)
			continue
		}

		// Extract dimensi
		width, height, err := c.Service.ExtractImageDimensions(file)
		if err == nil {
			c.Service.UpdateImageDimensions(image.ID, width, height)
		}

		uploadedImages = append(uploadedImages, *image)
	}

	response := fiber.Map{
		"success": len(uploadedImages) > 0,
		"message": "Upload selesai",
		"data": fiber.Map{
			"uploaded": uploadedImages,
			"total":    len(uploadedImages),
		},
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

// GetImageByID get image by ID
func (c *ImageController) GetImageByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("ID tidak valid"))
	}

	image, err := c.Service.GetImageByID(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse("Image tidak ditemukan"))
	}

	return ctx.JSON(utils.SuccessMessageResponse("Data berhasil diambil", image))
}

// GetImagesByEntity get images by entity type and ID
func (c *ImageController) GetImagesByEntity(ctx *fiber.Ctx) error {
	entityType := ctx.Query("entity_type")
	entityIDStr := ctx.Query("entity_id")

	if entityType == "" || entityIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("entity_type dan entity_id wajib diisi"))
	}

	entityID, err := strconv.ParseUint(entityIDStr, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("entity_id tidak valid"))
	}

	images, err := c.Service.GetImagesByEntity(entityType, uint(entityID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Gagal mengambil data"))
	}

	return ctx.JSON(utils.SuccessMessageResponse("Data berhasil diambil", images))
}

// GetAllImages get all images with pagination
func (c *ImageController) GetAllImages(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	entityType := ctx.Query("entity_type", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	images, total, err := c.Service.GetAllImages(page, limit, entityType)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Gagal mengambil data"))
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Data berhasil diambil",
		"data":    images,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// DeleteImage delete image by ID
func (c *ImageController) DeleteImage(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("ID tidak valid"))
	}

	if err := c.Service.DeleteImage(uint(id)); err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(err.Error()))
	}

	return ctx.JSON(utils.SuccessMessageResponse("Image berhasil dihapus", nil))
}

// UpdateImageEntity update entity relationship
func (c *ImageController) UpdateImageEntity(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("ID tidak valid"))
	}

	var input struct {
		EntityType string `json:"entity_type"`
		EntityID   uint   `json:"entity_id"`
	}

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	if err := c.Service.UpdateImageEntity(uint(id), input.EntityType, input.EntityID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Gagal update data"))
	}

	return ctx.JSON(utils.SuccessMessageResponse("Data berhasil diupdate", nil))
}
