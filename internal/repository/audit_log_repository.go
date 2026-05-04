package repository

import (
	"github.com/fkula5/remitly-internship/internal/models"
	"gorm.io/gorm"
)

type AuditLogRepository interface {
	Create(db *gorm.DB, log *models.AuditLog) error
	GetAll(db *gorm.DB) ([]models.AuditLog, error)
}

type gormAuditLogRepository struct{}

func NewAuditLogRepository() AuditLogRepository {
	return &gormAuditLogRepository{}
}

func (r *gormAuditLogRepository) Create(db *gorm.DB, log *models.AuditLog) error {
	return db.Create(log).Error
}

func (r *gormAuditLogRepository) GetAll(db *gorm.DB) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := db.Order("created_at asc").Find(&logs).Error
	return logs, err
}
