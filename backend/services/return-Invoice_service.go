package services

import (
	"context"
	"golang-pos-mongo-backend/backend/config"
	"golang-pos-mongo-backend/backend/dto"
	"golang-pos-mongo-backend/backend/models"
	"golang-pos-mongo-backend/backend/tools"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateReturnInvoice creates a return invoice along with its items atomically.
func CreateReturnInvoice(req dto.CreateReturnInvoiceDto) (interface{}, error) {

	invoiceResult, err := GetOneInvoiceDetailByInvNo(req.ReturnInvoiceNo)
	if err != nil {
		return gin.H{"error": "failed to fetch location: %w"}, err
	}
	if invoiceResult == nil {
		return gin.H{"error": "invoice not found"}, err
	}

	// verify Location exists
	locationResult, err := GetOneLocationById(req.ReturnLocationId)
	if err != nil {
		return gin.H{"error": "failed to fetch location: %w"}, err
	}
	if locationResult == nil {
		return gin.H{"error": "invoice not found"}, err
	}

	returnCaseNumber := GetReturnCaseNumber()
	if err != nil {
		return gin.H{"error": "failed to generate return case number: %w"}, err
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	returnInvoice := models.ReturnInvoice{
		ID:               primitive.NewObjectID(),
		ReturnCaseNumber: returnCaseNumber,
		InvoiceNumber:    req.ReturnInvoiceNo,
		ReturnReason:     req.ReturnReason,
		ReturnLocationId: req.ReturnLocationId,
		ReturnDetail:     req.ReturnDatail,
		ReturnMethod:     req.ReturnMethod,
		ProcessMethod:    req.ProcessMethod,
		Refund:           req.Refund,
		RefundAmount:     req.RefundAmount,
		RefundMethod:     req.RefundMethod,
		ReturnToVendor:   req.ReturnToVendor,
		Status:           1,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	// ReturnInvoice
	collection := config.GetCollection("return_invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, returnInvoice)
	if err != nil {
		CreateActionRecord("Create Return Invoice", "POST", "Return Invoice", returnInvoice, "Failed")
		return gin.H{"error": "failed to insert return invoice: %w"}, err
	}

	// batch to write return invoice items
	if len(req.ReturnItems) > 0 {
		returnInvoiceHex := returnInvoice.ID.Hex()
		itemsToInsert := make([]interface{}, len(req.ReturnItems))

		for i, item := range req.ReturnItems {
			itemsToInsert[i] = models.ReturnInvoiceItem{
				ID:              primitive.NewObjectID(),
				ReturnDataId:    returnInvoiceHex,
				ReturnInvoiceId: item.ReturnInvoiceId,
				ProductId:       item.ProductId,
				Qty:             item.Qty,
				Price:           item.Price,
				TaxType:         item.TaxType,
				TaxCode:         item.TaxCode,
				TaxRate:         item.TaxRate,
				TaxAmount:       item.TaxAmount,
			}
		}

		itemCollection := config.GetCollection("return_invoice_items")
		_, err = itemCollection.InsertMany(ctx, itemsToInsert)
		if err != nil {
			CreateActionRecord("Create Return Invoice", "POST", "Return Invoice", returnInvoice, "Failed")
			return gin.H{"error": "failed to insert return invoice items: %w"}, err
		}
	}

	CreateActionRecord("Create Return Invoice", "POST", "Return Invoice", returnInvoice, "Success")
	return &returnInvoice, nil
}

func GetReturnCaseNumber() string {
	answer := 1000000000 + rand.Int64N(9000000000)
	return strconv.FormatInt(answer, 10)
}

func GetReturnInvoiceById(id string) (interface{}, error) {
	collection := config.GetCollection("return_invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},

		// $lookup invoices
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "invoices"},
				{Key: "let", Value: bson.D{{Key: "locationIdStr", Value: "$invoiceNumber"}}},
				{Key: "pipeline", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$number", "$$locationIdStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "invoice"},
			}},
		},
		//$unwind
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$invoice"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},

		// $lookup member
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "members"},
			{Key: "let", Value: bson.D{
				{Key: "memberIdStr", Value: bson.D{{Key: "$toObjectId", Value: "$invoice.memberId"}}},
			}},
			{Key: "pipeline", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{
						{Key: "$and", Value: bson.A{
							bson.D{{Key: "$eq", Value: bson.A{"$_id", "$$memberIdStr"}}},
							bson.D{{Key: "$ne", Value: bson.A{"$$memberIdStr", nil}}},
						}},
					}},
				}}},
			}},
			{Key: "as", Value: "member"},
		}}},

		//$unwind
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$member"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var results []dto.ReturnInvoiceDetail
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	returnItemList, err := GetReturnItemList(id)

	if err != nil {
		return nil, err
	}

	return dto.ReturnInvoiceDetailResponse{
		ID:               results[0].ID,
		InvoiceNumber:    results[0].InvoiceNumber,
		ReturnDate:       results[0].ReturnDate,
		ReturnReason:     results[0].ReturnReason,
		ReturnLocationId: results[0].ReturnLocationId,
		ReturnCaseNumber: results[0].ReturnCaseNumber,
		ReturnDatail:     results[0].ReturnDatail,
		ReturnMethod:     results[0].ReturnMethod,
		ProcessMethod:    results[0].ProcessMethod,
		Refund:           results[0].Refund,
		RefundAmount:     results[0].RefundAmount,
		RefundMethod:     results[0].RefundMethod,
		ReturnToVendor:   results[0].ReturnToVendor,
		Invoice:          results[0].Invoice,
		// Member:           results[0].Member,
		ReturnItems: returnItemList,
	}, err
}

