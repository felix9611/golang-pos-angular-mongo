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

func BatchCreateVendors(Vendors []models.Vendors) (interface{}, error) {
	for _, vendor := range Vendors {
		result, _ := CreateVendor(&vendor)

		if result == nil {
			return "failed to create vendor", nil
		}
	}

	return "batch insert completed", nil
}

func CreateVendor(vendor *models.Vendors) (interface{}, error) {
	collection := config.GetCollection("vendors")

	filter := bson.M{"status": 1, "vendorName": vendor.VendorName, "vendorCode": vendor.VendorCode}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter)

	if err != nil {
		return nil, err
	}

	if count == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if vendor.Status == 0 {
			vendor.Status = 1
		}

		if vendor.CreatedAt.IsZero() {
			vendor.CreatedAt = time.Now()
		}

		if vendor.UpdatedAt.IsZero() {
			vendor.UpdatedAt = time.Now()
		}

		_, err := collection.InsertOne(ctx, vendor)
		if err != nil {
			return nil, err
		}

		CreateActionRecord("Vendor Create", "POST", "Vendor", vendor, "Success")

		return vendor, nil
	} else {
		CreateActionRecord("Vendor Create", "POST", "Vendor", vendor, "Failed")
		return "Vendor with the same name already exists", nil
	}

}

func GetOneVendorById(id string) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.GetCollection("vendors")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID, "status": 1}

	var vendor models.Vendors
	err = collection.FindOne(ctx, filter).Decode(&vendor)

	if err != nil {
		return nil, err
	}

	if vendor.Status == 0 {
		return "Vendor is inactive", nil
	} else {
		return &vendor, nil
	}

}

func UpdateVendorById(updateData *models.Vendors) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.GetCollection("vendors")
	filter := bson.M{"_id": updateData.ID, "status": 1}

	var existingVendor models.Vendors
	err := collection.FindOne(ctx, filter).Decode(&existingVendor)
	if err != nil {
		return nil, err
	}

	if existingVendor.Status == 1 {
		updateData.UpdatedAt = time.Now()

		result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": updateData})
		if err != nil {
			return nil, err
		}

		CreateActionRecord("Vendor Update", "POST", "Vendor", updateData, "Success")

		return result, nil
	} else {
		CreateActionRecord("Vendor Update", "POST", "Vendor", updateData, "Failed")
		return "Vendor is inactive", nil
	}
}

func InactiveVendorByID(id string) (interface{}, error) {
	collection := config.GetCollection("vendors")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID, "status": 1}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var vendor models.Vendors

	err = collection.FindOne(ctx, filter).Decode(&vendor)
	if err != nil {
		return nil, err
	}

	if vendor.Status == 1 {
		vendor.Status = 0
		vendor.UpdatedAt = time.Now()

		result, err := collection.UpdateOne(ctx, filter, bson.M{
			"$set": vendor,
		})

		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		return "Vendor is already inactive", nil
	}
}

func GetAllVendors() ([]models.Vendors, error) {
	collection := config.GetCollection("vendors")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var vendors []models.Vendors
	cursor, err := collection.Find(ctx, bson.M{"status": 1})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var vendor models.Vendors
		if err := cursor.Decode(&vendor); err != nil {
			return nil, err
		}
		vendors = append(vendors, vendor)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return vendors, nil
}

func VendorListwithFilter(pageDto *dto.VendorPageDto) (interface{}, error) {

	collection := config.GetCollection("vendors")
	filter := bson.M{
		"status": 1,
		"$or": []bson.M{
			{"name": bson.M{"$regex": pageDto.Name, "$options": "i"}},
			{"email": bson.M{"$regex": pageDto.Name, "$options": "i"}},
			{"phone": bson.M{"$regex": pageDto.Name, "$options": "i"}},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", -1}})
	cursor, err := collection.Find(ctx, filter, findOptions)

	if err != nil {
		return nil, err
	}

	var vendors []models.Vendors
	for cursor.Next(ctx) {
		var vendor models.Vendors
		if err := cursor.Decode(&vendor); err != nil {
			return nil, err
		}
		vendors = append(vendors, vendor)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return vendors, nil

}

func VendorList(pageDto *dto.VendorPageDto) (interface{}, error) {
	if pageDto.Page < 1 {
		pageDto.Page = 1
	}
	if pageDto.Limit < 1 {
		pageDto.Limit = 10
	}

	skip := (pageDto.Page - 1) * pageDto.Limit
	limit := pageDto.Limit

	collection := config.GetCollection("vendors")
	filter := bson.M{
		"status": 1,
		"$or": []bson.M{
			{"name": bson.M{"$regex": pageDto.Name, "$options": "i"}},
			{"email": bson.M{"$regex": pageDto.Name, "$options": "i"}},
			{"phone": bson.M{"$regex": pageDto.Name, "$options": "i"}},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, errCount := collection.CountDocuments(ctx, filter)
	if errCount != nil {
		return nil, errCount
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{"created_at", -1}})
	cursor, err := collection.Find(ctx, filter, findOptions)

	if err != nil {
		return nil, err
	}

	var vendors []models.Vendors
	for cursor.Next(ctx) {
		var vendor models.Vendors
		if err := cursor.Decode(&vendor); err != nil {
			return nil, err
		}
		vendors = append(vendors, vendor)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return gin.H{"lists": vendors, "total": count, "page": pageDto.Page, "limit": pageDto.Limit}, nil
}
