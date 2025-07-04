package utils

import (
	"CollabDoc-go/model/database"
	"archive/zip"
	"encoding/xml"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch" // 这是Go的一个常用文本差异库
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// computeDiffByFileType 根据文件扩展名分支处理
func ComputeDiffByFileType(oldFile, newFile string) ([]diffmatchpatch.Diff, error) {
	ext := strings.ToLower(filepath.Ext(oldFile))
	switch ext {
	case ".md", ".txt":
		return diffTextFiles(oldFile, newFile)
	case ".docx":
		return DiffDocxByFile(oldFile, newFile)
	// 其它格式先返回空或自定义
	default:
		return nil, nil
	}
}

func diffTextFiles(oldFile, newFile string) ([]diffmatchpatch.Diff, error) {
	oldBytes, err := ioutil.ReadFile(oldFile)
	if err != nil {
		return nil, err
	}
	newBytes, err := ioutil.ReadFile(newFile)
	if err != nil {
		return nil, err
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(oldBytes), string(newBytes), false)
	return diffs, nil
}

func DiffDocxByFile(oldFile, newFile string) ([]diffmatchpatch.Diff, error) {
	oldText, err := extractDocxText(oldFile)
	if err != nil {
		return nil, err
	}
	newText, err := extractDocxText(newFile)
	if err != nil {
		return nil, err
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(oldText, newText, false)
	return diffs, nil
}
func extractDocxText(docxPath string) (string, error) {
	// 打开 docx 文件 (zip)
	r, err := zip.OpenReader(docxPath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	// 找到 word/document.xml
	var documentFile *zip.File
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			documentFile = f
			break
		}
	}
	if documentFile == nil {
		return "", fmt.Errorf("document.xml not found in docx")
	}

	// 读取 document.xml
	rc, err := documentFile.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	xmlBytes, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}

	// 解析 XML
	var doc database.DocumentDoc
	if err := xml.Unmarshal(xmlBytes, &doc); err != nil {
		return "", err
	}

	// 拼接所有段落文本
	var sb strings.Builder
	for _, p := range doc.Body.Paragraphs {
		for _, r := range p.Runs {
			for _, t := range r.Texts {
				sb.WriteString(t.Content)
			}
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}
