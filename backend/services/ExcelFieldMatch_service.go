package services

import (
	"context"
	"golang-pos-mongo-backend/backend/config"
	"golang-pos-mongo-backend/backend/dto"
	"golang-pos-mongo-backend/backend/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateExcelFieldMatch(data *models.ExcelFieldMatchs) (interface{}, error) {
	collection := config.GetCollection("excel_field_matchs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	filter := bson.M{"status": 1, "functionCode": data.FunctionCode, "functionName": data.FunctionName, "functionType": data.FunctionType}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return gin.H{"error": "Failed to create Excel Field Match"}, err
	}

	if count == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if data.Status == 0 {
			data.Status = 1
		}

		if data.CreatedAt.IsZero() {
			data.CreatedAt = time.Now()
		}

		if data.UpdatedAt.IsZero() {
			data.UpdatedAt = time.Now()
		}

		result, err := collection.InsertOne(ctx, data)
		if err != nil {
			return gin.H{"error": "Failed to create Excel Field Match"}, err
		}

		CreateActionRecord("Excel Field Match Create", "POST", "Excel Field Match", data, "Success")

		return result, nil
	} else {
		CreateActionRecord("Excel Field Match Create", "POST", "Excel Field Match", data, "Failed")
		return gin.H{"error": "Excel Field Match already exists"}, nil
	}
}

func GetOneExcelFieldMatch(id string) (interface{}, error) {
	collection := config.GetCollection("excel_field_matchs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var excelFieldMatch models.ExcelFieldMatchs
	err = collection.FindOne(ctx, filter).Decode(&excelFieldMatch)
	if err != nil {
		return nil, err
	}

	if excelFieldMatch.Status == 0 {
		return gin.H{"error": "Excel Field Match not found or already voided"}, nil
	} else {
		return excelFieldMatch, nil
	}

}

func GetOneExcelFieldMatchByCode(code string) (interface{}, error) {
	collection := config.GetCollection("excel_field_matchs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"functionCode": code}

	var excelFieldMatch models.ExcelFieldMatchs
	err := collection.FindOne(ctx, filter).Decode(&excelFieldMatch)
	if err != nil {
		return nil, err
	}

	if excelFieldMatch.Status == 0 {
		return gin.H{"error": "Excel Field Match not found or already voided"}, nil
	} else {
		return excelFieldMatch, nil
	}
}

func InactiveOneExcelFieldMatch(id string) (interface{}, error) {

	collection := config.GetCollection("excel_field_matchs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var excelFieldMatch models.ExcelFieldMatchs
	err = collection.FindOne(ctx, filter).Decode(&excelFieldMatch)
	if err != nil {
		return nil, err
	}

	if excelFieldMatch.Status == 0 {
		CreateActionRecord("Excel Field Match Inactive", "DELETE", "Excel Field Match", excelFieldMatch, "Failed")
		return gin.H{"error": "Excel Field Match not found or already voided"}, nil
	} else {
		excelFieldMatch.Status = 0
		excelFieldMatch.UpdatedAt = time.Now()

		result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": excelFieldMatch})
		if err != nil {
			return nil, err
		}
		CreateActionRecord("Excel Field Match Inactive", "DELETE", "Excel Field Match", excelFieldMatch, "Success")
		return result, nil
	}
}

func UpdateOneExcelFieldMatch(data *models.ExcelFieldMatchs) (interface{}, error) {

	collection := config.GetCollection("excel_field_matchs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": data.ID}

	var excelFieldMatch models.ExcelFieldMatchs
	err := collection.FindOne(ctx, filter).Decode(&excelFieldMatch)
	if err != nil {
		return nil, err
	}

	if excelFieldMatch.Status == 0 {
		CreateActionRecord("Excel Field Match Update", "POST", "Excel Field Match", data, "Failed")
		return gin.H{"error": "Excel Field Match not found or already voided"}, nil
	} else {
		data.UpdatedAt = time.Now()
		result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": data})
		if err != nil {
			return nil, err
		}
		CreateActionRecord("Excel Field Match Update", "POST", "Excel Field Match", data, "Success")
		return result, nil
	}
}

func ExcelFieldMatchListAndPage(
	pageDto dto.ExcelFieldMatchPageDto,
) (interface{}, error) {

	collection := config.GetCollection("excel_field_matchs")

	if pageDto.Page <= 0 {
		pageDto.Page = 1
	}

	if pageDto.Limit <= 0 {
		pageDto.Limit = 10
	}

	skip := (pageDto.Page - 1) * pageDto.Limit
	limit := int64(pageDto.Limit)

	log.Printf(
		"ExcelFieldMatch Page=%d Limit=%d Skip=%d",
		pageDto.Page,
		pageDto.Limit,
		skip,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1}

	if pageDto.Name != "" {
		filter["FunctionNam"] = bson.M{
			"$regex":   pageDto.Name,
			"$options": "i",
		}
	}

	if pageDto.Type != "" {
		filter["functionType"] = bson.M{
			"$regex":   pageDto.Type,
			"$options": "i",
		}
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{
		{Key: "createdAt", Value: -1},
		{Key: "_id", Value: -1},
	})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	records := make([]models.ExcelFieldMatchs, 0)

	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"lists": records,
		"total": count,
		"page":  pageDto.Page,
		"limit": pageDto.Limit,
	}, nil
}
