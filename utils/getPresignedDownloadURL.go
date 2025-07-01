package utils

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"time"

	"CollabDoc-go/global"
)

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
