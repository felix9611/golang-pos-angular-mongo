package services

import (
	"context"
	"fmt"
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

func BatchCreateAssetItems(assetItems []dto.AssetListsPureDetails) (interface{}, error) {
	var results []interface{}
	for _, asset := range assetItems {

		assetType, _ := AssetTypeDataFinderReturn(asset.TypeCode, asset.TypeName)

		department, _ := DepartmentDataFinder(asset.DeptCode, asset.DeptName)

		location, _ := LocationDataFinder(asset.PlaceCode, asset.PlaceName)

		var sponsor bool
		if asset.Sponsor == "Yes" || asset.Sponsor == "YES" || asset.Sponsor == "yes" {
			sponsor = true
		} else {
			sponsor = false
		}

		var includeTax bool
		if asset.IncludeTax == "Yes" || asset.IncludeTax == "YES" || asset.IncludeTax == "yes" {
			includeTax = true
		} else {
			includeTax = false
		}

		newData := models.AssetLists{
			AssetCode:              asset.AssetCode,
			AssetName:              asset.AssetName,
			Unit:                   asset.Unit,
			PurchaseDate:           asset.PurchaseDate,
			Description:            asset.Description,
			Sponsor:                sponsor,
			SponsorName:            asset.SponsorName,
			Cost:                   asset.Cost,
			SerialNumber:           asset.SerialNumber,
			InvoiceNo:              asset.InvoiceNo,
			InvoiceDate:            asset.InvoiceDate,
			InvoiceRemark:          asset.InvoiceRemark,
			VendorId:               asset.VendorId,
			Remark:                 asset.Remark,
			TaxCountryCode:         asset.TaxCountryCode,
			TaxCode:                asset.TaxCode,
			TaxRate:                asset.TaxRate,
			IncludeTax:             includeTax,
			AfterBeforeTax:         asset.AfterBeforeTax,
			AccountCode:            asset.AccountCode,
			AccountName:            asset.AccountName,
			BrandCode:              asset.BrandCode,
			BrandName:              asset.BrandName,
			ChequeNo:               asset.ChequeNo,
			MaintenancePeriodStart: asset.MaintenancePeriodStart,
			MaintenancePeriodEnd:   asset.MaintenancePeriodEnd,
			VoucherNo:              asset.VoucherNo,
			VoucherUsedDate:        asset.VoucherUsedDate,
			StaffName:              asset.StaffName,
			Status:                 1,
			CreatedAt:              time.Now().Format(time.RFC3339),
			UpdatedAt:              time.Now().Format(time.RFC3339),
		}

		if !assetType.ID.IsZero() {
			newData.TypeID = assetType.ID.Hex()
		}

		if !department.ID.IsZero() {
			newData.DeptId = department.ID.Hex()
		}

		if !location.ID.IsZero() {
			newData.PlaceId = location.ID.Hex()
		}

		result, err := CreateAssetItem(&newData)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func GetOneAssetItemByID(id string) (interface{}, error) {
	collection := config.GetCollection("asset_lists")
	var assetItem models.AssetLists

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err2 := collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&assetItem)
	if err2 != nil {
		return nil, err2
	}
	if assetItem.Status == 0 {
		return "This Asset item was write off", nil
	} else {

		return &assetItem, nil
	}
}

func GetOneAssetItemByAssetCode(assetCode string) (interface{}, error) {
	collection := config.GetCollection("asset_lists")
	var assetItem models.AssetLists

	err := collection.FindOne(context.Background(), bson.M{"assetCode": assetCode}).Decode(&assetItem)
	if err != nil {
		return nil, err
	}
	if assetItem.Status == 0 {
		return "This Asset item was write off", nil
	} else {
		return &assetItem, nil
	}
}

func CreateAssetItem(assetItem *models.AssetLists) (interface{}, error) {
	collection := config.GetCollection("asset_lists")

	filter := bson.M{"assetName": assetItem.AssetName}

	var existingAsset models.AssetLists
	err := collection.FindOne(context.Background(), filter).Decode(&existingAsset)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No existing asset, proceed to insert
			newAssetCode, err := CreateNewAssetCode()
			if err != nil {
				return nil, err
			}

			assetItem.Status = 1
			assetItem.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
			assetItem.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
			assetItem.AssetCode = newAssetCode

			_, errInv := CreateInvRecord(newAssetCode, "", assetItem.PlaceId)
			if errInv != nil {
				return nil, errInv
			}

			updateFiles := assetItem.UploadAssetListFiles
			assetItem.UploadAssetListFiles = nil

			res, errInsert := collection.InsertOne(context.Background(), assetItem)
			if errInsert != nil {
				return nil, errInsert
			}

			if len(updateFiles) > 0 {
				_, errFile := UploadAssetFile(updateFiles, res.InsertedID.(primitive.ObjectID).Hex())
				if errFile != nil {
					return nil, errFile
				}
			}

			return res, nil
		} else {
			// Some other error
			return nil, err
		}
	}

	// If we get here, asset already exists
	return "This asset name already exists! Please check again!", nil
}

