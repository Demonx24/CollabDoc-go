package api

import (
	"CollabDoc-go/global"
	"CollabDoc-go/model/database"
	"CollabDoc-go/model/request"
	"CollabDoc-go/model/response"
	"CollabDoc-go/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"path"

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
		response.FailWithMessage("Invalid request", c)
		return
	}
	log.Printf(" 接收到文档回调（status=%d）：%s\n", req.Status, req.Key)
	var doc_version database.DocumentVersion
	var editlog database.DocumentEditLog
	//截取docuuid和版本号
	docuuid, _, err := utils.GetVersionFromDocKey(req.Key)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	//获取文档信息
	doc, err := documentService.GetPublicDoc(docuuid)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	//获取最新版本号
	latestVersion, err := documnet_vService.GetLatestVersionNumber(doc.ID)
	if err != nil {
		response.FailWithMessage("获取当前版本号失败", c)
		return
	}
	//保存编辑人信息
	editlog.UserUUID = doc.OwnerID
	editlog.DocumentUUID = doc.DocUUID
	editlog.CreatedAt = time.Now()
	editlog.VersionNumber = int(latestVersion + 1)
	//保存文档版本信息
	doc_version.VersionNumber = latestVersion + 1
	doc_version.CreatedAt = time.Now()
	doc_version.DocumentID = doc.ID
	doc_version.CreatedBy = doc.OwnerID
	doc_version.VersionName = fmt.Sprintf("%s-V%d", doc.Title, doc_version.VersionNumber)

	if req.Status == 2 {

		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("%s_%s.%s", req.Key, timestamp, req.FileType)
		savePath := path.Join("saved", filename)
		doc_version.FilePath = savePath
		//存入版本信息
		doc_version, err = documnet_vService.Createdoc_v(doc_version)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
		}
		//存入编辑人信息
		err = editlogService.CreateEditLog(&editlog)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
		}
		// 修复 URL 中的 \u0026
		req.URL = strings.ReplaceAll(req.URL, `\u0026`, `&`)

		go func() {
			err := utils.UploadFromURLToMinio(req.URL, savePath)
			if err != nil {
				log.Println(" 保存失败:", err)
			} else {
				log.Println("已保存文档:", savePath)
			}
		}()
	} else if req.Status == 6 {
		//更新版本数据
		latestVersion, err = documnet_vService.GetLatestVersionNumber(doc.ID)
		if err != nil {
			response.FailWithMessage("获取当前版本号失败", c)
			return
		}
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("%s_%s.%s", req.Key, timestamp, req.FileType)
		savePath := path.Join("saved", filename)
		doc_version.FilePath = savePath
		doc_version, err = documnet_vService.Createdoc_v(doc_version)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
		}
		// 修复 URL 中的 \u0026
		req.URL = strings.ReplaceAll(req.URL, `\u0026`, `&`)
		err = editlogService.CreateEditLog(&editlog)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
		}
		go func() {
			err := utils.UploadFromURLToMinio(req.URL, savePath)
			if err != nil {
				log.Println(" 保存失败:", err)
			} else {
				log.Println("已保存文档:", savePath)
			}
		}()
		//更新最新版本
		destDir := path.Join("documents", doc.OwnerID, fmt.Sprintf("doc_%d", doc.ID))
		newFileName := fmt.Sprintf("%s.%s", doc.DocUUID, doc.DocType)
		destPath := path.Join(destDir, newFileName)
		doc.CurrentVersionID = &doc_version.VersionNumber
		err = documentService.UpdatDocument(doc)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return // 这里必须return
		}
		go func() {
			err = utils.UploadFromURLToMinio(req.URL, destPath)
			if err != nil {
				log.Println(" 保存失败:", err)
			} else {
				log.Println("已保存文档:", savePath)
			}
		}()
	}

	c.JSON(http.StatusOK, request.CallbackResponse{Error: 0})
}
func (onlyofficeApi *OnlyofficeApi) GetConfig(c *gin.Context) {
	var req request.OnlyofficeUser
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 模拟从数据库获取文档信息
	var doc database.User_Documents
	if err := global.DB.Where("doc_uuid = ?", req.DocUUID).First(&doc).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if doc.CurrentVersionID == nil {
		response.FailWithMessage("CurrentVersionID is nil", c)
		return
	}
	// 可根据实际情况自定义 key，如 hash(docUUID + versionID)
	docKey := fmt.Sprintf("%s_v%d", doc.DocUUID, *doc.CurrentVersionID)
	//获取文档url
	destDir := path.Join("documents", doc.OwnerID, fmt.Sprintf("doc_%d", doc.ID))
	// 目标文件名，确保唯一
	newFileName := fmt.Sprintf("%s.%s", doc.DocUUID, doc.DocType)
	fileURL := path.Join(destDir, newFileName)
	//fileURL := fmt.Sprintf(
	//	"http://host.docker.internal:8080/api/file?owner_id=%s&doc_uuid=%s&id=%d&ext=%s&_t=%d",
	//	doc.OwnerID, doc.DocUUID, doc.ID, doc.DocType, time.Now().UnixNano(),
	//)
	fmt.Println(fileURL)
	fileURL, err := utils.GetPresignedDownloadURL(fileURL, "")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(fileURL)
	callbackURL := fmt.Sprintf("http://host.docker.internal:8080/api/onlyoffice/callback")

	// 假设你有登录用户信息
	user := database.Useronlyoffice{
		ID:   req.UserUUID, // 实际应为当前登录用户的 ID
		Name: req.Name,     // 实际应为当前用户的名字
	}

	config := database.OnlyOfficeConfig{
		Document: database.Document{
			FileType: doc.DocType,
			Key:      docKey,
			Title:    doc.Title,
			URL:      fileURL,
		},
		DocumentType: doc.DocType,
		EditorConfig: database.EditorConfig{
			AutoSave:    true,
			Mode:        "edit",
			CallbackUrl: callbackURL,
			User: database.Useronlyoffice{
				ID:   user.ID,
				Name: user.Name,
			},
			Permissions: database.Permissions{
				Edit:     true,
				Download: true,
				Print:    true,
				Review:   true,
				Comment:  true,
			},
			Lang: "zh-CN",
			Customization: database.Customization{
				ForceSave:         true,
				Chat:              true,
				Comments:          true,
				CompactHeader:     false,
				Feedback:          true,
				Help:              true,
				ToolbarNoTabs:     false,
				HideRightMenu:     false,
				HideRuler:         false,
				HideToolbar:       false,
				HideFileMenu:      false,
				HideReviewTab:     false,
				ShowReviewChanges: true,
				HideInsertTab:     false,
				HideHomeTab:       false,
				HideViewTab:       false,
			},
			CoEditing: database.CoEditing{
				Mode:   "fast",
				Change: true,
			},
		},
	}

	response.OkWithData(config, c)
}
