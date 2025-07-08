package service

import (
	"CollabDoc-go/global"
	"CollabDoc-go/model/database"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type DocumentService struct{}

func (s *DocumentService) GetUserMarkdownDocs(ownerID string) ([]database.User_Documents, error) {
	var docs []database.User_Documents
	err := global.DB.Where("doc_type = ? AND owner_id = ? AND status = ?", "md", ownerID, "active").
		Order("updated_at desc").
		Find(&docs).Error
	return docs, err
}

func (documentService *DocumentService) CreateDocument(userUUID string, title, docType string) (*database.User_Documents, error) {
	doc := database.User_Documents{
		Title:   title,
		OwnerID: userUUID,
		DocType: docType,
		Status:  "active",
		DocUUID: uuid.New().String(),
	}
	// 1. 创建数据库记录（文档元信息）
	if err := global.DB.Create(&doc).Error; err != nil {
		return nil, err
	}

	// 2. 拷贝模板文件生成文档实际文件
	filePath, err := CopyTemplateFile(docType, userUUID, doc.ID, doc.DocUUID)
	if err != nil {
		return nil, err
	}

	// 3. 插入初始版本记录
	version := database.DocumentVersion{
		DocumentID:  doc.ID,
		VersionName: "v1.0",
		FilePath:    filePath,
		CreatedBy:   userUUID,
	}
	if err := global.DB.Create(&version).Error; err != nil {
		return nil, err
	}

	// 4. 更新文档表 CurrentVersionID 字段
	doc.CurrentVersionID = &version.VersionNumber
	if err := global.DB.Save(&doc).Error; err != nil {
		return nil, err
	}

	return &doc, nil
}

// CopyTemplateFile 拷贝模板文件生成新文档文件，返回新文件路径
func CopyTemplateFile(docType string, userUUID string, documentID uint, documentUUID string) (string, error) {
	templateDir := "./templates"
	var templateFile string
	switch docType {
	case "docx":
		templateFile = "template.docx"
	case "pptx":
		templateFile = "template.pptx"
	case "xlsx":
		templateFile = "template.xlsx"
	default:
		return "", fmt.Errorf("unsupported document type: %s", docType)
	}

	srcPath := filepath.Join(templateDir, templateFile)
	// 目标目录按用户+文档ID分目录存储，便于管理
	destDir := filepath.Join("./documents", userUUID, fmt.Sprintf("doc_%d", documentID))
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return "", err
	}

	// 目标文件名，确保唯一
	newFileName := fmt.Sprintf("%s.%s", documentUUID, docType)
	destPath := filepath.Join(destDir, newFileName)

	// 拷贝文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return "", err
	}

	return destPath, nil
}

// 查询用户文档
func (documentService *DocumentService) GetUserDocument(userUUID string) ([]database.User_Documents, error) {
	db := global.DB
	var documents []database.User_Documents
	if err := db.Where("owner_id= ?", userUUID).Find(&documents).Error; err != nil {
		global.Log.Error(err.Error())
		return documents, err
	}
	return documents, nil
}

// 查询用户uuid文档
func (documentService *DocumentService) GetUUIdDocument(docUUid string) (database.User_Documents, error) {
	db := global.DB
	var document database.User_Documents
	if err := db.Where("doc_uuid= ?", docUUid).Find(&document).Error; err != nil {
		global.Log.Error(err.Error())
		return document, err
	}
	return document, nil
}

// GetByID 根据自增主键 ID 查询
func (documentService *DocumentService) GetByID(id uint) (*database.User_Documents, error) {
	var doc database.User_Documents
	db := global.DB
	err := db.First(&doc, id).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// 删除用户文档
// func (documentService *DocumentService) DeleteDocument(userUUID uuid.UUID, documentID uint) error {}
// 更新用户文档
func (documentService *DocumentService) UpdatDocument(document database.User_Documents) error {
	db := global.DB
	document.UpdatedAt = time.Now()
	if err := db.Model(&database.User_Documents{}).
		Where("doc_uuid = ?", document.DocUUID).
		Updates(document).Error; err != nil {
		return err
	}
	return nil
}

func (documentService *DocumentService) GetDocument(documnet database.User_Documents) (database.User_Documents, error) {

	db := global.DB
	err := db.Where("doc_uuid = ?", documnet.DocUUID).First(&documnet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return database.User_Documents{}, errors.New("文章不存在")
		}
		return database.User_Documents{}, errors.New("查询id失败")
	}
	return documnet, nil
}

func (documentService *DocumentService) GetPublicDoc(docUUID string) (database.User_Documents, error) {
	cacheKey := "public_doc:" + docUUID
	ctx := context.Background()
	// 先查缓存
	val, err := global.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var doc database.User_Documents
		if err := json.Unmarshal([]byte(val), &doc); err == nil {
			return doc, nil
		}
		// 解析失败走数据库
	}

	// 缓存未命中或解析失败，从数据库查
	doc, err := documentService.GetUUIdDocument(docUUID)
	if err != nil {
		return doc, err
	}

	// 写缓存（异步也行）
	data, _ := json.Marshal(doc)
	_ = global.Redis.Set(ctx, cacheKey, data, time.Minute*10).Err()

	return doc, nil
}
func (documentService *DocumentService) GetUUIDByIDWithCache(id uint) (string, error) {
	ctx := context.Background()
	cacheKeyIDToUUID := fmt.Sprintf("doc_id:%d", id)

	// 先查 id->uuid 缓存
	if val, err := global.Redis.Get(ctx, cacheKeyIDToUUID).Result(); err == nil {
		return val, nil
	}

	// 缓存没命中，从数据库查 uuid
	var doc database.User_Documents
	err := global.DB.Select("doc_uuid").First(&doc, id).Error
	if err != nil {
		return "", err
	}

	// 写缓存，设置过期时间，比如1小时
	_ = global.Redis.Set(ctx, cacheKeyIDToUUID, doc.DocUUID, time.Hour).Err()

	// 同时写 uuid->id 缓存（方便反向查找）
	cacheKeyUUIDToID := fmt.Sprintf("doc_uuid:%s", doc.DocUUID)
	_ = global.Redis.Set(ctx, cacheKeyUUIDToID, id, time.Hour).Err()

	return doc.DocUUID, nil
}

// 类似反向查 id 的函数
func (documentService *DocumentService) GetIDByUUIDWithCache(uuid string) (uint, error) {
	ctx := context.Background()
	cacheKeyUUIDToID := fmt.Sprintf("doc_uuid:%s", uuid)

	// 查 uuid->id 缓存
	if val, err := global.Redis.Get(ctx, cacheKeyUUIDToID).Result(); err == nil {
		id, err := strconv.ParseUint(val, 10, 64)
		if err == nil {
			return uint(id), nil
		}
	}

	// 缓存没命中，查数据库
	var doc database.User_Documents
	err := global.DB.Select("id").Where("doc_uuid = ?", uuid).First(&doc).Error
	if err != nil {
		return 0, err
	}

	// 写缓存
	cacheKeyIDToUUID := fmt.Sprintf("doc_id:%d", doc.ID)
	_ = global.Redis.Set(ctx, cacheKeyUUIDToID, doc.ID, time.Hour).Err()
	_ = global.Redis.Set(ctx, cacheKeyIDToUUID, uuid, time.Hour).Err()

	return doc.ID, nil
}
