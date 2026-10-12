package services

import (
	"context"
	"fmt"
	"golang-pos-mongo-backend/backend/config"
	"golang-pos-mongo-backend/backend/dto"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetInvRecordFilter(
	query dto.DashboardReqFilterDto,
) (bson.M, error) {

	filters := bson.M{}

	// Date Range
	if len(query.SalesDateRange) >= 2 {
		from, err := time.Parse(
			time.RFC3339,
			query.SalesDateRange[0],
		)
		if err != nil {
			return nil, err
		}

		to, err := time.Parse(
			time.RFC3339,
			query.SalesDateRange[1],
		)
		if err != nil {
			return nil, err
		}

		filters["createdAt"] = bson.M{
			"$gte": from,
			"$lte": to,
		}
	}

	// PlaceIds - OR condition
	if len(query.PlaceIds) > 0 {
		filters["$or"] = bson.A{
			bson.M{
				"locFrom": bson.M{
					"$in": query.PlaceIds,
				},
			},
			bson.M{
				"locTo": bson.M{
					"$in": query.PlaceIds,
				},
			},
		}
	}

	// Product Code - case insensitive regex
	if strings.TrimSpace(query.ProductCode) != "" {
		filters["product.productCode"] = bson.M{
			"$regex":   regexp.QuoteMeta(query.ProductCode),
			"$options": "i",
		}
	}

	return filters, nil
}

func GetInvRecordYearMonth() AggregationSet {
	return AggregationSet{
		Group: bson.M{
			"_id": bson.M{
				"year": bson.M{
					"$year": "$createdAt",
				},
				"month": bson.M{
					"$dateToString": bson.M{
						"format": "%B",
						"date":   "$createdAt",
					},
				},
				"monthNum": bson.M{
					"$month": "$createdAt",
				},
			},
		},

		Project: bson.M{
			"year":     "$_id.year",
			"month":    "$_id.month",
			"monthNum": "$_id.monthNum",
		},

		Sort: bson.D{
			{Key: "year", Value: 1},
			{Key: "monthNum", Value: 1},
		},
	}
}

func GetInvRecordQueryGroupBy() bson.M {
	return bson.M{
		"qtys": bson.M{
			"$sum": "$qty",
		},
		"costs": bson.M{
			"$sum": "$cost",
		},
	}
}

func GetInvRecordQueryProject() bson.M {
	return bson.M{
		"qtys":  1,
		"costs": 1,
	}
}

func GetInvRecordByDataType(
	ctx context.Context,
	query dto.InvoiceDashboardReqDto,
) (interface{}, error) {
	InventoryRecordsCollection := config.GetCollection("inventory_records")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if query.DataTypeValue == "" {
		return gin.H{"error": "dataType is required"}, nil
	}

	// Data Type Filter
	dataTypeFilter := bson.M{}

	switch query.DataTypeValue {

	case "stockIn":
		dataTypeFilter["locFrom"] = ""

	case "stockOut":
		dataTypeFilter["locTo"] = ""

	case "stockMove":
		dataTypeFilter["locFrom"] = bson.M{
			"$ne": "",
		}
		dataTypeFilter["locTo"] = bson.M{
			"$ne": "",
		}

	default:
		return nil, fmt.Errorf(
			"invalid dataType: %s",
			query.DataTypeValue,
		)
	}

	// Common Filter
	filterObj, err := GetInvRecordFilter(query.Filter)
	if err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	// Merge filters
	for key, value := range dataTypeFilter {
		filterObj[key] = value
	}

	// Year Month
	yearMonthObj := GetInvRecordYearMonth()

	// Group
	groupObj := GetInvRecordQueryGroupBy()

	for key, value := range yearMonthObj.Group {
		groupObj[key] = value
	}

	// Project
	projectObj := GetInvRecordQueryProject()

	for key, value := range yearMonthObj.Project {
		projectObj[key] = value
	}

	// Aggregation Pipeline
	pipeline := mongo.Pipeline{

		// $lookup products
		{{
			Key: "$lookup",
			Value: bson.M{
				"from": "product_lists",

				"let": bson.M{
					"productIdStr": bson.M{
						"$convert": bson.M{
							"input":   "$productId",
							"to":      "objectId",
							"onError": nil,
							"onNull":  nil,
						},
					},
				},

				"pipeline": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{
								"$eq": bson.A{
									"$_id",
									"$$productIdStr",
								},
							},
						},
					},
				},

				"as": "product",
			},
		}},

		// $unwind product
		{{
			Key: "$unwind",
			Value: bson.M{
				"path":                       "$product",
				"preserveNullAndEmptyArrays": true,
			},
		}},

		// $match
		{{
			Key:   "$match",
			Value: filterObj,
		}},

		// $group
		{{
			Key:   "$group",
			Value: groupObj,
		}},

		// $project
		{{
			Key:   "$project",
			Value: projectObj,
		}},

		// $sort
		{{
			Key:   "$sort",
			Value: yearMonthObj.Sort,
		}},
	}

	cursor, err := InventoryRecordsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode Result
	results := make([]dto.InvRecordYearMonthResult, 0)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
