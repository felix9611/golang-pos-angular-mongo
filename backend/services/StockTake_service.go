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
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateStockTakeForm(stockTakeForm *models.StockTakes) (interface{}, error) {

	collection := config.GetCollection("stock_takes")

	filter := bson.M{"actionName": stockTakeForm.ActionName, "actionPlaceId": stockTakeForm.ActionPlaceId}

	var existingStockTakeForm models.StockTakes

	err := collection.FindOne(context.Background(), filter).Decode(&existingStockTakeForm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			stockTakeForm.CreatedTime = time.Now()
			stockTakeForm.Status = 1 // 1 means "in progress"
			result, err := collection.InsertOne(context.Background(), stockTakeForm)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, err
	}
	return "A Stock Take form already exists", nil
}

func GetOneStockTakeFormByID(id string) (interface{}, error) {
	collection := config.GetCollection("stock_takes")
	var stockTakeForm models.StockTakes

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err2 := collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&stockTakeForm)
	if err2 != nil {
		return nil, err2
	}

	stockTakeItems, err := GetStockTakeItem(id)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"id":             stockTakeForm.ID,
		"actionName":     stockTakeForm.ActionName,
		"actionPlaceId":  stockTakeForm.ActionPlaceId,
		"remark":         stockTakeForm.Remark,
		"createdTime":    stockTakeForm.CreatedTime,
		"status":         stockTakeForm.Status,
		"createdBy":      stockTakeForm.CreatedBy,
		"finishTime":     stockTakeForm.FinishTime,
		"finishBy":       stockTakeForm.FinishBy,
		"stockTakeItems": stockTakeItems,
	}, nil

}

func ListStockTakeForms(pageDto *dto.ListStockTakeDto) (interface{}, error) {
	if pageDto.Page < 1 {
		pageDto.Page = 1
	}
	if pageDto.Limit < 1 {
		pageDto.Limit = 10
	}

	skip := (pageDto.Page - 1) * pageDto.Limit
	limit := pageDto.Limit

	filters := bson.M{}

	if pageDto.Name == "" {
		filters["actionName"] = bson.M{"$regex": pageDto.Name, "$options": "i"}
	}

	if len(pageDto.PlaceIds) > 0 {
		filters["actionPlaceId"] = bson.M{"$in": pageDto.PlaceIds}
	}

	collection := config.GetCollection("stock_takes")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filters}},
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "locations"},
				{Key: "let", Value: bson.D{
					{Key: "placeIdStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$actionPlaceId"},
							{Key: "to", Value: "objectId"},
							{Key: "onError", Value: nil},
							{Key: "onNull", Value: nil},
						}},
					}},
				}},
				{Key: "pipeline", Value: mongo.Pipeline{
					{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$_id", "$$placeIdStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "location"},
			},
		}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$location"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$sort", Value: bson.D{{"createdTime", -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	count, errCount := collection.CountDocuments(ctx, filters)
	if errCount != nil {
		return nil, errCount
	}

	return gin.H{"lists": results, "total": count, "page": pageDto.Page, "limit": pageDto.Limit}, nil
}

func UpdateStockTakeForm(updateData *models.StockTakes) (interface{}, error) {
	collection := config.GetCollection("stock_takes")

	filter := bson.M{"_id": updateData.ID}

	var existingStockTake models.StockTakes

	err := collection.FindOne(context.Background(), filter).Decode(&existingStockTake)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "Stock Take form not found", nil
		}
		return nil, err
	}

	if existingStockTake.Status == 0 {
		return "Ooooops! This stock take form no longer active! Please create a new form!", nil
	} else if existingStockTake.Status == 2 {
		return "Ooooops! This stock take form has been completed! Please create a new form!", nil
	} else {
		res, err2 := collection.UpdateOne(context.Background(), bson.M{"_id": updateData.ID}, bson.M{"$set": updateData})
		if err2 != nil {
			return nil, err2
		}

		return res, nil
	}
}

