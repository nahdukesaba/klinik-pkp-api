package service

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"klinik-api/models"
	"klinik-api/utils"
	"mime/multipart"
	"os"
	"path/filepath"

	_ "golang.org/x/image/webp"
	"gorm.io/gorm"
)

type ImageService struct {
	DB     *gorm.DB
	Config utils.ImageUploadConfig
}

func NewImageService(db *gorm.DB) *ImageService {
	config := utils.DefaultImageConfig()
	// Buat folder upload jika belum ada
	if err := os.MkdirAll(config.UploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}
	return &ImageService{
		DB:     db,
		Config: config,
	}
}

// UploadImage meng-handle upload image dengan best practices
func (s *ImageService) UploadImage(file *multipart.FileHeader, entityType string, entityID, uploadedBy uint) (*models.Image, error) {
	// Validasi file
	if err := utils.ValidateImageFile(file, s.Config); err != nil {
		return nil, err
	}

	// Generate nama file unik
	uniqueFilename := utils.GenerateUniqueFilename(file.Filename)
	filePath := filepath.Join(s.Config.UploadDir, uniqueFilename)

	// Buat record image di database
	image := &models.Image{
		Filename:     uniqueFilename,
		OriginalName: utils.SanitizeFilename(file.Filename),
		FilePath:     filePath,
		FileURL:      fmt.Sprintf("/uploads/images/%s", uniqueFilename),
		MimeType:     file.Header.Get("Content-Type"),
		FileSize:     file.Size,
		EntityType:   entityType,
		EntityID:     entityID,
		UploadedBy:   uploadedBy,
	}

	// Simpan ke database dulu
	if err := s.DB.Create(image).Error; err != nil {
		return nil, fmt.Errorf("gagal menyimpan data image: %v", err)
	}

	return image, nil
}

// SaveFile menyimpan file ke disk dan extract dimensi image
func (s *ImageService) SaveFile(file *multipart.FileHeader, filename string) error {
	filePath := filepath.Join(s.Config.UploadDir, filename)

	// Buka file yang diupload
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("gagal membuka file: %v", err)
	}
	defer src.Close()

	// Buat file tujuan
	dst, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("gagal membuat file: %v", err)
	}
	defer dst.Close()

	// Copy file
	if _, err := dst.ReadFrom(src); err != nil {
		return fmt.Errorf("gagal menyimpan file: %v", err)
	}

	return nil
}

// ExtractImageDimensions mengekstrak width dan height dari image
func (s *ImageService) ExtractImageDimensions(file *multipart.FileHeader) (int, int, error) {
	src, err := file.Open()
	if err != nil {
		return 0, 0, err
	}
	defer src.Close()

	img, _, err := image.DecodeConfig(src)
	if err != nil {
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}

// UpdateImageDimensions update width dan height setelah file disimpan
func (s *ImageService) UpdateImageDimensions(id uint, width, height int) error {
	return s.DB.Model(&models.Image{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"width":  width,
			"height": height,
		}).Error
}

// GetImageByID mendapatkan image berdasarkan ID
func (s *ImageService) GetImageByID(id uint) (*models.Image, error) {
	var image models.Image
	if err := s.DB.First(&image, id).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

// GetImagesByEntity mendapatkan semua images untuk entity tertentu
func (s *ImageService) GetImagesByEntity(entityType string, entityID uint) ([]models.Image, error) {
	var images []models.Image
	err := s.DB.Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Find(&images).Error
	if err != nil {
		return nil, err
	}
	return images, nil
}

// DeleteImage menghapus image (soft delete di database dan hapus file)
func (s *ImageService) DeleteImage(id uint) error {
	var image models.Image
	if err := s.DB.First(&image, id).Error; err != nil {
		return fmt.Errorf("image tidak ditemukan")
	}

	// Hapus file dari disk
	if err := os.Remove(image.FilePath); err != nil && !os.IsNotExist(err) {
		// Log error tapi tetap lanjut hapus dari database
		fmt.Printf("Warning: failed to delete file %s: %v\n", image.FilePath, err)
	}

	// Soft delete dari database
	if err := s.DB.Delete(&image).Error; err != nil {
		return fmt.Errorf("gagal menghapus data image: %v", err)
	}

	return nil
}

// UpdateImageEntity mengupdate entity yang terkait dengan image
func (s *ImageService) UpdateImageEntity(id uint, entityType string, entityID uint) error {
	return s.DB.Model(&models.Image{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"entity_type": entityType,
			"entity_id":   entityID,
		}).Error
}

// GetAllImages mendapatkan semua images dengan pagination
func (s *ImageService) GetAllImages(page, limit int, entityType string) ([]models.Image, int64, error) {
	var images []models.Image
	var total int64

	query := s.DB.Model(&models.Image{})

	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}

	// Hitung total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Ambil data dengan pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&images).Error; err != nil {
		return nil, 0, err
	}

	return images, total, nil
}