func formatNumber(num int, digits int) string {
	return fmt.Sprintf("%0*d", digits, num+1)
}

func CreateNewAssetCode() (string, error) {
	collection := config.GetCollection("asset_lists")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$addFields", Value: bson.D{
			{Key: "assetCodeInt", Value: bson.D{
				{Key: "$convert", Value: bson.D{
					{Key: "input", Value: "$assetCode"},
					{Key: "to", Value: "int"},
					{Key: "onError", Value: 0},
					{Key: "onNull", Value: 0},
				}},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "maxNumber", Value: bson.D{{Key: "$max", Value: "$assetCodeInt"}}},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return "", err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return "", err
	}

	var maxNumber int
	if len(results) > 0 {
		if val, ok := results[0]["maxNumber"]; ok {
			// BSON to Go int
			switch v := val.(type) {
			case int32:
				maxNumber = int(v)
			case int64:
				maxNumber = int(v)
			case float64:
				maxNumber = int(v)
			default:
				maxNumber = 0
			}
		}
	}

	// 格式化
	return formatNumber(maxNumber, 6), nil
}

func UpdateAssetItem(updateData *models.AssetLists) (interface{}, error) {
	collection := config.GetCollection("asset_lists")

	// Check if the asset item exists
	var existingAsset models.AssetLists
	err := collection.FindOne(context.Background(), bson.M{"_id": updateData.ID}).Decode(&existingAsset)
	if err != nil {
		return nil, err
	}

	if existingAsset.Status == 1 {

		updateData.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

		if updateData.PlaceId != existingAsset.PlaceId {
			_, err := CreateInvRecord(existingAsset.AssetCode, existingAsset.PlaceId, updateData.PlaceId)
			if err != nil {
				return nil, err
			}
		}

		updateFiles := updateData.UploadAssetListFiles
		updateData.UploadAssetListFiles = nil

		res, err2 := collection.UpdateOne(context.Background(), bson.M{"_id": updateData.ID}, bson.M{"$set": updateData})
		if err2 != nil {
			return nil, err2
		}

		if len(updateFiles) > 0 {
			_, errFile := UploadAssetFile(updateFiles, updateData.ID.Hex())
			if errFile != nil {
				return nil, errFile
			}
		}

		return res, nil
	} else {
		return "This asset item may be invalidated or not exist! Please contact admin!", nil
	}
}

