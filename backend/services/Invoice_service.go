package services

import (
	"context"
	"fmt"
	"golang-fixedasset-mongo-backend/backend/config"
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateInvoice(data dto.InvoiceDto) (interface{}, error) {
	collection := config.GetCollection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	invoiceNumber, err := CreateNewInvoiceNumber()

	if err != nil {
		return gin.H{"error": "Invoice not created!"}, err
	}
	newNumber := "INV-" + invoiceNumber
	taxRefNo := randomNumber()

	invoice := models.Invoice{
		Number:       newNumber,
		MemberId:     data.MemberId,
		LocationId:   data.LocationId,
		TotalAmount:  data.TotalAmount,
		Discount:     data.Discount,
		DiscountType: data.DiscountType,
		TaxTotal:     data.TaxTotal,
		TaxRefNo:     taxRefNo,
		Remark:       data.Remark,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Status:       1,
	}

	result, err := collection.InsertOne(ctx, invoice)
	if err != nil {
		return gin.H{"error": "Invoice not created!"}, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		invoice.ID = oid
	} else {
		return gin.H{"error": "Cannot get Invoice ID"}, fmt.Errorf("failed to cast InsertedID")
	}

	if len(data.InvoiceItems) > 0 {
		for _, item := range data.InvoiceItems {

			newInvoiceId := invoice.ID.Hex()

			newItem := models.InvoiceItemList{
				InvoiceId:    newInvoiceId,
				ProductId:    item.ProductId,
				Qty:          item.Qty,
				Price:        item.Price,
				Discount:     item.Discount,
				DiscountType: item.DiscountType,
				TaxType:      item.TaxType,
				TaxCode:      item.TaxCode,
				TaxRate:      item.TaxRate,
				TaxAmount:    item.TaxAmount,
			}

			_, err = config.GetCollection("invoice_items").InsertOne(ctx, &newItem)
			if err != nil {
				return gin.H{"error": "Invoice not created!"}, err
			}

			stockOutBody := dto.StockInOutProductLocationDto{
				ProductCode: "",
				ProductId:   item.ProductId,
				LocationId:  data.LocationId,
				Qty:         item.Qty,
				TotalCost:   item.Price,
				TotalPrice:  item.Price,
				PlaceCode:   "",
			}

			_, err = StockOut(stockOutBody)
			if err != nil {
				return gin.H{"error": "Invoice not created!"}, err
			}
		}
	}

	if len(data.InvoicePayment) > 0 {
		for _, item := range data.InvoicePayment {

			newInvoiceId := invoice.ID.Hex()

			newItem := models.InvoicePayment{
				InvoiceId:   newInvoiceId,
				Method:      item.Method,
				Amount:      item.Amount,
				PaymentTime: item.PaymentTime,
			}

			_, err = config.GetCollection("invoice_payments").InsertOne(ctx, &newItem)
			if err != nil {
				return gin.H{"error": "Invoice not created!"}, err
			}
		}
	}

	CreateActionRecord("Create Invoice", "POST", "Member", data, "Success")
	newInvoiceId := invoice.ID.Hex()
	return gin.H{"message": "Invoice created successfully!", "id": newInvoiceId}, nil
}

func randomNumber() string {
	answer := 1000000000 + rand.Int64N(9000000000)
	return strconv.FormatInt(answer, 10)
}

func CreateNewInvoiceNumber() (string, error) {
	collection := config.GetCollection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.M{
				"number": bson.M{"$regex": `^[iI][nN][vV]-\d+`},
			}},
		},
		{
			{Key: "$addFields", Value: bson.M{
				"numberInt": bson.M{
					"$cond": bson.M{
						"if": bson.M{
							"$regexMatch": bson.M{
								"input": "$number",
								"regex": `^[iI][nN][vV]-\d+$`,
							},
						},
						"then": bson.M{
							"$toInt": bson.M{
								"$substrBytes": bson.A{"$number", 4, -1},
							},
						},
						"else": 0,
					},
				},
			}},
		},
		{
			{Key: "$group", Value: bson.M{
				"_id":       nil,
				"maxNumber": bson.M{"$max": "$numberInt"},
			}},
		},
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

