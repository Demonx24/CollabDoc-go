package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var MongoClient *mongo.Client
var MongoColl *mongo.Collection

func InitMongo() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	MongoClient = client
	MongoColl = client.Database("collabdoc").Collection("docs")
}

func SaveDoc(name string) string {
	res, err := MongoColl.InsertOne(context.TODO(), map[string]string{"name": name})
	if err != nil {
		log.Printf("Insert error: %v", err)
		return ""
	}
	if res.InsertedID == nil {
		log.Println("InsertedID is nil")
		return ""
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println("InsertedID type assert failed")
		return ""
	}
	return oid.Hex()
}
func LoadDoc(id string) map[string]interface{} {
	objID, _ := primitive.ObjectIDFromHex(id)
	var doc map[string]interface{}
	MongoColl.FindOne(context.TODO(), map[string]interface{}{"_id": objID}).Decode(&doc)
	return doc
}
