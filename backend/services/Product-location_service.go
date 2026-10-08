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
	if err != nil || productInfo == nil {
		return gin.H{"error": "Product not found!"}, err
	}

	locationInfo, err := GetOneLocationById(entryData.LocationId)
	if err != nil || locationInfo == nil {
		return gin.H{"error": "Location not found!"}, err
	}

	collection := config.GetCollection("product_location")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"productId":  entryData.ProductId,
		"locationId": entryData.LocationId,
	}

	var existingProductLocation models.ProductLocations
	err = collection.FindOne(ctx, filter).Decode(&existingProductLocation)

	if err != nil {
		if err == mongo.ErrNoDocuments {

			productLocation := models.ProductLocations{
				ProductId:  entryData.ProductId,
				LocationId: entryData.LocationId,
				Qty:        entryData.Qty,
				TotalPrice: entryData.TotalPrice,
				TotalCost:  entryData.TotalCost,
			}

			if _, err := collection.InsertOne(ctx, productLocation); err != nil {
				return gin.H{"error": "Failed to create new product location record!"}, err
			}

			insertRecord := dto.InsertInvRecordDto{
				ProductId: entryData.ProductId,
				LocFrom:   entryData.LocationId,
				LocTo:     entryData.LocationId,
				Qty:       entryData.Qty,
				Cost:      entryData.TotalCost,
			}

			if _, err = InsertInventoryRecord(insertRecord); err != nil {
				return gin.H{"error": "Failed to create new inventory record!"}, err
			}

			CreateActionRecord("Product Location Record Create", "POST", "Product Location", productLocation, "Success")
			return gin.H{"message": "Product location record created successfully!"}, nil
		}

		return gin.H{"error": "Database query failed!"}, err
	}

	newQty := existingProductLocation.Qty + entryData.Qty
	newTotalPrice := existingProductLocation.TotalPrice + entryData.TotalPrice
	newTotalCost := existingProductLocation.TotalCost + entryData.TotalCost

	updateData := bson.M{
		"$set": bson.M{
			"qty":        newQty,
			"totalPrice": newTotalPrice,
			"totalCost":  newTotalCost,
		},
	}

	if _, err = collection.UpdateOne(ctx, filter, updateData); err != nil {
		return gin.H{"error": "Failed to update product location record!"}, err
	}

	insertRecord := dto.InsertInvRecordDto{
		ProductId: entryData.ProductId,
		LocFrom:   entryData.LocationId,
		LocTo:     entryData.LocationId,
		Qty:       entryData.Qty,
		Cost:      entryData.TotalCost,
	}

	if _, err = InsertInventoryRecord(insertRecord); err != nil {
		return gin.H{"error": "Failed to create new inventory record!"}, err
	}

	CreateActionRecord("Product Location Record Update", "PUT", "Product Location", existingProductLocation, "Success")

	return gin.H{"message": "Product location record updated successfully!", "status": true}, nil
}

