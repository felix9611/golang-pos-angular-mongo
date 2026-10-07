package services

import (
	"context"
	"golang-fixedasset-mongo-backend/backend/config"
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func BatchCreateTaxInformation(TaxInfos []models.TaxInformations) (interface{}, error) {
	for _, taxInfo := range TaxInfos {
		result, _ := CreateTaxInformation(&taxInfo)

		if result == nil {
			return "failed to create vendor", nil
		}
	}

	return "batch insert completed", nil
}

func CreateTaxInformation(taxInfo *models.TaxInformations) (interface{}, error) {
	filter := bson.M{
		"status":      1,
		"countryCode": taxInfo.CountryCode,
		"countryName": taxInfo.CountryName,
		"taxType":     taxInfo.TaxType,
		"taxCode":     taxInfo.TaxCode,
		"taxName":     taxInfo.TaxName,
	}

	collection := config.GetCollection("tax_informations")

	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		if taxInfo.Status == 0 {
			taxInfo.Status = 1
		}

		if taxInfo.CreatedAt.IsZero() {
			taxInfo.CreatedAt = time.Now()
		}

		if taxInfo.UpdatedAt.IsZero() {
			taxInfo.UpdatedAt = time.Now()
		}

		result, err := collection.InsertOne(context.Background(), taxInfo)
		if err != nil {
			return nil, err
		}

		return result, nil
	} else {
		return "Tax Information with the same name already exists", nil
	}
}

func GetOneTaxInformation(id string) (interface{}, error) {
	collection := config.GetCollection("tax_informations")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID, "status": 1}

	var taxInfo models.TaxInformations

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, filter).Decode(&taxInfo)
	if err != nil {
		return nil, err
	}

	if taxInfo.Status == 0 {
		return "This Tax Information is inactive", nil
	} else {
		return &taxInfo, nil
	}
}

func VoidOneTaxInformation(id string) (interface{}, error) {
	collection := config.GetCollection("tax_informations")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var taxInfo models.TaxInformations

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, filter).Decode(&taxInfo)

	if err != nil {
		return nil, err
	}

	if taxInfo.Status == 1 {

		result, err := collection.UpdateOne(ctx, filter, bson.M{
			"$set": bson.M{
				"status":     0,
				"updated_at": time.Now(),
			},
		})

		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		return "This Tax Information is already voided", nil
	}
}

func UpdateTaxInformation(updateData *models.TaxInformations) (interface{}, error) {
	collection := config.GetCollection("tax_informations")

	filter := bson.M{"_id": updateData.ID, "status": 1}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var checkTaxInfo models.TaxInformations
	err := collection.FindOne(ctx, filter).Decode(&checkTaxInfo)
	if err != nil {
		return nil, err
	}

	if checkTaxInfo.Status == 1 {
		updateData.UpdatedAt = time.Now()

		result, err := collection.UpdateOne(ctx, filter, bson.M{
			"$set": updateData,
		})
		if err != nil {
			return nil, err
		}

		return result, nil
	} else {
		return "This Tax Information is already voided", nil
	}
}

func TaxInformationListWithoutPagination(pageDto *dto.TaxInformationListDTO) (interface{}, error) {
	collection := config.GetCollection("tax_informations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1}

	if pageDto.NameCode != "" {
		filter = bson.M{
			"status": 1,
			"$or": []bson.M{
				{"nationCode": bson.M{"$regex": pageDto.NameCode, "$options": "i"}},
				{"nationName": bson.M{"$regex": pageDto.NameCode, "$options": "i"}},
				{"countryCode": bson.M{"$regex": pageDto.NameCode, "$options": "i"}},
				{"countryName": bson.M{"$regex": pageDto.NameCode, "$options": "i"}},
			},
		}
	}

	if pageDto.Tax != "" {
		filter = bson.M{
			"status": 1,
			"$or": []bson.M{
				{"taxType": bson.M{"$regex": pageDto.Tax, "$options": "i"}},
				{"taxCode": bson.M{"$regex": pageDto.Tax, "$options": "i"}},
				{"taxName": bson.M{"$regex": pageDto.Tax, "$options": "i"}},
			},
		}
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var TaxInfoList []models.TaxInformations
	for cursor.Next(ctx) {
		var taxInfo models.TaxInformations
		if err := cursor.Decode(&taxInfo); err != nil {
			return nil, err
		}
		TaxInfoList = append(TaxInfoList, taxInfo)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return TaxInfoList, nil
}

func TaxInformationList(pageDto *dto.TaxInformationListDTO) (interface{}, error) {
	if pageDto.Page < 1 {
		pageDto.Page = 1
	}

	if pageDto.Limit < 1 {
		pageDto.Limit = 10
	}

	collection := config.GetCollection("tax_informations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1}

	if pageDto.NameCode != "" {
		filter = bson.M{
			"status": 1,
			"$or": []bson.M{
				{"nationCode": bson.M{"$regex": pageDto.NameCode, "$options": "i"}},
				{"nationName": bson.M{"$regex": pageDto.NameCode, "$options": "i"}},
				{"countryCode": bson.M{"$regex": pageDto.NameCode, "$options": "i"}},
				{"countryName": bson.M{"$regex": pageDto.NameCode, "$options": "i"}},
			},
		}
	}

	if pageDto.Tax != "" {
		filter = bson.M{
			"status": 1,
			"$or": []bson.M{
				{"taxType": bson.M{"$regex": pageDto.Tax, "$options": "i"}},
				{"taxCode": bson.M{"$regex": pageDto.Tax, "$options": "i"}},
				{"taxName": bson.M{"$regex": pageDto.Tax, "$options": "i"}},
			},
		}
	}

	count, errCount := collection.CountDocuments(ctx, filter)

	if errCount != nil {
		return nil, errCount
	}

	skip := int64((pageDto.Page - 1) * pageDto.Limit)
	limit := int64(pageDto.Limit)

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var TaxInfoList []models.TaxInformations
	for cursor.Next(ctx) {
		var taxInfo models.TaxInformations
		if err := cursor.Decode(&taxInfo); err != nil {
			return nil, err
		}
		TaxInfoList = append(TaxInfoList, taxInfo)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return gin.H{"lists": TaxInfoList, "total": count, "page": pageDto.Page, "limit": pageDto.Limit}, nil
}

func ListAllTaxInformation() (interface{}, error) {
	collection := config.GetCollection("tax_informations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var TaxInfoList []models.TaxInformations
	for cursor.Next(ctx) {
		var taxInfo models.TaxInformations
		if err := cursor.Decode(&taxInfo); err != nil {
			return nil, err
		}
		TaxInfoList = append(TaxInfoList, taxInfo)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return TaxInfoList, nil
}
