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
	"log"
	"net/http"
)

type DocumentApi struct{}

func (api *DocumentApi) CreateDocumentByUserUUid(c *gin.Context) {
	var req request.CreateDocRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	doc, err := minioService.CreateDocument(req.UserUUID, req.Title, req.DocType)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(response.CreateDocResponse{Document: *doc}, c)
}
func (api *DocumentApi) GetRecommendDocumentsByUserUUID(c *gin.Context) {
	var req request.UserCard
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("请求参数错误", c)
		return
	}

	// 获取用户自己的文档
	docs, err := documentService.GetUserDocument(req.UUID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 筛选 DocType 为 md 的情况
	hasMarkdown := false
	for _, doc := range docs {
		if doc.DocType == "md" {
			hasMarkdown = true
			break
		}
	}

	// 如果没有 Markdown 文档，不推荐
	if !hasMarkdown {
		response.OkWithData([]response.GetDocResponse{}, c)
		return
	}

	// 获取推荐 Markdown 文档（例如公开的、状态正常的）
	recommendDocs, err := documentService.GetUserMarkdownDocs(req.UUID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var cards []response.GetDocResponse
	for _, d := range recommendDocs {
		card := response.GetDocResponse{
			ID:       d.ID,
			DocUUID:  d.DocUUID,
			Title:    d.Title,
			DocType:  d.DocType,
			Status:   d.Status,
			IsPublic: *d.IsPublic,
			Summary:  truncateDescription(d.Description, 80),
			Updated:  d.UpdatedAt.Format("2006-01-02 15:04"),
		}
		cards = append(cards, card)
	}

	response.OkWithData(cards, c)
}

func (api *DocumentApi) GetDocumentByUserUUId(c *gin.Context) {
	var req request.UserCard
	if err := c.ShouldBindQuery(&req); err != nil {
	}

	docs, err := documentService.GetUserDocument(req.UUID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var cards []response.GetDocResponse
	for _, d := range docs {
		card := response.GetDocResponse{
			ID:       d.ID,
			DocUUID:  d.DocUUID,
			Title:    d.Title,
			DocType:  d.DocType,
			Status:   d.Status,
			IsPublic: *d.IsPublic,
			Summary:  truncateDescription(d.Description, 80),
			Updated:  d.UpdatedAt.Format("2006-01-02 15:04"),
		}
		cards = append(cards, card)
	}
	response.OkWithData(cards, c)
}
func truncateDescription(desc *string, max int) *string {
	if desc == nil || len(*desc) <= max {
		return desc
	}
	s := (*desc)[:max] + "..."
	return &s
}

func (api *DocumentApi) UpdateDocument(c *gin.Context) {
	var req request.UpdateDocument
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	doc := database.User_Documents{
		Title:       req.Title,
		DocType:     req.DocType,
		Status:      req.Status,
		IsPublic:    &req.IsPublic,
		DocUUID:     req.DocUUID,
		Description: req.Description,
	}

	if err := documentService.UpdatDocument(doc); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	doc2, err := documentService.GetDocument(doc)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(doc2, c)
}
func (api *DocumentApi) GetPublicDocuments(c *gin.Context) {
	var docs []database.User_Documents

	if err := global.DB.Where("is_public = ?", true).Find(&docs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败：" + err.Error(),
		})
		return
	}

	response.OkWithData(docs, c)
}

func (api *DocumentApi) GetVersions(c *gin.Context) {
	// 1. 读取 query 参数
	var getVersions request.GetVersions
	if err := c.ShouldBindQuery(&getVersions); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 2. 查询版本记录
	versions, err := documnet_vService.GetVersionsByDocID(getVersions.DocumentID)
	if err != nil {
		response.FailWithMessage("获取版本失败："+err.Error(), c)
		return
	}

	// 3. 把 file_path 转为预签名 URL
	for i := range versions {
		// 假设 versions[i].FilePath 就是 MinIO objectKey
		url, err := utils.GetPresignedDownloadURL(versions[i].FilePath, versions[i].VersionName)
		if err != nil {
			// 这里可以日志一下，继续使用原 path
			log.Printf("生成预签名 URL 失败(%s): %v", versions[i].FilePath, err)
			continue
		}
		versions[i].FilePath = url
	}

	// 4. 返回给前端
	response.OkWithData(versions, c)
}

// 获取文档差异
func (api *DocumentApi) GetDiff(c *gin.Context) {
	var req request.GetDiff
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 查询 Mongo 是否已有缓存
	diff, err := mongoService.GetCachedDocDiff(req.DocUUID, req.FromVer, req.ToVer)
	if err == nil && diff != nil {
		//  有缓存，直接返回
		fmt.Println("正在从mongo中拿数据")
		response.OkWithData(diff, c)
		return
	}
	fmt.Println("没有从mongo中拿数据")
	//  没有缓存，异步推送任务
	task := database.DiffMessage{
		DocUUID:     req.DocUUID,
		FromVersion: req.FromVer,
		ToVersion:   req.ToVer,
	}
	taskJson, _ := json.Marshal(task)
	err = kafkaService.SendKafkaMessage(global.Config.Kafka.DiffTopic, req.DocUUID, string(taskJson))
	if err != nil {
		response.FailWithMessage("任务派发失败", c)
		return
	}

	response.OkWithMessage("正在生成差异，请稍后刷新", c)
}
