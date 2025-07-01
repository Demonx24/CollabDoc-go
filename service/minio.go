package service

import (
	"CollabDoc-go/global"
	"CollabDoc-go/model/database"
	"CollabDoc-go/utils"
	"fmt"
	"github.com/google/uuid"
)

type MinioService struct {
}

func (minioService *MinioService) CreateDocument(userUUID string, title, docType string) (*database.User_Documents, error) {
	// 1. 初始化文档记录（元信息）
	doc := database.User_Documents{
		Title:   title,
		OwnerID: userUUID,
		DocType: docType,
		Status:  "active",
		DocUUID: uuid.New().String(),
	}

	if err := global.DB.Create(&doc).Error; err != nil {
		return nil, err
	}

	// 2. 从 MinIO 模板复制为初始文件（返回对象路径）
	objectKey, err := utils.CopyTemplateFileToMinIO(docType, userUUID, doc.ID, doc.DocUUID)
	if err != nil {
		return nil, fmt.Errorf("模板文件复制失败: %w", err)
	}

	// 3. 插入文档版本信息（记录对象路径）
	version := database.DocumentVersion{
		DocumentID:  doc.ID,
		VersionName: "v1.0",
		FilePath:    objectKey, // 这里是 MinIO 中的 Key
		CreatedBy:   userUUID,
	}

	if err := global.DB.Create(&version).Error; err != nil {
		return nil, err
	}

	// 4. 更新主文档表中的当前版本号
	doc.CurrentVersionID = &version.VersionNumber
	if err := global.DB.Save(&doc).Error; err != nil {
		return nil, err
	}

	return &doc, nil
}
