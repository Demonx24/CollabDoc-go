package service

import (
	"CollabDoc-go/global"
	"CollabDoc-go/model/database"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
