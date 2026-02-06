package webhook

import "gorm.io/gorm"

type gormRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) WebhookLogRepository {
	return &gormRepo{db: db}
}
