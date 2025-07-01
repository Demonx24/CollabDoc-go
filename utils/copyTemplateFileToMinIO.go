package utils

import (
	"CollabDoc-go/global"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
)

func CopyTemplateFileToMinIO(docType, userUUID string, documentID uint, documentUUID string) (string, error) {
	var templateFile string
	switch docType {
	case "docx":
		templateFile = "templates/template.docx"
	case "pptx":
		templateFile = "templates/template.pptx"
	case "xlsx":
		templateFile = "templates/template.xlsx"
	default:
		return "", fmt.Errorf("unsupported document type: %s", docType)
	}

	// 构造目标路径
	objectKey := fmt.Sprintf("documents/%s/doc_%d/%s.%s", userUUID, documentID, documentUUID, docType)

	// 构造复制源/目标
	src := minio.CopySrcOptions{
		Bucket: global.Config.Minio.Bucket,
		Object: templateFile,
	}
	dst := minio.CopyDestOptions{
		Bucket: global.Config.Minio.Bucket,
		Object: objectKey,
	}

	// 执行复制操作
	ctx := context.Background()
	_, err := global.Minio.CopyObject(ctx, dst, src)
	if err != nil {
		return "", fmt.Errorf("failed to copy template to %s: %w", objectKey, err)
	}

	return objectKey, nil
}
