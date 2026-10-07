package services

import (
	"context"
	"fmt"
	"golang-fixedasset-mongo-backend/backend/config"
	"golang-fixedasset-mongo-backend/backend/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupConfig struct {
	Lookup  bson.D
	Unwind  bson.D
	GroupBy bson.D
	Project bson.D
	Sort    bson.D
}

// ---------------- Filter Builder ----------------
func getFilter(filter *dto.DashboardReqFilterDto) bson.D {
	f := bson.D{}
	if len(filter.PurchaseDates) == 2 {
		start := filter.PurchaseDates[0]
		end := filter.PurchaseDates[1]
		f = append(f, bson.E{Key: "purchaseDate", Value: bson.D{{Key: "$gte", Value: start}, {Key: "$lte", Value: end}}})
	}
	if len(filter.DeptIds) > 0 {
		f = append(f, bson.E{Key: "deptId", Value: bson.D{{Key: "$in", Value: filter.DeptIds}}})
	}
	if len(filter.TypeIds) > 0 {
		f = append(f, bson.E{Key: "typeId", Value: bson.D{{Key: "$in", Value: filter.TypeIds}}})
	}
	if len(filter.PlaceIds) > 0 {
		f = append(f, bson.E{Key: "placeId", Value: bson.D{{Key: "$in", Value: filter.PlaceIds}}})
	}
	return f
}

// ---------------- GroupBy / Lookup ----------------
func getGroupByDept() (lookup, unwind, groupBy, project bson.D) {
	lookup = bson.D{
		{Key: "from", Value: "departments"},
		{Key: "let", Value: bson.D{{Key: "deptIdStr", Value: bson.D{{Key: "$toObjectId", Value: "$deptId"}}}}},
		{Key: "pipeline", Value: mongo.Pipeline{
			{{Key: "$match", Value: bson.D{{Key: "$expr", Value: bson.D{{Key: "$eq", Value: bson.A{"$_id", "$$deptIdStr"}}}}}}},
		}},
		{Key: "as", Value: "department"},
	}
	unwind = bson.D{{Key: "path", Value: "$department"}, {Key: "preserveNullAndEmptyArrays", Value: true}}
	groupBy = bson.D{{Key: "_id", Value: bson.D{{Key: "deptName", Value: "$department.deptName"}}}}
	project = bson.D{{Key: "deptName", Value: "$_id.deptName"}}
	return lookup, unwind, groupBy, project
}

func getGroupByType() (lookup, unwind, groupBy, project bson.D) {
	lookup = bson.D{
		{Key: "from", Value: "asset_types"},
		{Key: "let", Value: bson.D{{Key: "typeIdStr", Value: bson.D{{Key: "$toObjectId", Value: "$typeId"}}}}},
		{Key: "pipeline", Value: mongo.Pipeline{
			{{Key: "$match", Value: bson.D{{Key: "$expr", Value: bson.D{{Key: "$eq", Value: bson.A{"$_id", "$$typeIdStr"}}}}}}},
		}},
		{Key: "as", Value: "assettype"},
	}
	unwind = bson.D{{Key: "path", Value: "$assettype"}, {Key: "preserveNullAndEmptyArrays", Value: true}}
	groupBy = bson.D{{Key: "_id", Value: bson.D{{Key: "typeName", Value: "$assettype.typeName"}}}}
	project = bson.D{{Key: "typeName", Value: "$_id.typeName"}}
	return lookup, unwind, groupBy, project
}

func getGroupByLocations() (lookup, unwind, groupBy, project bson.D) {
	lookup = bson.D{
		{Key: "from", Value: "locations"},
		{Key: "let", Value: bson.D{{Key: "placeIdStr", Value: bson.D{{Key: "$toObjectId", Value: "$placeId"}}}}},
		{Key: "pipeline", Value: mongo.Pipeline{
			{{Key: "$match", Value: bson.D{{Key: "$expr", Value: bson.D{{Key: "$eq", Value: bson.A{"$_id", "$$placeIdStr"}}}}}}},
		}},
		{Key: "as", Value: "location"},
	}
	unwind = bson.D{{Key: "path", Value: "$location"}, {Key: "preserveNullAndEmptyArrays", Value: true}}
	groupBy = bson.D{{Key: "_id", Value: bson.D{{Key: "placeName", Value: "$location.placeName"}}}}
	project = bson.D{{Key: "placeName", Value: "$_id.placeName"}}
	return lookup, unwind, groupBy, project
}

