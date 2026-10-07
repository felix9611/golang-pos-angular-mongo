package services

import (
	"context"
	"golang-fixedasset-mongo-backend/backend/config"
	"golang-fixedasset-mongo-backend/backend/models"
	"golang-fixedasset-mongo-backend/backend/dto"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gin-gonic/gin"
)

func CreateActionRecord(
	actionName string,
	actionMethod string,
	actionFrom string,
	actionData any,
	actionSuccess string,
) (interface{}, error) {
	collection := config.GetCollection("action_records")
	result, err := collection.InsertOne(context.Background(), models.ActionRecords{
		ActionName:   actionName,
		ActionMethod: actionMethod,
		ActionFrom:   actionFrom,
		ActionData:   actionData,
		ActionSuccess: actionSuccess,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func ListActionRecords(pageDto *dto.ListActionRecordReqDto) (interface{}, error) {
	collection := config.GetCollection("action_records")
	skip := (pageDto.Page - 1) * pageDto.Limit
	limit := int64(pageDto.Limit)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{}
	options := options.Find()
	options.SetSkip(int64(skip))
	options.SetLimit(limit)
	options.SetSort(bson.D{{"createdAt", -1}})

	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}

	var actionRecords []models.ActionRecords
	if err = cursor.All(ctx, &actionRecords); err != nil {
		return nil, err
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	return gin.H{"lists": actionRecords, "total": count, "page": pageDto.Page, "limit": pageDto.Limit}, nil
}