func GetOneInvoiceDetailByInvNo(invoiceNumber string) (interface{}, error) {
	collection := config.GetCollection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"number": invoiceNumber}

	var invoice models.Invoice

	err := collection.FindOne(ctx, filter).Decode(&invoice)
	if err != nil {
		return nil, err
	}

	InvoiceDataId := invoice.ID.Hex()

	var invoiceItems []dto.InvoiceItemListBody
	invoiceItemCollection := config.GetCollection("invoice_items")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"invoiceId": InvoiceDataId}}},
		// lookup
		{{
			Key: "$lookup", Value: bson.D{
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
				{Key: "pipeline", Value: mongo.Pipeline{
					{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$_id", "$$productIdStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "product"},
			},
		}},

		// unwind
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$product"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
	}

	cursor, err := invoiceItemCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &invoiceItems); err != nil {
		return nil, err
	}

	var invoicePayments []models.InvoicePayment
	invoicePaymentCollection := config.GetCollection("invoice_payments")

	pipeline = mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"invoiceId": InvoiceDataId}}},
	}

	cursor, err = invoicePaymentCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &invoicePayments); err != nil {
		return nil, err
	}

	finalItem := dto.InvoiceBody{
		ID:              invoice.ID,
		Number:          invoice.Number,
		MemberId:        invoice.MemberId,
		LocationId:      invoice.LocationId,
		TotalAmount:     invoice.TotalAmount,
		Discount:        invoice.Discount,
		DiscountType:    invoice.DiscountType,
		TaxTotal:        invoice.TaxTotal,
		TaxRefNo:        invoice.TaxRefNo,
		Remark:          invoice.Remark,
		CreatedAt:       invoice.CreatedAt,
		UpdatedAt:       invoice.UpdatedAt,
		Status:          invoice.Status,
		InvoiceItems:    invoiceItems,
		InvoicePayments: invoicePayments,
	}

	return finalItem, nil
}

func GetOneInvoiceDetailById(id string) (interface{}, error) {
	collection := config.GetCollection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": objectID}
	if err != nil {
		return nil, err
	}
	var invoice models.Invoice

	err = collection.FindOne(ctx, filter).Decode(&invoice)
	if err != nil {
		return nil, err
	}

	var invoiceItems []dto.InvoiceItemListBody
	invoiceItemCollection := config.GetCollection("invoice_items")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"invoiceId": objectID.Hex()}}},
		// lookup
		{{
			Key: "$lookup", Value: bson.D{
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
				{Key: "pipeline", Value: mongo.Pipeline{
					{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$_id", "$$productIdStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "product"},
			},
		}},

		// unwind
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$product"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
	}

	cursor, err := invoiceItemCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &invoiceItems); err != nil {
		return nil, err
	}

	var invoicePayments []models.InvoicePayment
	invoicePaymentCollection := config.GetCollection("invoice_payments")

	pipeline = mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"invoiceId": objectID.Hex()}}},
	}

	cursor, err = invoicePaymentCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &invoicePayments); err != nil {
		return nil, err
	}

	finalItem := dto.InvoiceBody{
		ID:              objectID,
		Number:          invoice.Number,
		MemberId:        invoice.MemberId,
		LocationId:      invoice.LocationId,
		TotalAmount:     invoice.TotalAmount,
		Discount:        invoice.Discount,
		DiscountType:    invoice.DiscountType,
		TaxTotal:        invoice.TaxTotal,
		TaxRefNo:        invoice.TaxRefNo,
		Remark:          invoice.Remark,
		CreatedAt:       invoice.CreatedAt,
		UpdatedAt:       invoice.UpdatedAt,
		Status:          invoice.Status,
		InvoiceItems:    invoiceItems,
		InvoicePayments: invoicePayments,
	}

	return &finalItem, nil
}

func ListInvoices(req dto.ListInvoiceRequestDto) (dto.ListInvoices, error) {
	collection := config.GetCollection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filters := bson.M{
		"status": 1,
	}

	if len(req.LocationIds) > 0 {
		filters["locationId"] = bson.M{"$in": req.LocationIds}
	}

	if len(req.Date) > 0 {
		filters["createdAt"] = bson.M{"$gte": req.Date[0], "$lte": req.Date[1]}
	}

	skip := int64((req.Page - 1) * req.Limit)
	limit := int64(req.Limit)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filters}},

		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "locations"},
				{Key: "let", Value: bson.D{
					{Key: "locationIdStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$locationId"},
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
			},
		}},

		// unwind
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$location"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},

		{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return dto.ListInvoices{}, err
	}
	defer cursor.Close(ctx)

	var invoices []dto.InvoiceListBody
	if err := cursor.All(ctx, &invoices); err != nil {
		return dto.ListInvoices{}, err
	}

	count, errCount := collection.CountDocuments(ctx, filters)
	if errCount != nil {
		return dto.ListInvoices{}, errCount
	}

	return dto.ListInvoices{Lists: invoices, Total: count, Page: req.Page, Limit: req.Limit}, nil
}

func VoidOneInvoice(invoiceId string) (interface{}, error) {
	collection := config.GetCollection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"_id": invoiceId}, bson.M{"$set": bson.M{"status": 0}})
	if err != nil {
		return nil, err
	}

	CreateActionRecord("Void Invoice", "POST", "Invoice", bson.M{"invoiceId": invoiceId}, "Success")

	return gin.H{"status": "Void Invoice success"}, nil
}
