package services

import (
	"context"
	"golang-fixedasset-mongo-backend/backend/config"
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// "go.mongodb.org/mongo-driver/mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func BatchCreateLocation(locations []models.Locations) (interface{}, error) {
	for _, location := range locations {
		result, _ := CreateLocation(&location)

		if result == nil {
			return "failed to create vendor", nil
		}
	}

	return "batch insert completed", nil
}

func CreateLocation(location *models.Locations) (interface{}, error) {
	collection := config.GetCollection("locations")

	filter := bson.M{"status": 1, "vendorName": location.PlaceCode, "placeCode": location.PlaceName}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter)

	if err != nil {
		return nil, err
	}

	if count == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if location.Status == 0 {
			location.Status = 1
		}

		if location.CreatedAt.IsZero() {
			location.CreatedAt = time.Now()
		}

		if location.UpdatedAt.IsZero() {
			location.UpdatedAt = time.Now()
		}

		_, err := collection.InsertOne(ctx, location)

		if err != nil {
			return nil, err
		}

		CreateActionRecord("Location Create", "POST", "Location", location, "Success")

		return location, nil

	} else {
		CreateActionRecord("Location Create", "POST", "Location", location, "Failed")
		return "Location with the same name already exists", nil
	}
}

func GetOneLocationById(id string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.GetCollection("locations")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID, "status": 1}

	var location models.Locations
	err = collection.FindOne(ctx, filter).Decode(&location)

	if err != nil {
		return nil, err
	}

	if location.Status == 0 {
		return "Location is inactive", nil
	} else {
		return &location, nil
	}
}

func UpdateLocationById(updateData *models.Locations) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.GetCollection("locations")
	filter := bson.M{"_id": updateData.ID, "status": 1}

	var existingLocation models.Locations
	err := collection.FindOne(ctx, filter).Decode(&existingLocation)
	if err != nil {
		return nil, err
	}

	if existingLocation.Status == 1 {
		updateData.UpdatedAt = time.Now()

		result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": updateData})
		if err != nil {
			return nil, err
		}

		CreateActionRecord("Location Update", "POST", "Location", updateData, "Success")

		return result, nil
	} else {
		CreateActionRecord("Location Update", "POST", "Location", updateData, "Failed")
		return "Location is inactive", nil
	}
}

func InactiveLocationByID(id string) (interface{}, error) {
	collection := config.GetCollection("locations")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID, "status": 1}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingLocation models.Locations
	err = collection.FindOne(ctx, filter).Decode(&existingLocation)
	if err != nil {
		return nil, err
	}

	if existingLocation.Status == 1 {
		existingLocation.Status = 0
		existingLocation.UpdatedAt = time.Now()
		_, err := collection.UpdateOne(ctx, filter, bson.M{
			"$set": existingLocation,
		})

		if err != nil {
			return nil, err
		}

		CreateActionRecord("Location Inactive", "DELETE", "Location", existingLocation, "Success")

		return "Location deactivated successfully", nil
	} else {
		CreateActionRecord("Location Inactive", "DELETE", "Location", existingLocation, "Failed")
		return "Location is already inactive", nil
	}
}

func LocationListWithFilter(pageDto *dto.LocationPageDto) (interface{}, error) {
	collection := config.GetCollection("locations")
	filter := bson.M{
		"status": 1,
		"$or": []bson.M{
			{"placeCode": bson.M{"$regex": pageDto.Name, "$options": "i"}},
			{"placeName": bson.M{"$regex": pageDto.Name, "$options": "i"}},
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

	var locations []models.Locations
	for cursor.Next(ctx) {
		var location models.Locations
		if err := cursor.Decode(&location); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(ctx)
	return locations, nil

}

func LocationList(pageDto *dto.LocationPageDto) (interface{}, error) {
	if pageDto.Page < 1 {
		pageDto.Page = 1
	}
	if pageDto.Limit < 1 {
		pageDto.Limit = 10
	}

	skip := (pageDto.Page - 1) * pageDto.Limit
	limit := pageDto.Limit

	collection := config.GetCollection("locations")
	filter := bson.M{
		"status": 1,
		"$or": []bson.M{
			{"placeCode": bson.M{"$regex": pageDto.Name, "$options": "i"}},
			{"placeName": bson.M{"$regex": pageDto.Name, "$options": "i"}},
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

	var locations []models.Locations
	for cursor.Next(ctx) {
		var location models.Locations
		if err := cursor.Decode(&location); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return gin.H{"lists": locations, "total": count, "page": pageDto.Page, "limit": pageDto.Limit}, nil
}

func ListAllLocation() (interface{}, error) {
	collection := config.GetCollection("locations")

	filter := bson.M{"status": 1}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var locations []models.Locations
	for cursor.Next(ctx) {
		var location models.Locations
		if err := cursor.Decode(&location); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return gin.H{"data": locations}, nil
}

func LocationDataFinder(placeCode string, placeName string) (models.Locations, error) {
	collection := config.GetCollection("locations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1, "placeCode": placeCode, "placeName": placeName}

	var location models.Locations

	err := collection.FindOne(ctx, filter).Decode(&location)
	if err != nil {
		return models.Locations{}, err
	}

	if location.Status == 0 {
		return models.Locations{}, nil
	} else {
		return location, nil
	}
}
