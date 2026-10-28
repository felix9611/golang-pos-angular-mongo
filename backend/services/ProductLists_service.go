package services

import (
	"context"
	"fmt"
	"golang-pos-mongo-backend/backend/config"
	"golang-pos-mongo-backend/backend/dto"
	"golang-pos-mongo-backend/backend/models"
	"golang-pos-mongo-backend/backend/tools"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetProductListsCollection() *mongo.Collection {
	return config.GetCollection("product_lists")
}

func GetProductFileListsCollection() *mongo.Collection {
	return config.GetCollection("product_file_lists")
}

func CreateProductItem(productItem *models.ProductLists) (interface{}, error) {
	collection := GetProductListsCollection()

	filter := bson.M{"productName": productItem.ProductName, "status": 1}

	var existingAsset models.ProductLists
	err := collection.FindOne(context.Background(), filter).Decode(&existingAsset)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No existing asset, proceed to insert
			newProductCode, err := CreateNewProductCode()
			if err != nil {
				return nil, err
			}

			productItem.Status = 1
			productItem.CreatedAt = time.Now()
			productItem.UpdatedAt = time.Now()
			productItem.ProductCode = newProductCode

			updateFiles := productItem.UploadProductListFiles
			productItem.UploadProductListFiles = nil

			resNew, errInsert := collection.InsertOne(context.Background(), productItem)
			if errInsert != nil {
				return nil, errInsert
			}

			CreateActionRecord("Product Item Create", "POST", "Product Item", productItem, "Success")

			if len(updateFiles) > 0 {
				_, errFile := UploadProductFile(updateFiles, resNew.InsertedID.(primitive.ObjectID).Hex())
				if errFile != nil {
					return nil, errFile
				}
			}
			return resNew, nil
		} else {
			// Some other error
			return nil, err
		}
	}
	return gin.H{"error": "Product name already exists!"}, nil
}

func UploadProductFile(files []models.ProductListFiles, id string) (interface{}, error) {
	collection := GetProductFileListsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var itemsData []interface{}
	now := time.Now()

	for _, item := range files {
		doc := bson.M{
			"productId": id,
			"fileName":  item.FileName,
			"fileType":  item.FileType,
			"base64":    item,
			"status":    1,
			"createdAt": now,
			"updatedAt": now,
		}
		itemsData = append(itemsData, doc)
	}

	res, err := collection.InsertMany(ctx, itemsData)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func formatNumberNew(num int, digits int) string {
	return fmt.Sprintf("%0*d", digits, num+1)
}

func CreateNewProductCode() (string, error) {
	collection := GetProductListsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$addFields", Value: bson.D{
			{Key: "productCodeInt", Value: bson.D{
				{Key: "$convert", Value: bson.D{
					{Key: "input", Value: "$productCode"},
					{Key: "to", Value: "int"},
					{Key: "onError", Value: 0},
					{Key: "onNull", Value: 0},
				}},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "maxNumber", Value: bson.D{{Key: "$max", Value: "$productCodeInt"}}},
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

	// format
	return formatNumberNew(maxNumber, 6), nil
}

func GetOneProductItemByID(id string) (interface{}, error) {
	collection := GetProductListsCollection()
	var productItem models.ProductLists

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err2 := collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&productItem)
	if err2 != nil {
		return nil, err2
	}
	if productItem.Status == 0 {
		return "This Product item was write off", nil
	} else {

		return &productItem, nil
	}
}

func GetOneProductItemByCode(code string) (interface{}, error) {
	collection := GetProductListsCollection()

	var productItem models.ProductLists

	err2 := collection.FindOne(context.Background(), bson.M{"productCode": code}).Decode(&productItem)
	if err2 != nil {
		return nil, err2
	}

	if productItem.Status == 0 {
		return "This Product item was write off", nil
	} else {
		return &productItem, nil
	}
}

func UpdateProductItem(updateData *models.ProductLists) (interface{}, error) {
	collection := GetProductListsCollection()

	// Check if the product item exists
	var existingProduct models.ProductLists
	err := collection.FindOne(context.Background(), bson.M{"_id": updateData.ID}).Decode(&existingProduct)
	if err != nil {
		return nil, err
	}

	if existingProduct.Status == 1 {
		// Update the product item
		updateData.UpdatedAt = time.Now()

		updateFiles := updateData.UploadProductListFiles
		updateData.UploadProductListFiles = nil

		res, err2 := collection.UpdateOne(context.Background(), bson.M{"_id": updateData.ID}, bson.M{"$set": updateData})
		if err2 != nil {
			return nil, err2
		}

		CreateActionRecord("Product Item Update", "POST", "Product Item", updateData, "Success")

		if len(updateFiles) > 0 {
			_, errFile := UploadProductFile(updateFiles, updateData.ID.Hex())
			if errFile != nil {
				return nil, errFile
			}
		}

		return res, nil
	} else {
		return gin.H{"error": "This Product item was voided"}, nil
	}
}

func GetListProductFiles(id string) ([]models.ProductListFiles, error) {
	collection := GetProductFileListsCollection()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"productId": id, "status": 1})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.ProductListFiles
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func DeleteProductFile(fileId string) (interface{}, error) {

	objectID, err := primitive.ObjectIDFromHex(fileId)

	collectionFiles := GetProductFileListsCollection()
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
		return gin.H{"error": "File not found"}, nil
	}
}

