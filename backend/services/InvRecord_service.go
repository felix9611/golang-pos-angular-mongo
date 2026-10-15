package services

import (
	"context"
	"golang-pos-mongo-backend/backend/config"
	"golang-pos-mongo-backend/backend/dto"
	"golang-pos-mongo-backend/backend/models"
	"golang-pos-mongo-backend/backend/tools"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ListInvRecordsWithFilter(dtoData *dto.InvRecordListDto) ([]dto.InvRecordPureList, error) {
	collection := config.GetCollection("inventory_records")

	filters := bson.M{}

	if len(dtoData.Date) > 0 {
		filters["createdAt"] = bson.M{"$gte": dtoData.Date[0], "$lte": dtoData.Date[1]}
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filters}},

		// $lookup place from
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "locations"},
				{Key: "let", Value: bson.D{
					{Key: "placeFromStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$locFrom"},
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
				{Key: "as", Value: "locFromData"},
			},
		}},

		// $lookup place to
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "locations"},
				{Key: "let", Value: bson.D{
					{Key: "placeToStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$locTo"},
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
				{Key: "as", Value: "locToData"},
			},
		}},

		// $lookup product
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "product_lists"},
				{Key: "let", Value: bson.D{
					{Key: "productToStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$ProductId"},
							{Key: "to", Value: "objectId"},
							{Key: "onError", Value: nil},
							{Key: "onNull", Value: nil},
						}},
					}},
				}},
				{Key: "pipeline", Value: mongo.Pipeline{
					{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$_id", "$$productToStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "product"},
			},
		}},

		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$product"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$locFromData"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$locToData"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},

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

		var InvRecordDetails dto.InvRecordPureList
		if err := bson.Unmarshal(resultData, &InvRecordDetails); err != nil {
			return nil, err
		}

		InvRecordDetails.ProductCode = tools.GetNestedValue(item, "product", "productCode").(string)
		InvRecordDetails.ProductName = tools.GetNestedValue(item, "product", "productName").(string)
		InvRecordDetails.Unit = tools.GetNestedValue(item, "product", "unit").(string)
		InvRecordDetails.ItemCode = tools.GetNestedValue(item, "product", "itemCode").(string)
		InvRecordDetails.BrandCode = tools.GetNestedValue(item, "product", "brandCode").(string)
		InvRecordDetails.BrandName = tools.GetNestedValue(item, "product", "brandName").(string)
		InvRecordDetails.CostPrice = tools.GetFloat64(item, "product", "costPrice").(float64)
		InvRecordDetails.RetailPrice = tools.GetFloat64(item, "product", "retailPrice").(float64)

		if item["locFromData"] != nil {
			InvRecordDetails.LocFromCode = tools.GetNestedValue(item, "locFromData", "placeCode").(string)
			InvRecordDetails.LocFromName = tools.GetNestedValue(item, "locFromData", "placeName").(string)
		}

		if item["locToData"] != nil {
			InvRecordDetails.LocToCode = tools.GetNestedValue(item, "locToData", "placeCode").(string)
			InvRecordDetails.LocToName = tools.GetNestedValue(item, "locToData", "placeName").(string)
		}

		LastResults = append(LastResults, InvRecordDetails)
	}

	return LastResults, nil
}

func ListInvRecords(req *dto.InvRecordListDto) (interface{}, error) {
	collection := config.GetCollection("inventory_records")

	filters := bson.M{}

	if len(req.Date) > 0 {
		filters["createdAt"] = bson.M{"$gte": req.Date[0], "$lte": req.Date[1]}
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	skip := int64((req.Page - 1) * req.Limit)
	limit := int64(req.Limit)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filters}},

		// $lookup place from
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "locations"},
				{Key: "let", Value: bson.D{
					{Key: "placeFromStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$locFrom"},
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
				{Key: "as", Value: "locFromData"},
			},
		}},

		// $lookup place to
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "locations"},
				{Key: "let", Value: bson.D{
					{Key: "placeToStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$locTo"},
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
				{Key: "as", Value: "locToData"},
			},
		}},

		// $lookup product
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "product_lists"},
				{Key: "let", Value: bson.D{
					{Key: "productToStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$ProductId"},
							{Key: "to", Value: "objectId"},
							{Key: "onError", Value: nil},
							{Key: "onNull", Value: nil},
						}},
					}},
				}},
				{Key: "pipeline", Value: mongo.Pipeline{
					{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$_id", "$$productToStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "product"},
			},
		}},

		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$product"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$locFromData"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$locToData"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},

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

	var results []dto.InvRecordList
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	count, errCount := collection.CountDocuments(ctx, filters)
	if errCount != nil {
		return nil, errCount
	}

	return dto.InvRecordListResponse{Lists: results, Total: count, Page: req.Page, Limit: req.Limit}, nil

}

func InsertInventoryRecord(dto dto.InsertInvRecordDto) (interface{}, error) {

	collection := config.GetCollection("inventory_records")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	finalData := models.InventoryRecords{
		ProductId: dto.ProductId,
		LocFrom:   dto.LocFrom,
		LocTo:     dto.LocTo,
		Qty:       dto.Qty,
		Cost:      dto.Cost,
		CreatedAt: time.Now(),
	}

	_, err := collection.InsertOne(ctx, finalData)
	if err != nil {
		return nil, err
	}

	return finalData, nil
}
