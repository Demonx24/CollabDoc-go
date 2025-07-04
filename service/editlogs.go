package service

import (
	"CollabDoc-go/global"
	"CollabDoc-go/model/database"
	"errors"
	"gorm.io/gorm"
	"time"
)

type EditLogService struct{}

// CreateEditLog 创建编辑日志记录
func (s *EditLogService) CreateEditLog(log *database.DocumentEditLog) error {
	log.CreatedAt = time.Now()
	return global.DB.Create(log).Error
}

// GetLogsByDocUUID 根据文档 UUID 获取所有编辑记录（按时间倒序）
func (s *EditLogService) GetLogsByDocUUID(docUUID string) ([]database.DocumentEditLog, error) {
	var logs []database.DocumentEditLog
	err := global.DB.
		Where("document_uuid = ?", docUUID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// GetLogsByUserUUID 获取某个用户所有编辑记录
func (s *EditLogService) GetLogsByUserUUID(userUUID string) ([]database.DocumentEditLog, error) {
	var logs []database.DocumentEditLog
	err := global.DB.
		Where("user_uuid = ?", userUUID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// DeleteLogByID 根据 ID 删除日志
func (s *EditLogService) DeleteLogByID(id int64) error {
	return global.DB.Delete(&database.DocumentEditLog{}, id).Error
}

// GetLogByID 获取单条日志记录
func (s *EditLogService) GetLogByID(id int64) (*database.DocumentEditLog, error) {
	var log database.DocumentEditLog
	err := global.DB.First(&log, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &log, err
}