func ListAsseetItemsWithFilter(req *dto.ListAssetReqDto) (interface{}, error) {
	collection := config.GetCollection("asset_lists")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filters := bson.M{
		"status": 1,
	}
	if req.AssetCode != "" {
		filters["asset_code"] = req.AssetCode
	}
	if req.AssetName != "" {
		filters["asset_name"] = bson.M{"$regex": req.AssetName, "$options": "i"}
	}
	if len(req.TypeIds) > 0 {
		filters["type_ids"] = bson.M{"$in": req.TypeIds}
	}
	if len(req.PlaceIds) > 0 {
		filters["place_ids"] = bson.M{"$in": req.PlaceIds}
	}
	if len(req.DeptIds) > 0 {
		filters["dept_ids"] = bson.M{"$in": req.DeptIds}
	}
	if len(req.PurchaseDates) > 0 {
		filters["purchase_dates"] = bson.M{"$gte": req.PurchaseDates[0], "$lte": req.PurchaseDates[1]}
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filters}}, // filters must be bson.D or bson.M

		// $lookup locations
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

		// $lookup departments
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

		// $lookup assettypes
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
				{Key: "as", Value: "assettype"},
			},
		}},

		// $addFields assetCodeInt
		{{
			Key: "$addFields", Value: bson.D{
				{Key: "assetCodeInt", Value: bson.D{
					{Key: "$convert", Value: bson.D{
						{Key: "input", Value: "$assetCode"},
						{Key: "to", Value: "int"},
						{Key: "onError", Value: 0}, // 防止 "" 出錯
						{Key: "onNull", Value: 0},
					}},
				}},
			},
		}},

		// unwind
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$location"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$department"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$assettype"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},

		// sort
		{{Key: "$sort", Value: bson.D{{Key: "assetCodeInt", Value: 1}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var FinalResults []bson.M
	if err := cursor.All(ctx, &FinalResults); err != nil {
		return nil, err
	}

	var LastResults []dto.AssetListsPureDetails
	for _, item := range FinalResults {
		if isSponsor, ok := item["sponsor"].(bool); ok && isSponsor {
			item["sponsor"] = "Yes"
		} else {
			item["sponsor"] = "No"
		}

		if isIncludeTax, ok := item["includeTax"].(bool); ok && isIncludeTax {
			item["includeTax"] = "Yes"
		} else {
			item["includeTax"] = "No"
		}

		resultData, err := bson.Marshal(item)
		if err != nil {
			return nil, err
		}

		var assetDetails dto.AssetListsPureDetails
		if err := bson.Unmarshal(resultData, &assetDetails); err != nil {
			return nil, err
		}

		assetDetails.PlaceName = tools.GetNestedString(item, "location", "placeName")
		assetDetails.PlaceCode = tools.GetNestedString(item, "location", "placeCode")

		assetDetails.DeptName = tools.GetNestedString(item, "department", "deptName")
		assetDetails.DeptCode = tools.GetNestedString(item, "department", "deptCode")

		assetDetails.TypeName = tools.GetNestedString(item, "assettype", "typeName")
		assetDetails.TypeCode = tools.GetNestedString(item, "assettype", "typeCode")

		LastResults = append(LastResults, assetDetails)
	}

	return LastResults, nil
}

func ListAllAssetItems() (interface{}, error) {
	collection := config.GetCollection("asset_lists")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		// $lookup locations
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

		// $lookup departments
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

		// $lookup assettypes
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
				{Key: "as", Value: "assettype"},
			},
		}},

		// $addFields assetCodeInt
		{{
			Key: "$addFields", Value: bson.D{
				{Key: "assetCodeInt", Value: bson.D{
					{Key: "$convert", Value: bson.D{
						{Key: "input", Value: "$assetCode"},
						{Key: "to", Value: "int"},
						{Key: "onError", Value: 0}, // 防止 "" 出錯
						{Key: "onNull", Value: 0},
					}},
				}},
			},
		}},

		// unwind
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$location"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$department"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$assettype"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},

		// sort
		{{Key: "$sort", Value: bson.D{{Key: "assetCodeInt", Value: 1}}}},
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

	return results, nil
}

