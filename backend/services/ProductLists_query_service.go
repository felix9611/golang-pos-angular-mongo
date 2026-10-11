package services

import (
	"context"
	"errors"
	"golang-fixedasset-mongo-backend/backend/config"
	"golang-fixedasset-mongo-backend/backend/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetProductLocationFilter(req dto.DashboardReqFilterDto) (interface{}, error) {
	filters := bson.M{}

	if len(req.DeptIds) > 0 {
		filters["deptId"] = bson.M{"$in": req.DeptIds}
	}
	if len(req.TypeIds) > 0 {
		filters["typeId"] = bson.M{"$in": req.TypeIds}
	}
	if len(req.PlaceIds) > 0 {
		filters["placeId"] = bson.M{"$in": req.PlaceIds}
	}
	return filters, nil
}

func ProductLocationGetLocationSet() AggregationSet {
	return AggregationSet{
		Lookup: mongo.Pipeline{
			buildLookup(
				"locations",
				"$locationId",
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

func ProductLocationGetByDataType(
	ctx context.Context,
	query dto.InvoiceDashboardReqDto,
) ([]bson.M, error) {

	if query.DataTypeValue == "" {
		return nil, errors.New("Data Type is required")
	}

	var dataTypeObj AggregationSet

	switch query.DataTypeValue {
	case "location":
		dataTypeObj = ProductLocationGetLocationSet()
	case "dept":
		dataTypeObj = GetDeptSet()
	case "type":
		dataTypeObj = GetTypeSet()
	default:
		return nil, errors.New("Invalid Data Type")
	}

	valueFieldObj := bson.M{
		"qtys": bson.M{
			"$sum": "$qty",
		},
		"cost": bson.M{
			"$sum": "$totalCost",
		},
		"price": bson.M{
			"$sum": "$totalPrice",
		},
	}

	filters, err := GetProductLocationFilter(query.Filter)

	pipeline := mongo.Pipeline{
		buildLookup(
			"product_lists",
			"$productId",
			"product",
		),
		buildUnwind("product"),
	}

	// Dynamic lookup
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

	// Group
	groupFields := bson.M{}

	// Include _id
	for k, v := range dataTypeObj.Group {
		groupFields[k] = v
	}

	// Include qtys, cost, price
	for k, v := range valueFieldObj {
		groupFields[k] = v
	}

	pipeline = append(pipeline,
		bson.D{
			{Key: "$group", Value: groupFields},
		},
	)

	// Project
	projectFields := bson.M{
		"_id":   0,
		"qtys":  1,
		"cost":  1,
		"price": 1,
	}

	for k, v := range dataTypeObj.Project {
		projectFields[k] = v
	}

	pipeline = append(pipeline,
		bson.D{
			{Key: "$project", Value: projectFields},
		},
	)

	// Sort
	pipeline = append(pipeline,
		bson.D{
			{Key: "$sort", Value: dataTypeObj.Sort},
		},
	)

	collection := config.GetCollection("product_location")

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make([]bson.M, 0)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