func GetReturnItemList(id string) ([]dto.ReturnInvoiceItemDetail, error) {
	collection := config.GetCollection("return_invoice_items")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"returnDataId": id}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},

		// $lookup product
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "product_lists"},
				{Key: "let", Value: bson.D{
					{Key: "productIdStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$productId"},
							{Key: "to", Value: "objectId"},
							{Key: "onError", Value: nil},
							{Key: "onNull", Value: nil},
						}},
					}},
				}},
				{Key: "pipeline", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$_id", "$$productIdStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "product"},
			}},
		},
		//$unwind
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$product"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var results []dto.ReturnInvoiceItemDetail
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func ListReturnInvoice(req dto.InvoiceListRequestDto) (dto.ReturnInvoiceList, error) {
	collection := config.GetCollection("return_invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filters := bson.M{
		"status": 1,
	}

	if len(req.DateRange) > 0 {
		filters["createdAt"] = bson.M{"$gte": req.DateRange[0], "$lte": req.DateRange[1]}
	}

	if req.InvoiceNumber != "" {
		filters["invoiceNumber"] = bson.M{"$regex": req.InvoiceNumber, "$options": "i"}
	}

	if req.CaseNumber != "" {
		filters["returnCaseNumber"] = bson.M{"$regex": req.CaseNumber, "$options": "i"}
	}

	if len(req.ReturnLocationIds) > 0 {
		filters["returnLocationId"] = bson.M{"$in": req.ReturnLocationIds}
	}

	skip := int64((req.Page - 1) * req.Limit)
	limit := int64(req.Limit)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filters}},

		// $lookup invoices
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "invoices"},
				{Key: "let", Value: bson.D{{Key: "invoiceIdStr", Value: "$invoiceNumber"}}},
				{Key: "pipeline", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$number", "$$invoiceIdStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "invoice"},
			},
		}},
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "locations"},
				{Key: "let", Value: bson.D{
					{Key: "locationIdStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$returnLocationId"},
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
			}}},
		//$unwind
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$invoice"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$location"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		// pagination
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return dto.ReturnInvoiceList{}, err
	}

	defer cursor.Close(ctx)

	var results []dto.ReturnInvoiceDetail
	if err := cursor.All(ctx, &results); err != nil {
		return dto.ReturnInvoiceList{}, err
	}

	count, errCount := collection.CountDocuments(ctx, filters)
	if errCount != nil {
		return dto.ReturnInvoiceList{}, errCount
	}

	return dto.ReturnInvoiceList{
		Lists: results,
		Total: count,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

func ListReturnInvoiceWithFilter(req dto.InvoiceListRequestDto) ([]dto.ReturnInvoicePureList, error) {
	collection := config.GetCollection("return_invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filters := bson.M{
		"status": 1,
	}

	if len(req.DateRange) > 0 {
		filters["createdAt"] = bson.M{"$gte": req.DateRange[0], "$lte": req.DateRange[1]}
	}

	if req.InvoiceNumber != "" {
		filters["invoiceNumber"] = bson.M{"$regex": req.InvoiceNumber, "$options": "i"}
	}

	if req.CaseNumber != "" {
		filters["returnCaseNumber"] = bson.M{"$regex": req.CaseNumber, "$options": "i"}
	}

	if len(req.ReturnLocationIds) > 0 {
		filters["returnLocationId"] = bson.M{"$in": req.ReturnLocationIds}
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filters}},

		// $lookup invoices
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "invoices"},
				{Key: "let", Value: bson.D{{Key: "invoiceIdStr", Value: "$invoiceNumber"}}},
				{Key: "pipeline", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$number", "$$invoiceIdStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "invoice"},
			},
		}},
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "locations"},
				{Key: "let", Value: bson.D{
					{Key: "locationIdStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$returnLocationId"},
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
			}}},
		//$unwind
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$invoice"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$location"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		// pagination
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

	var LastResults []dto.ReturnInvoicePureList

	for _, item := range FinalResults {

		var result dto.ReturnInvoicePureList

		switch item["refund"] {
		case true:
			result.Refund = "Yes"
		case false:
			result.Refund = "No"
		default:
			result.Refund = ""
		}

		switch item["returnToVendor"] {
		case true:
			result.ReturnToVendor = "Yes"
		case false:
			result.ReturnToVendor = "No"
		default:
			result.ReturnToVendor = ""
		}

		result.InvoiceNumber = item["invoiceNumber"].(string)
		result.ReturnCaseNumber = item["returnCaseNumber"].(string)
		result.ReturnReason = item["returnReason"].(string)
		result.ReturnDate = item["returnDate"].(string)
		result.ReturnDetail = item["returnDatail"].(string)
		result.ReturnMethod = item["returnMethod"].(string)
		result.ProcessMethod = item["processMethod"].(string)
		result.RefundAmount = item["refundAmount"].(float64)
		result.RefundMethod = item["refundMethod"].(string)
		result.CreatedAt = item["createdAt"].(string)
		result.UpdatedAt = item["updatedAt"].(string)
		result.TaxRefNo = tools.GetNestedValue(item, "invoice", "taxRefNo").(string)
		result.PaidAmount = tools.GetFloat64(item, "invoice", "totalAmount").(float64)
		result.PaidAmountWithTax = tools.GetFloat64(item, "invoice", "taxTotal").(float64)
		result.InvoiceRemark = tools.GetNestedValue(item, "invoice", "remark").(string)
		result.InvoiceCreatedAt = tools.GetNestedValue(item, "invoice", "createdAt").(primitive.DateTime).Time().Format("2006-01-02 15:04:05")
		result.PlaceName = tools.GetNestedValue(item, "location", "placeCode").(string)
		result.PlaceCode = tools.GetNestedValue(item, "location", "placeName").(string)

		LastResults = append(LastResults, result)
	}

	return LastResults, nil
}
