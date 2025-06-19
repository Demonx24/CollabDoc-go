package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadAndSaveFile(url, savePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("下载失败: %v", err)
	}
	defer resp.Body.Close()

	// 创建目录
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	f, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
