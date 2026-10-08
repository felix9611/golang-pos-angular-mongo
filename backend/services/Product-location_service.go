package services

import (
	"context"
	"golang-fixedasset-mongo-backend/backend/config"
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProductLocationList(req dto.ListProductLocationtRequestDto) (interface{}, error) {
	skip := int64((req.Page - 1) * req.Limit)
	limit := int64(req.Limit)

	collection := config.GetCollection("product_location")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filters := bson.M{}

	if len(req.LocationIds) > 0 {
		filters["locationId"] = bson.M{"$in": req.LocationIds}
	}

	if req.ProductCode != "" {
		filters["productCode"] = req.ProductCode
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filters}},

		// lookup product
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "product_lists"},
				{Key: "let", Value: bson.D{
					{Key: "productIdStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$productId"},
							{Key: "to", Value: "objectId"},
							{Key: "onError", Value: nil},
							{Key: "onNull", Value: nil},
						}},
					}},
				}},
				{Key: "pipeline", Value: mongo.Pipeline{
					{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$_id", "$$productIdStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "product"},
			},
		}},
		// lookup location
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "locations"},
				{Key: "let", Value: bson.D{
					{Key: "locationIdStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$locationId"},
							{Key: "to", Value: "objectId"},
							{Key: "onError", Value: nil},
							{Key: "onNull", Value: nil},
						}},
					}},
				}},
				{Key: "pipeline", Value: mongo.Pipeline{
					{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$_id", "$$locationIdStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "location"},
			}}},
		// unwind
		{{Key: "$unwind", Value: bson.D{bson.E{Key: "path", Value: "$product"}, bson.E{Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{bson.E{Key: "path", Value: "$location"}, bson.E{Key: "preserveNullAndEmptyArrays", Value: true}}}},

		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	count, errCount := collection.CountDocuments(ctx, filters)
	if errCount != nil {
		return nil, errCount
	}

	var results []dto.ProductLocationLists
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return dto.ListProductLocationts{Lists: results, Total: count, Page: req.Page, Limit: req.Limit}, nil
}

func StockIn(entryData dto.StockInOutProductLocationDto) (interface{}, error) {

	productInfo, err := GetOneProductItemByID(entryData.ProductId)
	if err != nil {
		return nil, err
	}

	locationInfo, err := GetOneLocationById(entryData.LocationId)
	if err != nil {
		return nil, err
	}

	if productInfo == nil || locationInfo == nil {
		return gin.H{"error": "Product or Location not found!"}, err
	}

	collection := config.GetCollection("product_location")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"productId": entryData.ProductId, "locationId": entryData.LocationId, "status": 1}

	var existingProductLocation models.ProductLocations
	err = collection.FindOne(ctx, filter).Decode(&existingProductLocation)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			var productLocation models.ProductLocations

			productLocation.ProductId = entryData.ProductId
			productLocation.LocationId = entryData.LocationId
			productLocation.Qty = entryData.Qty
			productLocation.TotalPrice = entryData.TotalPrice
			productLocation.TotalCost = entryData.TotalCost

			_, err := collection.InsertOne(ctx, productLocation)
			if err != nil {
				return gin.H{"error": "Failed to create new product location record!"}, err
			}

			insertRecord := dto.InsertInvRecordDto{
				ProductId: entryData.ProductId,
				LocFrom:   entryData.LocationId,
				LocTo:     entryData.LocationId,
				Qty:       entryData.Qty,
				Cost:      entryData.TotalCost,
			}

			_, err = InsertInventoryRecord(insertRecord)
			if err != nil {
				return gin.H{"error": "Failed to create new inventory record!"}, err
			}

			CreateActionRecord("Product Location Record Create", "POST", "Product Location", productLocation, "Success")

			return gin.H{"message": "Product location record created successfully!"}, nil
		} else {

			updateData := bson.M{
				"$set": bson.M{
					"qty":        entryData.Qty + existingProductLocation.Qty,
					"totalPrice": entryData.TotalPrice + existingProductLocation.TotalPrice,
					"totalCost":  entryData.TotalCost + existingProductLocation.TotalCost,
				},
			}

			_, err := collection.UpdateOne(ctx, filter, updateData)
			if err != nil {
				return gin.H{"error": "Failed to update product location record!"}, err
			}

			insertRecord := dto.InsertInvRecordDto{
				ProductId: entryData.ProductId,
				LocFrom:   entryData.LocationId,
				LocTo:     entryData.LocationId,
				Qty:       entryData.Qty,
				Cost:      entryData.TotalCost,
			}

			_, err = InsertInventoryRecord(insertRecord)
			if err != nil {
				return gin.H{"error": "Failed to create new inventory record!"}, err
			}

			CreateActionRecord("Product Location Record Create", "POST", "Product Location", existingProductLocation, "Success")

			return gin.H{"message": "Product location record updated successfully!", "status": true}, nil

		}
	}

	return gin.H{"error": "Oooops! Sometimes ie error! Please try again", "status": false}, err
}
