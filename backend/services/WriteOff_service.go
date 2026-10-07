package services

import (
	"context"
	"golang-fixedasset-mongo-backend/backend/config"
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"golang-fixedasset-mongo-backend/backend/tools"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func BatchCreateWriteOff(datas []dto.WriteOffPureList) (interface{}, error) {
	var results []interface{}
	for _, asset := range datas {

		assetList, _ := AssetDataFinderReturn(asset.AssetCode, asset.AssetName)
		location, _ := LocationDataFinder(asset.LastPlaceCode, asset.LastPlaceName)

		newData := models.WriteOffs{
			Reason:         asset.Reason,
			LastDay:        asset.LastDay,
			DisposalMethod: asset.DisposalMethod,
			RemainingValue: asset.RemainingValue,
			Status:         1,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if !assetList.ID.IsZero() {
			newData.AssetId = assetList.ID.Hex()
		}

		if !location.ID.IsZero() {
			newData.LastPlaceId = location.ID.Hex()
		}

		/*	result, err := CreateWriteOff(&newData)
			if err != nil {
				return nil, err
			} */
		results = append(results, newData)
	}
	return results, nil
}

func CreateWriteOff(data dto.CreateWriteOffRecrod) (interface{}, error) {
	var assetItem models.AssetLists

	assetCollection := config.GetCollection("asset_lists")
	assetCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	assetObjectID, err := primitive.ObjectIDFromHex(data.AssetId)

	if err != nil {
		return nil, err
	}

	_ = assetCollection.FindOne(assetCtx, bson.M{"_id": assetObjectID}).Decode(&assetItem)

	if assetItem.Status == 0 {
		CreateActionRecord("Write Off Create", "POST", "Write Off", data, "Failed")
		return "This asset item is already written off", nil
	} else {

		if data.LastDay == "" {
			data.LastDay = time.Now().Format("2006-01-02")
		}

		finalData := bson.M{
			"assetId":        data.AssetId,
			"status":         1,
			"reason":         data.Reason,
			"lastDay":        data.LastDay,
			"lastPlaceId":    data.LastPlaceId,
			"disposalMethod": data.DisposalMethod,
			"remainingValue": data.RemainingValue,
			"createdAt":      time.Now(),
			"updatedAt":      time.Now(),
		}

		WriteOffInactiveAsset(data.AssetId)

		CreateInvRecord(assetItem.AssetCode, data.LastPlaceId, "")

		collection := config.GetCollection("write_offs")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, errCreate := collection.InsertOne(ctx, finalData)

		if errCreate != nil {
			CreateActionRecord("Write Off Create", "POST", "Write Off", finalData, "Failed")
			return nil, errCreate
		}

		CreateActionRecord("Write Off Create", "POST", "Write Off", finalData, "Success")
		return res, nil
	}
}

func ListPageWriteOffWithFilter(dataReq dto.ListWriteOffReqDto) (interface{}, error) {
	collection := config.GetCollection("write_offs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filters := bson.M{"status": 1}

	if len(dataReq.DateRange) == 2 {
		filters["createdAt"] = bson.M{
			"$gte": dataReq.DateRange[0], // Ensure these are time.Time values
			"$lte": dataReq.DateRange[1],
		}
	}

	if len(dataReq.PlaceIds) > 0 {
		filters["lastPlaceId"] = bson.M{"$in": dataReq.PlaceIds}
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filters}},

		// Lookup Asset List
		bson.D{{Key: "$lookup", Value: bson.D{
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
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{
						{Key: "$eq", Value: bson.A{"$_id", "$$assetIdStr"}},
					}},
				}}},
			}},
			{Key: "as", Value: "assetlist"},
		}}},

		// Lookup Location
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "locations"},
			{Key: "let", Value: bson.D{
				{Key: "placeIdStr", Value: bson.D{
					{Key: "$convert", Value: bson.D{
						{Key: "input", Value: "$lastPlaceId"},
						{Key: "to", Value: "objectId"},
						{Key: "onError", Value: nil},
						{Key: "onNull", Value: nil},
					}},
				}},
			}},
			{Key: "pipeline", Value: mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{
						{Key: "$eq", Value: bson.A{"$_id", "$$placeIdStr"}},
					}},
				}}},
			}},
			{Key: "as", Value: "location"},
		}}},

		// Unwind joined arrays
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$assetlist"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$location"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},

		bson.D{{Key: "$sort", Value: bson.M{"createdAt": -1}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode directly into struct slice
	var FinalResults []bson.M
	if err := cursor.All(ctx, &FinalResults); err != nil {
		return nil, err
	}

	var LastResults []bson.M
	for _, item := range FinalResults {
		resultData, err := bson.Marshal(item)
		if err != nil {
			return nil, err
		}

		var writeOff bson.M
		if err = bson.Unmarshal(resultData, &writeOff); err != nil {
			return nil, err
		}

		writeOff["assetCode"] = tools.GetNestedString(item, "assetlist", "assetCode")
		writeOff["assetName"] = tools.GetNestedString(item, "assetlist", "assetName")
		writeOff["purchaseDate"] = tools.GetNestedString(item, "assetlist", "purchaseDate")
		writeOff["lastPlaceCode"] = tools.GetNestedString(item, "location", "placeCode")
		writeOff["lastPlaceName"] = tools.GetNestedString(item, "location", "placeName")

		LastResults = append(LastResults, writeOff)
	}

	return LastResults, nil
}

func ListPageWriteOff(dataReq dto.ListWriteOffReqDto) (interface{}, error) {
	collection := config.GetCollection("write_offs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filters := bson.M{"status": 1}
	if len(dataReq.DateRange) == 2 {
		filters["createdAt"] = bson.M{
			"$gte": dataReq.DateRange[0],
			"$lte": dataReq.DateRange[1],
		}
	}

	if len(dataReq.PlaceIds) > 0 {
		filters["lastPlaceId"] = bson.M{"$in": dataReq.PlaceIds}
	}

	skip := int64((dataReq.Page - 1) * dataReq.Limit)
	limit := int64(dataReq.Limit)

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filters}},
		bson.D{{Key: "$lookup", Value: bson.D{
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
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{
						{Key: "$eq", Value: bson.A{"$_id", "$$assetIdStr"}},
					}},
				}}},
			}},
			{Key: "as", Value: "assetlist"},
		}}},
		bson.D{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "locations"},
				{Key: "let", Value: bson.D{
					{Key: "placeIdStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$lastPlaceId"},
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
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$assetlist"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$location"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		bson.D{{Key: "$sort", Value: bson.M{"createdAt": -1}}},
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},
	}

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

	return gin.H{"lists": results, "total": count, "page": dataReq.Page, "limit": dataReq.Limit}, nil

}