func ListAssetItems(req *dto.ListAssetReqDto) (interface{}, error) {
	collection := config.GetCollection("asset_lists")
	//var assetItems []models.AssetLists

	filters := bson.M{
		"status": 1,
	}
	if req.AssetCode != "" {
		filters["asset_code"] = req.AssetCode
	}
	if req.AssetName != "" {
		filters["asset_name"] = bson.M{"$regex": req.AssetName, "$options": "i"}
	}
	if len(req.TypeIds) > 0 {
		filters["type_ids"] = bson.M{"$in": req.TypeIds}
	}
	if len(req.PlaceIds) > 0 {
		filters["place_ids"] = bson.M{"$in": req.PlaceIds}
	}
	if len(req.DeptIds) > 0 {
		filters["dept_ids"] = bson.M{"$in": req.DeptIds}
	}
	if len(req.PurchaseDates) > 0 {
		filters["purchase_dates"] = bson.M{"$gte": req.PurchaseDates[0], "$lte": req.PurchaseDates[1]}
	}

	skip := int64((req.Page - 1) * req.Limit)
	limit := int64(req.Limit)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filters}}, // filters must be bson.D or bson.M

		// $lookup locations
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

		// $lookup departments
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

		// $lookup assettypes
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
				{Key: "as", Value: "assettype"},
			},
		}},

		// $addFields assetCodeInt
		{{
			Key: "$addFields", Value: bson.D{
				{Key: "assetCodeInt", Value: bson.D{
					{Key: "$convert", Value: bson.D{
						{Key: "input", Value: "$assetCode"},
						{Key: "to", Value: "int"},
						{Key: "onError", Value: 0}, // 防止 "" 出錯
						{Key: "onNull", Value: 0},
					}},
				}},
			},
		}},

		// unwind
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$location"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$department"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$assettype"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},

		// sort
		{{Key: "$sort", Value: bson.D{{Key: "assetCodeInt", Value: 1}}}},

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

	return gin.H{"lists": results, "total": count, "page": req.Page, "limit": req.Limit}, nil
}

func UploadAssetFile(files []models.AssetListFiles, assetId string) (interface{}, error) {

	collectionFiles := config.GetCollection("asset_list_files")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var itemsData []interface{}
	now := time.Now()

	for _, item := range files {
		doc := bson.M{
			"assetId":   assetId,
			"fileName":  item.FileName,
			"fileTyp":   item.FileType,
			"base64":    item.Base64,
			"status":    1,
			"createdAt": now,
			"updatedAt": now,
		}
		itemsData = append(itemsData, doc)
	}

	res, err := collectionFiles.InsertMany(ctx, itemsData)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetListAssetFiles(assetId string) (interface{}, error) {

	collectionFiles := config.GetCollection("asset_list_files")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collectionFiles.Find(ctx, bson.M{"assetId": assetId, "status": 1})
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

func DeleteAssetFile(fileId string) (interface{}, error) {

	objectID, err := primitive.ObjectIDFromHex(fileId)

	collectionFiles := config.GetCollection("asset_list_files")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": objectID}

	count, err := collectionFiles.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	if count > 0 {

		updateRes, err := collectionFiles.UpdateOne(
			ctx,
			filter,
			bson.M{
				"$set": bson.M{
					"status":    0,
					"updatedAt": time.Now(),
				},
			},
		)

		if err != nil {
			return nil, err
		}

		return updateRes, nil
	} else {
		return "File not found", nil
	}

}

func WriteOffInactiveAsset(assetId string) (interface{}, error) {

	collection := config.GetCollection("asset_lists")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(assetId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	res := collection.FindOne(ctx, filter)

	var asset models.AssetLists
	err = res.Decode(&asset)
	if err != nil {
		return nil, err
	}

	if asset.Status == 0 {
		return "This asset have been written off", nil
	} else {
		asset.Status = 0
		asset.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

		_, err := collection.UpdateOne(ctx, filter, asset)
		if err != nil {
			return nil, err
		}
		return asset, nil
	}
}

func AssetDataFinderReturn(assetCode string, assetName string) (models.AssetLists, error) {
	collection := config.GetCollection("asset_lists")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1, "assetCode": assetCode, "assetName": assetName}

	var asset models.AssetLists

	err := collection.FindOne(ctx, filter).Decode(&asset)
	if err != nil {
		return models.AssetLists{}, err
	}

	if asset.Status == 0 {
		return models.AssetLists{}, nil
	}

	return asset, nil
}
