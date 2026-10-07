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
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateInvRecord(assetCode string, placeFrom string, placeTo string) (interface{}, error) {
	collection := config.GetCollection("inv_records")

	body := models.InvRecords{
		AssetCode: assetCode,
		PlaceFrom: placeFrom,
		PlaceTo:   placeTo,
		CreatedAt: time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), body)
	if err != nil {
		CreateActionRecord("Inventory Record Create", "POST", "Inventory Record", body, "Failed")
		return nil, err
	}

	CreateActionRecord("Inventory Record Create", "POST", "Inventory Record", body, "Success")

	return body, nil
}

func ListInvRecordsWithFilter(dtoData *dto.ListRecordReqDto) (interface{}, error) {
	collection := config.GetCollection("inv_records")

	filters := bson.M{}

	if dtoData.AssetCode != "" {
		filters["asset_code"] = dtoData.AssetCode
	}

	if len(dtoData.DateRange) > 0 {
		filters["createdAt"] = bson.M{"$gte": dtoData.DateRange[0], "$lte": dtoData.DateRange[1]}
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filters}},

		// lookup asset item data
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "asset_lists"},
				{Key: "localField", Value: "assetCode"},
				{Key: "foreignField", Value: "assetCode"},
				{Key: "as", Value: "assetlist"},
			},
		}},

		// $lookup place from
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "locations"},
				{Key: "let", Value: bson.D{
					{Key: "placeFromStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$placeFrom"},
							{Key: "to", Value: "objectId"},
							{Key: "onError", Value: nil},
							{Key: "onNull", Value: nil},
						}},
					}},
				}},
				{Key: "pipeline", Value: mongo.Pipeline{
					{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$_id", "$$placeFromStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "placeFromData"},
			},
		}},

		// $lookup place to
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "locations"},
				{Key: "let", Value: bson.D{
					{Key: "placeToStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$placeTo"},
							{Key: "to", Value: "objectId"},
							{Key: "onError", Value: nil},
							{Key: "onNull", Value: nil},
						}},
					}},
				}},
				{Key: "pipeline", Value: mongo.Pipeline{
					{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$_id", "$$placeToStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "placeToData"},
			},
		}},

		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$assetlist"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$placeFromData"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$placeToData"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},

		// sort
		{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: 1}}}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var FinalResults []bson.M
	if err := cursor.All(ctx, &FinalResults); err != nil {
		return nil, err
	}

	var LastResults []dto.InvRecordPureList
	for _, item := range FinalResults {
		resultData, err := bson.Marshal(item)
		if err != nil {
			return nil, err
		}

		var assetDetails dto.InvRecordPureList
		if err := bson.Unmarshal(resultData, &assetDetails); err != nil {
			return nil, err
		}

		assetDetails.PlaceFromName = tools.GetNestedString(item, "placeFromData", "placeName")
		assetDetails.PlaceFromCode = tools.GetNestedString(item, "placeFromData", "placeCode")
		assetDetails.PlaceToCode = tools.GetNestedString(item, "placeToData", "placeCode")
		assetDetails.PlaceToName = tools.GetNestedString(item, "placeToData", "placeName")
		assetDetails.AssetName = tools.GetNestedString(item, "assetlist", "assetName")
		assetDetails.AssetCode = tools.GetNestedString(item, "assetlist", "assetCode")

		LastResults = append(LastResults, assetDetails)
	}

	return LastResults, nil

}

func ListInvRecords(dtoData *dto.ListRecordReqDto) (interface{}, error) {
	collection := config.GetCollection("inv_records")

	filters := bson.M{}

	if dtoData.AssetCode != "" {
		filters["asset_code"] = dtoData.AssetCode
	}

	if len(dtoData.DateRange) > 0 {
		filters["createdAt"] = bson.M{"$gte": dtoData.DateRange[0], "$lte": dtoData.DateRange[1]}
	}

	skip := int64((dtoData.Page - 1) * dtoData.Limit)
	limit := int64(dtoData.Limit)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filters}},

		// lookup asset item data
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "asset_lists"},
				{Key: "localField", Value: "assetCode"},
				{Key: "foreignField", Value: "assetCode"},
				{Key: "as", Value: "assetlist"},
			},
		}},

		// $lookup place from
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "locations"},
				{Key: "let", Value: bson.D{
					{Key: "placeFromStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$placeFrom"},
							{Key: "to", Value: "objectId"},
							{Key: "onError", Value: nil},
							{Key: "onNull", Value: nil},
						}},
					}},
				}},
				{Key: "pipeline", Value: mongo.Pipeline{
					{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$_id", "$$placeFromStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "placeFromData"},
			},
		}},

		// $lookup place to
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "locations"},
				{Key: "let", Value: bson.D{
					{Key: "placeToStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$placeTo"},
							{Key: "to", Value: "objectId"},
							{Key: "onError", Value: nil},
							{Key: "onNull", Value: nil},
						}},
					}},
				}},
				{Key: "pipeline", Value: mongo.Pipeline{
					{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$_id", "$$placeToStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "placeToData"},
			},
		}},

		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$assetlist"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$placeFromData"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$placeToData"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},

		// sort
		{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: 1}}}},

		// pagination
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

	return gin.H{"lists": results, "total": count, "page": dtoData.Page, "limit": dtoData.Limit}, nil

}
