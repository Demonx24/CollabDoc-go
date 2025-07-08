package service

import (
	"CollabDoc-go/global"
	"CollabDoc-go/model/database"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoService struct{}

func (mongoService *MongoService) SaveDocDiff(ctx context.Context, diff *database.DocDiff) error {
	coll := global.Mongo.
		Database(global.Config.Mongo.Database).
		Collection("doc_diff")

	filter := bson.M{
		"doc_uuid":     diff.DocUUID,
		"from_version": diff.FromVersion,
		"to_version":   diff.ToVersion,
	}

	update := bson.M{
		"$set": diff,
	}

	opts := options.Update().SetUpsert(true)

	_, err := coll.UpdateOne(ctx, filter, update, opts)
	return err
}
func (mongoService *MongoService) GetCachedDocDiff(docUUID string, from, to int) (*database.DocDiff, error) {
	collection := global.Mongo.
		Database(global.Config.Mongo.Database).
		Collection("doc_diff")

	filter := bson.M{
		"doc_uuid":     docUUID,
		"from_version": from,
		"to_version":   to,
	}

	var result database.DocDiff
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SaveDocumentMd 保存markdown文档最新文本、版本和操作历史
func SaveDocumentMd(docId string, text string, version int, op database.Operation) error {
	collection := global.Mongo.
		Database(global.Config.Mongo.Database).
		Collection("Md")
	filter := bson.M{"doc_id": docId}
	update := bson.M{
		"$set": bson.M{
			"text":       text,
			"version":    version,
			"updated_at": time.Now(),
		},
		"$push": bson.M{
			"history": op,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Printf("SaveDocumentMd Mongo UpdateOne error: %v\n", err)
		return err
	}
	return nil
}
func GetDocumentMd(docId string) (*database.DocumentMd, error) {
	collection := global.Mongo.
		Database(global.Config.Mongo.Database).
		Collection("Md")

	filter := bson.M{"doc_id": docId}

	var doc database.DocumentMd
	err := collection.FindOne(context.Background(), filter).Decode(&doc)
	if err != nil {
		log.Printf("GetDocumentMd Mongo FindOne error: %v\n", err)
		return nil, err
	}

	return &doc, nil
}
