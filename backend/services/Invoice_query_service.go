package services

import (
	"context"
	"errors"
	"golang-pos-mongo-backend/backend/config"
	"golang-pos-mongo-backend/backend/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo bd Aggregation Set interface
type AggregationSet struct {
	Lookup  mongo.Pipeline
	Group   bson.M
	Project bson.M
	Sort    bson.D
}

// Safe conversion for string -> ObjectId
func toObjectID(field string) bson.M {
	return bson.M{
		"$convert": bson.M{
			"input":   field,
			"to":      "objectId",
			"onError": nil,
			"onNull":  nil,
		},
	}
}

// setting lookup for query
func buildLookup(
	collection string,
	localField string,
	alias string,
) bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.M{
			"from": collection,
			"let": bson.M{
				"convertedId": toObjectID(localField),
			},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr": bson.M{
						"$eq": bson.A{
							"$_id",
							"$$convertedId",
						},
					},
				}}},
			},
			"as": alias,
		}},
	}
}

// building unwind
func buildUnwind(path string) bson.D {
	return bson.D{
		{Key: "$unwind", Value: bson.M{
			"path":                       "$" + path,
			"preserveNullAndEmptyArrays": true,
		}},
	}
}

// get filters setup
func GetFilter(
	query dto.DashboardReqFilterDto,
) bson.M {

	filters := bson.M{}

	// createdAt != null
	dateFilter := bson.M{
		"$ne": nil,
	}

	if len(query.SalesDateRange) >= 2 {
		dateFilter["$gte"] = query.SalesDateRange[0]
		dateFilter["$lte"] = query.SalesDateRange[1]
	}

	filters["invoice.createdAt"] = dateFilter

	if len(query.TypeIds) > 0 {
		filters["product.typeId"] = bson.M{
			"$in": query.TypeIds,
		}
	}

	if len(query.PlaceIds) > 0 {
		filters["invoice.locationId"] = bson.M{
			"$in": query.PlaceIds,
		}
	}

	if len(query.DeptIds) > 0 {
		filters["product.deptId"] = bson.M{
			"$in": query.DeptIds,
		}
	}

	if query.ProductCode != "" {
		filters["product.productCode"] = bson.M{
			"$regex":   query.ProductCode,
			"$options": "i",
		}
	}

	if query.ProductName != "" {
		filters["product.productName"] = bson.M{
			"$regex":   query.ProductName,
			"$options": "i",
		}
	}

	return filters
}

// get department set
func GetDeptSet() AggregationSet {
	return AggregationSet{
		Lookup: mongo.Pipeline{
			buildLookup(
				"departments",
				"$product.deptId",
				"department",
			),
			buildUnwind("department"),
		},
		Group: bson.M{
			"_id": "$department.deptName",
		},
		Project: bson.M{
			"deptName": "$_id",
		},
		Sort: bson.D{
			{Key: "deptName", Value: 1},
		},
	}
}

// get set product type
func GetTypeSet() AggregationSet {
	return AggregationSet{
		Lookup: mongo.Pipeline{
			buildLookup(
				"product_types",
				"$product.typeId",
				"producttype",
			),
			buildUnwind("producttype"),
		},
		Group: bson.M{
			"_id": "$producttype.typeName",
		},
		Project: bson.M{
			"typeName": "$_id",
		},
		Sort: bson.D{
			{Key: "typeName", Value: 1},
		},
	}
}

// get set location
func GetLocationSet() AggregationSet {
	return AggregationSet{
		Lookup: mongo.Pipeline{
			buildLookup(
				"locations",
				"$invoice.locationId",
				"location",
			),
			buildUnwind("location"),
		},
		Group: bson.M{
			"_id": "$location.placeName",
		},
		Project: bson.M{
			"placeName": "$_id",
		},
		Sort: bson.D{
			{Key: "placeName", Value: 1},
		},
	}
}

// get year & month
func GetYearMonth() AggregationSet {
	return AggregationSet{
		Lookup: mongo.Pipeline{},

		Group: bson.M{
			"_id": bson.M{
				"year": bson.M{
					"$year": "$invoice.createdAt",
				},
				"month": bson.M{
					"$dateToString": bson.M{
						"format": "%B",
						"date":   "$invoice.createdAt",
					},
				},
				"monthNum": bson.M{
					"$month": "$invoice.createdAt",
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

// running query
func InvoiceGetByDataType(
	ctx context.Context,
	query dto.InvoiceDashboardReqDto,
) ([]bson.M, error) {

	if query.ValueField == "" {
		return nil, errors.New("Value Field is required")
	}

	if query.DataTypeValue == "" {
		return nil, errors.New("Data Type is required")
	}

	var dataTypeObj AggregationSet

	// Select Data Type
	switch query.DataTypeValue {
	case "dept":
		dataTypeObj = GetDeptSet()

	case "type":
		dataTypeObj = GetTypeSet()

	case "location":
		dataTypeObj = GetLocationSet()

	case "year-month":
		dataTypeObj = GetYearMonth()

	default:
		return nil, errors.New("Invalid Data Type")
	}

	// Select Value Field
	var valueFieldObj bson.M

	switch query.ValueField {
	case "qtys":
		valueFieldObj = bson.M{
			"qtys": bson.M{
				"$sum": "$qty",
			},
		}

	case "price":
		valueFieldObj = bson.M{
			"price": bson.M{
				"$sum": "$price",
			},
		}

	default:
		return nil, errors.New("Invalid Value Field")
	}

	filters := GetFilter(query.Filter)

	// Base aggregation pipeline
	pipeline := mongo.Pipeline{
		buildLookup(
			"invoices",
			"$invoiceId",
			"invoice",
		),

		buildUnwind("invoice"),

		buildLookup(
			"product_lists",
			"$productId",
			"product",
		),

		buildUnwind("product"),
	}

	// Append dynamic lookup
	pipeline = append(
		pipeline,
		dataTypeObj.Lookup...,
	)

	// Match filters
	pipeline = append(pipeline,
		bson.D{
			{Key: "$match", Value: filters},
		},
	)

	// Merge group fields
	groupFields := bson.M{}

	for k, v := range dataTypeObj.Group {
		groupFields[k] = v
	}

	for k, v := range valueFieldObj {
		groupFields[k] = v
	}

	pipeline = append(pipeline,
		bson.D{
			{Key: "$group", Value: groupFields},
		},
	)

	// Merge project fields
	projectFields := bson.M{
		"_id": 0,
	}

	for k, v := range dataTypeObj.Project {
		projectFields[k] = v
	}

	// Include selected value field
	projectFields[query.ValueField] = 1

	pipeline = append(pipeline,
		bson.D{
			{Key: "$project", Value: projectFields},
		},
	)

	// Sorting
	pipeline = append(pipeline,
		bson.D{
			{Key: "$sort", Value: dataTypeObj.Sort},
		},
	)

	InvoiceItemCollection := config.GetCollection("invoice_items")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute
	cursor, err := InvoiceItemCollection.Aggregate(
		ctx,
		pipeline,
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	// Decode aggregation results
	results := make([]bson.M, 0)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
