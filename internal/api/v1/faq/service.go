package faq

import "gorm.io/gorm"

type Service struct {
	db *gorm.DB
}

type FAQPayload struct {
	Question string `json:"question" validate:"required"`
	Answer   string `json:"answer" validate:"required"`
	IsActive *bool  `json:"is_active"`
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetFAQs(page int, limit int, search string) ([]FAQ, int64, error) {
	var faqs []FAQ
	var total int64

	query := s.db.Model(&FAQ{})

	// SEARCH (case-insensitive)
	if search != "" {
		searchQuery := "%" + search + "%"
		query = query.Where("question ILIKE ? OR answer ILIKE ?", searchQuery, searchQuery)
	}

	// COUNT TOTAL (before pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// PAGINATION
	offset := (page - 1) * limit

	if err := query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&faqs).Error; err != nil {
		return nil, 0, err
	}

	return faqs, total, nil
}

func (s *Service) GetFAQById(id uint64) (*FAQ, error) {
	var faq FAQ
	err := s.db.First(&faq, id).Error
	return &faq, err
}

func (s *Service) CreateFAQ(payload *FAQPayload) (*FAQ, error) {
	faq := FAQ{
		Question: payload.Question,
		Answer:   payload.Answer,
	}

	if payload.IsActive != nil {
		faq.IsActive = *payload.IsActive
	}

	if err := s.db.Create(&faq).Error; err != nil {
		return nil, err
	}

	return &faq, nil
}

func (s *Service) UpdateFAQ(id uint64, payload *FAQPayload) error {
	var faq FAQ

	if err := s.db.First(&faq, id).Error; err != nil {
		return err
	}

	update := FAQ{
		Question: payload.Question,
		Answer:   payload.Answer,
	}

	if payload.IsActive != nil {
		update.IsActive = *payload.IsActive
	}

	return s.db.Model(&faq).Updates(update).Error
}

func (s *Service) DeleteFAQ(id uint64) error {
	result := s.db.Delete(&FAQ{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
