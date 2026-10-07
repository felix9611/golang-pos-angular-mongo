package services

import (
	"context"
	"golang-fixedasset-mongo-backend/backend/config"
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func BatchCreateAssetTypes(assetTypes []models.AssetTypes) (interface{}, error) {
	for _, assetType := range assetTypes {
		result, _ := CreateAssetType(&assetType)

		if result == nil {
			return "failed to create code type", nil
		}
	}

	return "batch insert completed", nil
}

func CreateAssetType(assetType *models.AssetTypes) (interface{}, error) {
	filter := bson.M{"status": 1, "typeCode": assetType.TypeCode, "typeName": assetType.TypeName}

	collection := config.GetCollection("asset_types")

	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		if assetType.Status == 0 {
			assetType.Status = 1
		}

		if assetType.CreatedAt.IsZero() {
			assetType.CreatedAt = time.Now()
		}

		if assetType.UpdatedAt.IsZero() {
			assetType.UpdatedAt = time.Now()
		}

		CreateActionRecord("Asset Type Create", "POST", "Asset Type", assetType, "Success")

		result, err := collection.InsertOne(context.Background(), assetType)
		if err != nil {
			return nil, err
		}

		return result, nil
	} else {
		CreateActionRecord("Asset Type Create", "POST", "Asset Type", assetType, "Failed")
		return "This Asset Type with the same name already exists", nil
	}
}

func GetOneAssetType(id string) (interface{}, error) {
	collection := config.GetCollection("asset_types")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var assetType models.AssetTypes
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, filter).Decode(&assetType)
	if err != nil {
		return nil, err
	}

	if assetType.Status == 0 {
		return "This Asset Type is inactive", nil
	} else {
		return &assetType, nil
	}
}

func VoidOneAssetType(id string) (interface{}, error) {
	collection := config.GetCollection("asset_types")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var assetType models.AssetTypes
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, filter).Decode(&assetType)
	if err != nil {
		return nil, err
	}

	if assetType.Status == 0 {
		CreateActionRecord("Asset Type Void", "DELETE", "Asset Type", assetType, "Failed")
		return "This Asset Type is inactive", nil
	} else {
		assetType.Status = 0
		assetType.UpdatedAt = time.Now()

		_, err := collection.UpdateOne(ctx, filter, assetType)
		if err != nil {
			return nil, err
		}
		CreateActionRecord("Asset Type Void", "DELETE", "Asset Type", assetType, "Success")
		return "This Asset Type has been voided just now", nil
	}
}

func UpdateAssetType(updateData *models.AssetTypes) (interface{}, error) {
	collection := config.GetCollection("asset_types")

	filter := bson.M{"_id": updateData.ID}

	var assetType models.AssetTypes
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&assetType)
	if err != nil {
		return nil, err
	}

	if assetType.Status == 0 {
		CreateActionRecord("Asset Type Update", "POST", "Asset Type", updateData, "Failed")
		return "This Asset Type is inactive", nil
	} else {
		updateData.UpdatedAt = time.Now()

		result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": updateData})

		CreateActionRecord("Asset Type Update", "POST", "Asset Type", updateData, "Success")

		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func ListAllAssetType() (interface{}, error) {
	collection := config.GetCollection("asset_types")

	filter := bson.M{"status": 1}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assetTypes []models.AssetTypes
	for cursor.Next(ctx) {
		var assetType models.AssetTypes
		if err := cursor.Decode(&assetType); err != nil {
			return nil, err
		}
		assetTypes = append(assetTypes, assetType)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return gin.H{"datas": assetTypes}, nil
}

func ListAssetTypeNoPaging(pageDto *dto.AssetTypeListDto) (interface{}, error) {

	collection := config.GetCollection("asset_types")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1}

	if pageDto.Name != "" {
		filter = bson.M{
			"status": 1,
			"$or": []bson.M{
				{"typeCode": bson.M{"$regex": pageDto.Name, "$options": "i"}},
				{"typeName": bson.M{"$regex": pageDto.Name, "$options": "i"}},
			},
		}
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var assetTypes []models.AssetTypes
	for cursor.Next(ctx) {
		var assetType models.AssetTypes
		if err := cursor.Decode(&assetType); err != nil {
			return nil, err
		}
		assetTypes = append(assetTypes, assetType)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return assetTypes, nil
}

func ListAssetType(pageDto *dto.AssetTypeListDto) (interface{}, error) {

	if pageDto.Page < 1 {
		pageDto.Page = 1
	}

	if pageDto.Limit < 1 {
		pageDto.Limit = 10
	}

	collection := config.GetCollection("asset_types")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1}

	if pageDto.Name != "" {
		filter = bson.M{
			"status": 1,
			"$or": []bson.M{
				{"typeCode": bson.M{"$regex": pageDto.Name, "$options": "i"}},
				{"typeName": bson.M{"$regex": pageDto.Name, "$options": "i"}},
			},
		}
	}

	count, errCount := collection.CountDocuments(ctx, filter)
	if errCount != nil {
		return nil, errCount
	}

	skip := int64((pageDto.Page - 1) * pageDto.Limit)
	limit := int64(pageDto.Limit)

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var assetTypes []models.AssetTypes
	for cursor.Next(ctx) {
		var assetType models.AssetTypes
		if err := cursor.Decode(&assetType); err != nil {
			return nil, err
		}
		assetTypes = append(assetTypes, assetType)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return gin.H{"lists": assetTypes, "total": count, "page": pageDto.Page, "limit": pageDto.Limit}, nil
}

func AssetTypeDataFinderReturn(typeCode string, typeName string) (models.AssetTypes, error) {
	collection := config.GetCollection("asset_types")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1, "typeCode": typeCode, "typeName": typeName}

	var assetType models.AssetTypes

	err := collection.FindOne(ctx, filter).Decode(&assetType)
	if err != nil {
		return models.AssetTypes{}, err
	}

	if assetType.Status == 0 {
		return models.AssetTypes{}, nil
	}

	return assetType, nil
}
