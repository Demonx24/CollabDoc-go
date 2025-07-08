package utils

import (
	"CollabDoc-go/global"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func DownloadAndSaveFromMinio(objectKey, savePath string) error {
	ctx := context.Background()
	bucket := global.Config.Minio.Bucket

	// 从 MinIO 获取对象
	obj, err := global.Minio.GetObject(ctx, bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("从 MinIO 获取对象失败: %v", err)
	}
	defer obj.Close()

	// 创建保存目录
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 创建本地文件
	f, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer f.Close()

	// 拷贝数据流到本地文件
	_, err = io.Copy(f, obj)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

func DownloadFile(url string) (string, error) {
	tmpDir := os.TempDir()

	base := filepath.Base(url)
	if idx := strings.Index(base, "?"); idx != -1 {
		base = base[:idx]
	}
	if base == "" {
		base = "tempfile"
	}

	tmpFilePath := filepath.Join(tmpDir, base)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	outFile, err := os.Create(tmpFilePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return "", err
	}

	return tmpFilePath, nil
}

// 从minio中新建文件，使用的模板，传入文档消息在minio创建一个带目录和文件
func CopyTemplateFileToMinIO(docType, userUUID string, documentID uint, documentUUID string) (string, error) {
	var templateFile string
	switch docType {
	case "docx":
		templateFile = "templates/template.docx"
	case "pptx":
		templateFile = "templates/template.pptx"
	case "xlsx":
		templateFile = "templates/template.xlsx"
	case "md":
		templateFile = "templates/template.md"
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

// GetPresignedDownloadURL 返回带下载文件名的预签名 URL，保留原文件后缀
func GetPresignedDownloadURL(objectKey, filename string) (string, error) {
	ctx := context.Background()
	reqParams := make(url.Values)
	timeNow := time.Now().UnixNano()
	reqParams.Set("_t", fmt.Sprintf("%d", timeNow)) // 防缓存

	// 如果传入了自定义下载文件名
	if filename != "" {
		// 从 objectKey 中提取后缀
		ext := path.Ext(objectKey)
		// 如果自定义文件名不包含后缀，则追加原始后缀
		if ext != "" && path.Ext(filename) == "" {
			filename = fmt.Sprintf("%s%s", filename, ext)
		}
		// 设置 HTTP 响应头 Content-Disposition
		disposition := fmt.Sprintf(`attachment; filename="%s"`, filename)
		reqParams.Set("response-content-disposition", disposition)
	}

	// 生成预签名 URL
	presignedURL, err := global.Minio.PresignedGetObject(
		ctx,
		global.Config.Minio.Bucket,
		objectKey,
		time.Minute*15,
		reqParams,
	)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// 上传文件到minio中
func UploadFromURLToMinio(fileURL, objectKey string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("下载失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("URL请求失败，状态码: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = global.Minio.PutObject(
		context.Background(),
		global.Config.Minio.Bucket,
		objectKey,
		resp.Body,
		resp.ContentLength,
		minio.PutObjectOptions{ContentType: contentType},
	)

	if err != nil {
		return fmt.Errorf("上传到MinIO失败: %v", err)
	}

	log.Printf("✅ 文件已从 URL 上传至 MinIO：%s", objectKey)
	return nil
}
