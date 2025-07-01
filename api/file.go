package api

import (
	"CollabDoc-go/model/request"
	"CollabDoc-go/model/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
)

type FileApi struct {
}

func (FileApi *FileApi) FilePath(c *gin.Context) {
	var req request.FilePath
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	baseDir := "E:\\go代码\\CollabDoc-go\\documents"
	filePath := filepath.Join(
		baseDir,
		req.OwnerID,
		fmt.Sprintf("doc_%d", req.ID),
		req.DocUUID+"."+req.Ext,
	)
	fmt.Println("🧾 拼接路径：", filePath)
	fmt.Printf("📥 请求参数：OwnerID=%s, DocUUID=%s, ID=%d, Ext=%s\n", req.OwnerID, req.DocUUID, req.ID, req.Ext)

	// 返回文件，如果找不到会自动 404
	c.File(filePath)
}
