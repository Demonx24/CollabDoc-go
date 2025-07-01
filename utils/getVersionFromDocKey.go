package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func GetVersionFromDocKey(docKey string) (string, int, error) {
	// 匹配形如 "..._v123" 结尾，提取 docID 和 version
	re := regexp.MustCompile(`^(.*)_v(\d+)$`)
	matches := re.FindStringSubmatch(docKey)
	if len(matches) != 3 {
		return "", 0, fmt.Errorf("格式错误，无法解析 docKey: %s", docKey)
	}
	docID := matches[1]
	versionStr := matches[2]

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return "", 0, fmt.Errorf("版本号转换失败: %v", err)
	}
	return docID, version, nil
}