func GetStockTakeItem(stockTakeId string) (interface{}, error) {
	collection := config.GetCollection("stock_take_items")

	objectID, err := primitive.ObjectIDFromHex(stockTakeId)

	if err != nil {
		return nil, err
	}

	filter := bson.M{"stockTakeId": objectID}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "asset_lists"},
				{Key: "let", Value: bson.D{
					{Key: "assetIdStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$assetId"},
							{Key: "to", Value: "objectId"},
							{Key: "onError", Value: nil},
							{Key: "onNull", Value: nil},
						}},
					}},
				}},
				{Key: "pipeline", Value: mongo.Pipeline{
					{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$_id", "$$assetIdStr"}},
						}},
					}}},
					{{
						Key: "$lookup", Value: bson.D{
							{Key: "from", Value: "locations"},
							{Key: "let", Value: bson.D{
								{Key: "placeIdStr", Value: bson.D{
									{Key: "$convert", Value: bson.D{
										{Key: "input", Value: "$placeId"},
										{Key: "to", Value: "objectId"},
										{Key: "onError", Value: nil},
										{Key: "onNull", Value: nil},
									}},
								}},
							}},
							{Key: "pipeline", Value: mongo.Pipeline{
								{{Key: "$match", Value: bson.D{
									{Key: "$expr", Value: bson.D{
										{Key: "$eq", Value: bson.A{"$_id", "$$placeIdStr"}},
									}},
								}}},
							}},
							{Key: "as", Value: "location"},
						},
					}},
					{{
						Key: "$lookup", Value: bson.D{
							{Key: "from", Value: "departments"},
							{Key: "let", Value: bson.D{
								{Key: "deptIdStr", Value: bson.D{
									{Key: "$convert", Value: bson.D{
										{Key: "input", Value: "$deptId"},
										{Key: "to", Value: "objectId"},
										{Key: "onError", Value: nil},
										{Key: "onNull", Value: nil},
									}},
								}},
							}},
							{Key: "pipeline", Value: mongo.Pipeline{
								{{Key: "$match", Value: bson.D{
									{Key: "$expr", Value: bson.D{
										{Key: "$eq", Value: bson.A{"$_id", "$$deptIdStr"}},
									}},
								}}},
							}},
							{Key: "as", Value: "department"},
						},
					}},
					{{
						Key: "$lookup", Value: bson.D{
							{Key: "from", Value: "asset_types"},
							{Key: "let", Value: bson.D{
								{Key: "typeIdStr", Value: bson.D{
									{Key: "$convert", Value: bson.D{
										{Key: "input", Value: "$typeId"},
										{Key: "to", Value: "objectId"},
										{Key: "onError", Value: nil},
										{Key: "onNull", Value: nil},
									}},
								}},
							}},
							{Key: "pipeline", Value: mongo.Pipeline{
								{{Key: "$match", Value: bson.D{
									{Key: "$expr", Value: bson.D{
										{Key: "$eq", Value: bson.A{"$_id", "$$typeIdStr"}},
									}},
								}}},
							}},
							{Key: "as", Value: "assetType"},
						},
					}},
					{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$location"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
					{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$department"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
					{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$assetType"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
				}},
				{Key: "as", Value: "assetlist"},
			},
		}},
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "locations"},

				{Key: "let", Value: bson.D{
					{Key: "locationIdStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$placeId"},
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
			},
		}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$location"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$assetlist"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
	}

	collectionItems := config.GetCollection("stock_take_items")

	cursor, err3 := collectionItems.Find(context.Background(), bson.M{"stockTakeId ": stockTakeId})
	if err3 != nil {
		return nil, err3
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err = collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func StockTakeItemSubmit(stockTakeItem *models.StockTakeItems) (interface{}, error) {

	collection := config.GetCollection("stock_take_items")

	stockTakeItem.CheckTime = time.Now()

	res, errInsert := collection.InsertOne(context.Background(), stockTakeItem)

	if errInsert != nil {
		return nil, errInsert
	}

	return res, nil
}

func FinishOrVoidStockTakeForm(_id string, status int, username string) (interface{}, error) {
	collection := config.GetCollection("stock_takes")
	objectID, err := primitive.ObjectIDFromHex(_id)

	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID, "status": 1}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingStockTake models.StockTakes
	err = collection.FindOne(ctx, filter).Decode(&existingStockTake)
	if err != nil {
		return nil, err
	}

	if existingStockTake.Status == 0 {
		return "Ooooops! This stock take form no longer active! Please create a new form!", nil
	} else if existingStockTake.Status == 2 {
		return "Ooooops! This stock take form has been completed! Please create a new form!", nil
	} else {
		finalData := bson.M{
			"finishBy":   username,
			"finishTime": time.Now(),
			"status":     status,
		}

		res, err2 := collection.UpdateOne(ctx, filter, bson.M{"$set": finalData})
		if err2 != nil {
			return nil, err2
		}

		return gin.H{"result": res, "message": "Updated successfully!", "finished": true}, nil
	}
}
