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

func BatchCreateProductTypes(productTypes []models.ProductTypes) (interface{}, error) {
	for _, productType := range productTypes {
		result, _ := CreateProductType(&productType)

		if result == nil {
			return gin.H{"error": "failed to create product type"}, nil
		}
	}

	return gin.H{"message": "batch insert completed"}, nil
}

func CreateProductType(productType *models.ProductTypes) (interface{}, error) {
	filter := bson.M{
		"typeCode": productType.TypeCode,
		"typeName": productType.TypeName,
	}

	collection := config.GetCollection("product_types")
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		productType.Status = 1

		if productType.CreatedAt.IsZero() {
			productType.CreatedAt = time.Now()
		}

		if productType.UpdatedAt.IsZero() {
			productType.UpdatedAt = time.Now()
		}

		CreateActionRecord("Product Type Create", "POST", "Product Type", productType, "Success")

		result, err := collection.InsertOne(context.Background(), productType)
		if err != nil {
			return nil, err
		}

		return result.InsertedID, nil
	} else {
		return gin.H{"error": "This product type Record with the same details already exists"}, nil
	}
}

func GetOneProductType(id string) (interface{}, error) {

	collection := config.GetCollection("product_types")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var productType models.ProductTypes

	err = collection.FindOne(ctx, filter).Decode(&productType)
	if err != nil {
		return nil, err
	}

	if productType.Status == 0 {
		return gin.H{"error": "This Product Type Record is inactive"}, nil
	} else {
		return &productType, nil
	}
}

func VoidOneProductType(id string) (interface{}, error) {
	collection := config.GetCollection("product_types")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var productType models.ProductTypes
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, filter).Decode(&productType)
	if err != nil {
		return nil, err
	}

	if productType.Status == 0 {
		CreateActionRecord("Product Type Void", "DELETE", "Product Type", productType, "Failed")
		return gin.H{"error": "This Product Type Record is inactive"}, nil
	} else {
		productType.Status = 0
		productType.UpdatedAt = time.Now()

		_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": productType})
		if err != nil {
			return nil, err
		}

		CreateActionRecord("Product Type Void", "DELETE", "Product Type", productType, "Success")
		return gin.H{"message": "Product Type Record Voided Successfully"}, nil
	}
}

func UpdateProductType(updateData *models.ProductTypes) (interface{}, error) {
	collection := config.GetCollection("product_types")

	filter := bson.M{"_id": updateData.ID}

	var productType models.ProductTypes
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&productType)
	if err != nil {
		return nil, err
	}

	if productType.Status == 0 {
		CreateActionRecord("Product Type Update", "POST", "Product Type", updateData, "Failed")
		return gin.H{"error": "This Product Type Record is inactive"}, nil
	} else {
		updateData.UpdatedAt = time.Now()

		result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": updateData})
		if err != nil {
			return nil, err
		}
		CreateActionRecord("Product Type Update", "POST", "Product Type", updateData, "Success")
		return result, nil
	}
}

func ListAlProductType() (interface{}, error) {
	collection := config.GetCollection("product_types")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var productTypes []models.ProductTypes
	for cursor.Next(ctx) {
		var productType models.ProductTypes
		if err := cursor.Decode(&productType); err != nil {
			return nil, err
		}
		productTypes = append(productTypes, productType)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return productTypes, nil
}

func ListProductTypeWithFilter(pageDto *dto.ProductTypeListDto) (interface{}, error) {
	collection := config.GetCollection("product_types")
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

	var productTypes []models.ProductTypes
	for cursor.Next(ctx) {
		var productType models.ProductTypes
		if err := cursor.Decode(&productType); err != nil {
			return nil, err
		}
		productTypes = append(productTypes, productType)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return productTypes, nil
}

func ListProductType(pageDto *dto.ProductTypeListDto) (interface{}, error) {
	collection := config.GetCollection("product_types")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if pageDto.Page < 1 {
		pageDto.Page = 1
	}

	if pageDto.Limit < 1 {
		pageDto.Limit = 10
	}

	filter := bson.M{"status": 1}

	skip := (pageDto.Page - 1) * pageDto.Limit
	limit := pageDto.Limit

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

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var productTypes []models.ProductTypes
	for cursor.Next(ctx) {
		var productType models.ProductTypes
		if err := cursor.Decode(&productType); err != nil {
			return nil, err
		}
		productTypes = append(productTypes, productType)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return gin.H{"lists": productTypes, "total": count, "page": pageDto.Page, "limit": pageDto.Limit}, nil
}

func ProductTypeDataFinderReturn(typeCode string, typeName string) (models.ProductTypes, error) {
	collection := config.GetCollection("product_types")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1, "typeCode": typeCode, "typeName": typeName}

	var productType models.ProductTypes

	err := collection.FindOne(ctx, filter).Decode(&productType)
	if err != nil {
		return models.ProductTypes{}, err
	}

	if productType.Status == 0 {
		return models.ProductTypes{}, nil
	}

	return productType, nil
}