func getGroupByYearMonth() (addFields, groupBy, project, sort bson.D) {
	// new field for purchaseDateObj firstly
	addFields = bson.D{
		{Key: "$addFields", Value: bson.D{
			{Key: "purchaseDateObj", Value: bson.D{
				{Key: "$dateFromString", Value: bson.D{
					{Key: "dateString", Value: "$purchaseDate"},
				}},
			}},
		}},
	}

	// new field for groupBy
	groupBy = bson.D{
		{Key: "_id", Value: bson.D{
			{Key: "year", Value: bson.D{{Key: "$year", Value: "$purchaseDateObj"}}},
			{Key: "month", Value: bson.D{{Key: "$month", Value: "$purchaseDateObj"}}},
		}},
	}

	sort = bson.D{
		{Key: "_id.year", Value: 1},
		{Key: "_id.month", Value: 1},
	}

	project = bson.D{
		{Key: "year", Value: bson.D{{Key: "$toString", Value: "$_id.year"}}},
		{Key: "month", Value: bson.D{{Key: "$toString", Value: "$_id.month"}}},
		{Key: "monthString", Value: bson.D{
			{Key: "$arrayElemAt", Value: bson.A{
				bson.A{"", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
				"$_id.month",
			}},
		}},
	}

	return addFields, groupBy, project, sort
}

type UnitField struct {
	Group bson.D
	Proj  bson.D
}

// 支援同時返回 costs + count
func getGroupByUnit() map[string]UnitField {
	return map[string]UnitField{
		"costs": {Group: bson.D{{Key: "costs", Value: bson.D{{Key: "$sum", Value: "$cost"}}}}, Proj: bson.D{{Key: "costs", Value: 1}}},
		"count": {Group: bson.D{{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}, Proj: bson.D{{Key: "count", Value: 1}}},
		"all": {
			Group: bson.D{
				{Key: "costs", Value: bson.D{{Key: "$sum", Value: "$cost"}}},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			},
			Proj: bson.D{
				{Key: "costs", Value: 1},
				{Key: "count", Value: 1},
			},
		},
	}
}

func QueryMakerForData(query *dto.DashboardReqDto) ([]bson.M, error) {
	collection := config.GetCollection("asset_lists")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var pipeline mongo.Pipeline

	// ---------------- Filter ----------------
	matchFilter := bson.D{{Key: "status", Value: 1}}
	if query.Filter != nil {
		filters := getFilter(query.Filter)
		matchFilter = append(matchFilter, filters...)
	}
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: matchFilter}})

	// ---------------- Date Conversion ----------------
	// 將 purchaseDate (string) 轉成 Date
	if query.DateType && query.DateTypeValue == "YearMonth" {
		addFields := bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "purchaseDateObj", Value: bson.D{
					{Key: "$dateFromString", Value: bson.D{{Key: "dateString", Value: "$purchaseDate"}}},
				}},
			}},
		}
		pipeline = append(pipeline, addFields)
	}

	// ---------------- DataType ----------------
	var lookup, unwind, dataGroupBy, dataProject bson.D
	if query.DataType {
		switch query.DataTypeValue {
		case "dept":
			lookup, unwind, dataGroupBy, dataProject = getGroupByDept()
		case "type":
			lookup, unwind, dataGroupBy, dataProject = getGroupByType()
		case "location":
			lookup, unwind, dataGroupBy, dataProject = getGroupByLocations()
		}
		pipeline = append(pipeline, bson.D{{Key: "$lookup", Value: lookup}})
		pipeline = append(pipeline, bson.D{{Key: "$unwind", Value: unwind}})
	}

	// ---------------- DateType ----------------
	var dateGroupBy, dateProject, dateSort bson.D
	if query.DateType {
		switch query.DateTypeValue {
		case "YearMonth":
			_, dateGroupBy, dateProject, dateSort = getGroupByYearMonth()
		}
	}

	// ---------------- Value Field ----------------
	unitMap := getGroupByUnit()
	unitKey := query.ValueField
	switch unitKey {
	case "counts":
		unitKey = "count"
	case "", "all":
		unitKey = "all"
	}
	unitFields, ok := unitMap[unitKey]
	if !ok {
		return nil, fmt.Errorf("invalid valueField: %s", unitKey)
	}

	// ---------------- Group Stage ----------------
	groupID := bson.D{}
	if len(dateGroupBy) > 0 {
		if idField, ok := dateGroupBy.Map()["_id"].(bson.D); ok {
			groupID = append(groupID, idField...)
		}
	}
	if len(dataGroupBy) > 0 {
		if idField, ok := dataGroupBy.Map()["_id"].(bson.D); ok {
			groupID = append(groupID, idField...)
		}
	}

	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: groupID},
	}}}

	// add value fields
	for _, elem := range unitFields.Group {
		groupStage[0].Value = append(groupStage[0].Value.(bson.D), elem)
	}

	pipeline = append(pipeline, groupStage)

	// ---------------- Project Stage ----------------
	projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0}}}}
	for _, elem := range dateProject {
		projectStage[0].Value = append(projectStage[0].Value.(bson.D), elem)
	}
	for _, elem := range dataProject {
		projectStage[0].Value = append(projectStage[0].Value.(bson.D), elem)
	}
	for _, elem := range unitFields.Proj {
		projectStage[0].Value = append(projectStage[0].Value.(bson.D), elem)
	}
	pipeline = append(pipeline, projectStage)

	// ---------------- Sort Stage ----------------
	if len(dateSort) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: dateSort}})
	}

	// ---------------- Execute ----------------
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}
