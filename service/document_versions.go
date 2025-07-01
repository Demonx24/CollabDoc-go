package service

import (
	"CollabDoc-go/global"
	"CollabDoc-go/model/database"
	"errors"
	"gorm.io/gorm"
	"time"
)

type Document_vService struct {
}

func (document_vService *Document_vService) doc_vgerId(doc_v database.DocumentVersion) (database.DocumentVersion, error) {
	db := global.DB
	err := db.Where("id = ?", doc_v.ID).First(&doc_v).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return database.DocumentVersion{}, errors.New("文章不存在")
		}
		return database.DocumentVersion{}, errors.New("查询id失败")
	}
	return doc_v, nil
}
func (document_vService *Document_vService) Createdoc_v(doc_v database.DocumentVersion) (database.DocumentVersion, error) {
	db := global.DB
	doc_v.CreatedAt = time.Now()
	if err := db.Create(&doc_v).Error; err != nil {
		return database.DocumentVersion{}, err
	}
	return doc_v, nil
}
func (document_vService *Document_vService) Deletedoc_v(doc_v database.DocumentVersion) error {
	db := global.DB
	if err := db.Where("id=? ", doc_v.ID).Find(&doc_v).Error; err != nil {
		return errors.New("查询失败")
	}
	if doc_v.DeletedAt.Valid {
		return errors.New("以被删除")
	}

	if err := db.Model(&doc_v).Where("id=?", doc_v.ID).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}
func (document_vService *Document_vService) Updatedoc_v(doc_v database.DocumentVersion) error {
	db := global.DB
	doc_v.UpdatedAt = time.Now()
	if err := db.Model(&doc_v).Where("id=?", doc_v.ID).Updates(doc_v).Error; err != nil {
		return err
	}
	return nil
}

// 获取最新的版本号
func (document_vService *Document_vService) GetLatestVersionNumber(documentID uint) (uint, error) {
	var version database.DocumentVersion
	err := global.DB.
		Where("document_id = ?", documentID).
		Order("version_number DESC").
		First(&version).Error
	if err != nil {
		return 0, err
	}
	return version.VersionNumber, nil
}
func (document_vService *Document_vService) GetVersionsByDocID(documentID uint) ([]database.DocumentVersion, error) {
	var versions []database.DocumentVersion
	err := global.DB.Where("document_id = ?", documentID).
		Order("version_number desc").
		Limit(3).
		Find(&versions).Error
	return versions, err
}