func StockMove(entryData dto.StockMoveProductLocationDto) (interface{}, error) {
	productInfo, err := GetOneProductItemByID(entryData.ProductId)
	if err != nil {
		return nil, err
	}

	locationFromInfo, err := GetOneLocationById(entryData.FromLocationId)
	if err != nil {
		return nil, err
	}

	locationFToInfo, err := GetOneLocationById(entryData.ToLocationId)
	if err != nil {
		return nil, err
	}

	if productInfo == nil || locationFromInfo == nil || locationFToInfo == nil {
		return gin.H{"error": "Product or Location not found!"}, err
	}

	collection := config.GetCollection("product_location")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fromFilter := bson.M{"productId": entryData.ProductId, "locationId": entryData.FromLocationId}
	toFilter := bson.M{"productId": entryData.ProductId, "locationId": entryData.ToLocationId}

	var fromProductLocation models.ProductLocations
	var toProductLocation models.ProductLocations

	err = collection.FindOne(ctx, fromFilter).Decode(&fromProductLocation)
	if err != nil {
		return gin.H{"error": "Product location not found!"}, err
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return gin.H{"error": "Oooops! No any data matched", "status": false}, err
		}
	}

	err = collection.FindOne(ctx, toFilter).Decode(&toProductLocation)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			if err == mongo.ErrNoDocuments {

				newQty := fromProductLocation.Qty - entryData.Qty
				newTotalPrice := fromProductLocation.TotalPrice - entryData.TotalPrice
				newTotalCost := fromProductLocation.TotalCost - entryData.TotalCost

				updateData := bson.M{
					"$set": bson.M{
						"qty":        newQty,
						"totalPrice": newTotalPrice,
						"totalCost":  newTotalCost,
					},
				}

				if _, err = collection.UpdateOne(ctx, fromFilter, updateData); err != nil {
					return gin.H{"error": "Failed to update product location record!"}, err
				}

				productLocation := models.ProductLocations{
					ProductId:  entryData.ProductId,
					LocationId: entryData.ToLocationId,
					Qty:        entryData.Qty,
					TotalPrice: entryData.TotalPrice,
					TotalCost:  entryData.TotalCost,
				}

				if _, err := collection.InsertOne(ctx, productLocation); err != nil {
					return gin.H{"error": "Failed to create new product location record!"}, err
				}

				recordQtys := 0 - entryData.Qty
				insertRecord := dto.InsertInvRecordDto{
					ProductId: entryData.ProductId,
					LocFrom:   entryData.FromLocationId,
					LocTo:     entryData.ToLocationId,
					Qty:       recordQtys,
					Cost:      entryData.TotalCost,
				}

				if _, err = InsertInventoryRecord(insertRecord); err != nil {
					return gin.H{"error": "Failed to create new inventory record!"}, err
				}

				CreateActionRecord("Product Location Record Create", "POST", "Product Location", productLocation, "Success")
				return gin.H{"message": "Product location record created successfully!"}, nil
			}

			return gin.H{"error": "Database query failed!"}, err
		}
	}

	newQty := fromProductLocation.Qty - entryData.Qty
	newTotalPrice := fromProductLocation.TotalPrice - entryData.TotalPrice
	newTotalCost := fromProductLocation.TotalCost - entryData.TotalCost

	updateData := bson.M{
		"$set": bson.M{
			"qty":        newQty,
			"totalPrice": newTotalPrice,
			"totalCost":  newTotalCost,
		},
	}

	if _, err = collection.UpdateOne(ctx, fromFilter, updateData); err != nil {
		return gin.H{"error": "Failed to update product location record!"}, err
	}

	newPlaceQty := toProductLocation.Qty + entryData.Qty
	newPlaceTotalPrice := toProductLocation.TotalPrice + entryData.TotalPrice
	newPlaceTotalCost := toProductLocation.TotalCost + entryData.TotalCost

	updateNewData := bson.M{
		"$set": bson.M{
			"qty":        newPlaceQty,
			"totalPrice": newPlaceTotalPrice,
			"totalCost":  newPlaceTotalCost,
		},
	}

	if _, err = collection.UpdateOne(ctx, toFilter, updateNewData); err != nil {
		return gin.H{"error": "Failed to update product location record!"}, err
	}

	recordQtys := 0 - entryData.Qty
	insertRecord := dto.InsertInvRecordDto{
		ProductId: entryData.ProductId,
		LocFrom:   entryData.FromLocationId,
		LocTo:     entryData.ToLocationId,
		Qty:       recordQtys,
		Cost:      entryData.TotalCost,
	}

	if _, err = InsertInventoryRecord(insertRecord); err != nil {
		return gin.H{"error": "Failed to create new inventory record!"}, err
	}

	CreateActionRecord("Product Location Record Update", "PUT", "Product Location", toProductLocation, "Success")

	return gin.H{"message": "Product location record updated successfully!", "status": true}, nil
}

func StockOut(entryData dto.StockInOutProductLocationDto) (interface{}, error) {
	productInfo, err := GetOneProductItemByID(entryData.ProductId)
	if err != nil || productInfo == nil {
		return gin.H{"error": "Product not found!"}, err
	}

	locationInfo, err := GetOneLocationById(entryData.LocationId)
	if err != nil || locationInfo == nil {
		return gin.H{"error": "Location not found!"}, err
	}

	if productInfo == nil || locationInfo == nil {
		return gin.H{"error": "Invalid data!"}, err
	}

	collection := config.GetCollection("product_location")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fromFilter := bson.M{"productId": entryData.ProductId, "locationId": entryData.LocationId}

	var fromProductLocation models.ProductLocations
	err = collection.FindOne(ctx, fromFilter).Decode(&fromProductLocation)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return gin.H{"error": "Product location record not found!"}, err
		}

		return gin.H{"errore": "Database query failed!"}, err
	}

	if fromProductLocation.Qty < entryData.Qty {
		return gin.H{"error": "Insufficient stock!"}, err
	}

	newQty := fromProductLocation.Qty - entryData.Qty
	newTotalPrice := fromProductLocation.TotalPrice - entryData.TotalPrice
	newTotalCost := fromProductLocation.TotalCost - entryData.TotalCost

	updateData := bson.M{
		"$set": bson.M{
			"qty":        newQty,
			"totalPrice": newTotalPrice,
			"totalCost":  newTotalCost,
		},
	}

	if _, err = collection.UpdateOne(ctx, fromFilter, updateData); err != nil {
		return gin.H{"error": "Failed to update product location record!"}, err
	}

	recordQtys := 0 - entryData.Qty
	insertRecord := dto.InsertInvRecordDto{
		ProductId: entryData.ProductId,
		LocFrom:   entryData.LocationId,
		LocTo:     "",
		Qty:       recordQtys,
		Cost:      entryData.TotalCost,
	}

	if _, err = InsertInventoryRecord(insertRecord); err != nil {
		return gin.H{"error": "Failed to create new inventory record!"}, err
	}

	CreateActionRecord("Product Location Record Update", "PUT", "Product Location", fromProductLocation, "Success")

	return gin.H{"message": "Product location record updated successfully!", "status": true}, nil

}
