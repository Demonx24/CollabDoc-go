package initialize

import (
	"CollabDoc-go/global"
	"CollabDoc-go/model/database"
	"CollabDoc-go/service"
	"CollabDoc-go/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sergi/go-diff/diffmatchpatch"
	"go.uber.org/zap"
)

// DiffMessage 代表 diff_topic 中的消息结构

// StartDiffConsumer 启动一个 Kafka 消费协程，监听文档差异主题
func StartDiffConsumer(ctx context.Context) {
	kafkaCfg := global.Config.Kafka
	if !kafkaCfg.Enabled {
		global.Log.Warn("Kafka 消费端未启用，跳过初始化")
		return
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	client, err := sarama.NewConsumerGroup(kafkaCfg.Brokers, kafkaCfg.ConsumerGroup, config)
	if err != nil {
		global.Log.Fatal("Kafka 消费者组创建失败", zap.Error(err))
	}

	handler := &diffConsumerGroupHandler{}

	go func() {
		for {
			if err := client.Consume(ctx, []string{kafkaCfg.DiffTopic}, handler); err != nil {
				global.Log.Error("Kafka 消费失败", zap.Error(err))
			}
		}
	}()
	global.Log.Info("Kafka Diff 消费者已启动", zap.String("topic", kafkaCfg.DiffTopic))
}

type diffConsumerGroupHandler struct{}

func (h *diffConsumerGroupHandler) Setup(s sarama.ConsumerGroupSession) error   { return nil }
func (h *diffConsumerGroupHandler) Cleanup(s sarama.ConsumerGroupSession) error { return nil }

func (h *diffConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var diff database.DiffMessage
		if err := json.Unmarshal(msg.Value, &diff); err != nil {
			global.Log.Warn("解析 Kafka 消息失败", zap.ByteString("msg", msg.Value), zap.Error(err))
			continue
		}
		var docdiff database.DocDiff
		doc, err := service.ServiceGroupApp.DocumentService.GetUUIdDocument(diff.DocUUID)
		if err != nil {
			global.Log.Error("根据 UUID 查询文档失败", zap.Error(err))
			continue
		}

		documents, err := service.ServiceGroupApp.Document_vService.GetVersionsByDocID(doc.ID)
		if err != nil {
			global.Log.Error("查询文档版本失败", zap.Error(err))
			continue
		}

		if len(documents) < 2 {
			global.Log.Warn("文档版本不足，无法计算差异", zap.Int("version_count", len(documents)))
			continue
		}
		docdiff.FromVersion = diff.FromVersion
		docdiff.ToVersion = diff.ToVersion
		NewVersionUrl, _ := utils.GetPresignedDownloadURL(documents[0].FilePath, "")
		OldVersionUrl, _ := utils.GetPresignedDownloadURL(documents[1].FilePath, "")
		oldFilePath, err1 := utils.DownloadFile(OldVersionUrl)
		newFilePath, err2 := utils.DownloadFile(NewVersionUrl)
		if err1 != nil || err2 != nil {
			global.Log.Error("下载文档失败", zap.Error(err1), zap.Error(err2))
			continue
		}
		//存入差别信息
		diffs, err := utils.ComputeDiffByFileType(oldFilePath, newFilePath)
		if err != nil {
			global.Log.Error("计算差异失败", zap.Error(err))
			continue
		}
		doc, err = service.ServiceGroupApp.DocumentService.GetUUIdDocument(diff.DocUUID)
		if err != nil {
			global.Log.Error("根据 UUID 查询文档失败", zap.Error(err))
			continue
		}

		documents, err = service.ServiceGroupApp.Document_vService.GetVersionsByDocID(doc.ID)
		if err != nil {
			global.Log.Error("查询文档版本失败", zap.Error(err))
			continue
		}

		if len(documents) < 2 {
			global.Log.Warn("文档版本不足，无法计算差异", zap.Int("version_count", len(documents)))
			continue
		}
		changedFields := ConvertDiffsToMap(diffs)
		docdiff.DocUUID = diff.DocUUID

		docdiff.ChangedFields = changedFields

		ctx := context.Background()
		if err := service.ServiceGroupApp.MongoService.SaveDocDiff(ctx, &docdiff); err != nil {
			global.Log.Error("保存文档差异失败", zap.Error(err))
		}

		//  示例处理：打印/记录/后续保存到 Mongo
		fmt.Printf("收到文档更新: %+v\n", diff)
		global.Log.Info("接收到文档 diff", zap.Any("diff", diff))

		//  在这里调用存入 Mongo 的逻辑（可自定义为 service 层调用）

		session.MarkMessage(msg, "")
	}
	return nil
}
func ConvertDiffsToMap(diffs []diffmatchpatch.Diff) map[string]interface{} {
	diffItems := make([]database.DiffItem, 0, len(diffs))
	for _, d := range diffs {
		var op string
		switch d.Type {
		case diffmatchpatch.DiffInsert:
			op = "INSERT"
		case diffmatchpatch.DiffDelete:
			op = "DELETE"
		case diffmatchpatch.DiffEqual:
			op = "EQUAL"
		}
		diffItems = append(diffItems, database.DiffItem{
			Operation: op,
			Text:      d.Text,
		})
	}
	return map[string]interface{}{"diffs": diffItems}
}
