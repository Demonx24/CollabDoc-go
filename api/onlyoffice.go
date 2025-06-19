package api

import (
	"CollabDoc-go/model/request"
	"CollabDoc-go/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type OnlyofficeApi struct{}

func (onlyofficeApi *OnlyofficeApi) Callback(c *gin.Context) {
	var req request.CallbackRequest
	body, _ := io.ReadAll(c.Request.Body)

	if err := json.Unmarshal(body, &req); err != nil {
		log.Println("JSON decode error:", err)
		log.Println("Raw request body:", string(body))
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	log.Printf("✅ 接收到文档回调（status=%d）：%s\n", req.Status, req.Key)

	if req.Status == 2 || req.Status == 6 {
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("%s_%s.%s", req.Key, timestamp, req.FileType)
		savePath := filepath.Join("saved", filename)

		// 修复 URL 中的 \u0026
		req.URL = strings.ReplaceAll(req.URL, `\u0026`, `&`)

		go func() {
			err := utils.DownloadAndSaveFile(req.URL, savePath)
			if err != nil {
				log.Println("❌ 保存失败:", err)
			} else {
				log.Println("✅ 已保存文档:", savePath)
			}
		}()
	}

	c.JSON(http.StatusOK, request.CallbackResponse{Error: 0})
}
