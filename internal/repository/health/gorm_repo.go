package health

import "gorm.io/gorm"

type gormRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) HealthRepository {
	return &gormRepo{db: db}
}