func ListProductItemsWithFilters(req *dto.ListProductReqDto) ([]dto.ListProductDetailPure, error) {
	collection := GetProductListsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filters := bson.M{
		"status": 1,
	}

	if req.ProductCode != "" {
		filters["productCode"] = req.ProductCode
	}

	if req.ProductName != "" {
		filters["productName"] = bson.M{"$regex": req.ProductName, "$options": "i"}
	}

	if len(req.TypeIds) > 0 {
		filters["typeIds"] = bson.M{"$in": req.TypeIds}
	}

	if len(req.DeptIds) > 0 {
		filters["deptIds"] = bson.M{"$in": req.DeptIds}
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filters}}, // filters must be bson.D or bson.M
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
				{Key: "from", Value: "product_types"},
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
				{Key: "as", Value: "producttype"},
			},
		}},

		// $addFields assetCodeInt
		{{
			Key: "$addFields", Value: bson.D{
				{Key: "productCodeInt", Value: bson.D{
					{Key: "$convert", Value: bson.D{
						{Key: "input", Value: "$productCode"},
						{Key: "to", Value: "int"},
						{Key: "onError", Value: 0},
						{Key: "onNull", Value: 0},
					}},
				}},
			},
		}},

		// unwind
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$department"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$producttype"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},

		// sort
		{{Key: "$sort", Value: bson.D{{Key: "productCodeInt", Value: 1}}}},
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

	var LastResults []dto.ListProductDetailPure
	for _, item := range FinalResults {
		resultData, err := bson.Marshal(item)
		if err != nil {
			return nil, err
		}

		var productDetails dto.ListProductDetailPure
		if err := bson.Unmarshal(resultData, &productDetails); err != nil {
			return nil, err
		}

		productDetails.TypeName = tools.GetNestedValue(item, "producttype", "typeName").(string)
		productDetails.TypeCode = tools.GetNestedValue(item, "producttype", "typeCode").(string)
		productDetails.DeptName = tools.GetNestedValue(item, "department", "deptName").(string)
		productDetails.DeptCode = tools.GetNestedValue(item, "department", "deptCode").(string)

		LastResults = append(LastResults, productDetails)
	}

	return LastResults, err
}

func ListProductItems(req *dto.ListProductReqDto) (interface{}, error) {
	collection := GetProductListsCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filters := bson.M{
		"status": 1,
	}

	if req.ProductCode != "" {
		filters["productCode"] = req.ProductCode
	}

	if req.ProductName != "" {
		filters["productName"] = bson.M{"$regex": req.ProductName, "$options": "i"}
	}

	if len(req.TypeIds) > 0 {
		filters["typeIds"] = bson.M{"$in": req.TypeIds}
	}

	if len(req.DeptIds) > 0 {
		filters["deptIds"] = bson.M{"$in": req.DeptIds}
	}

	skip := int64((req.Page - 1) * req.Limit)
	limit := int64(req.Limit)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filters}}, // filters must be bson.D or bson.M
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
				{Key: "from", Value: "product_types"},
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
				{Key: "as", Value: "producttype"},
			},
		}},

		// $addFields assetCodeInt
		{{
			Key: "$addFields", Value: bson.D{
				{Key: "productCodeInt", Value: bson.D{
					{Key: "$convert", Value: bson.D{
						{Key: "input", Value: "$productCode"},
						{Key: "to", Value: "int"},
						{Key: "onError", Value: 0},
						{Key: "onNull", Value: 0},
					}},
				}},
			},
		}},

		// unwind
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$department"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$producttype"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},

		// sort
		{{Key: "$sort", Value: bson.D{{Key: "productCodeInt", Value: 1}}}},

		// pagination
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
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

	return gin.H{"lists": results, "total": count, "page": req.Page, "limit": req.Limit}, nil
}

func InactiveProductItemsByID(id string) (interface{}, error) {
	collection := GetProductListsCollection()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingProduct models.ProductLists
	err = collection.FindOne(ctx, filter).Decode(&existingProduct)
	if err != nil {
		return nil, err
	}

	if existingProduct.Status == 1 {
		existingProduct.Status = 0
		existingProduct.UpdatedAt = time.Now()
		_, err := collection.UpdateOne(ctx, filter, bson.M{
			"$set": existingProduct,
		})

		if err != nil {
			return gin.H{"error": "Something went wrong, please try again"}, err
		}

		CreateActionRecord("Product Item Inactive", "DELETE", "Product Item", existingProduct, "Success")

		return gin.H{"message": "This product deleted successfully"}, nil
	} else {
		CreateActionRecord("Product Item Inactive", "DELETE", "Product Item", existingProduct, "Failed")
		return gin.H{"error": "This product is already inactive"}, nil
	}
}
