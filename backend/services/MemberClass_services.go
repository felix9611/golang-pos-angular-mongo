package services

import (
	"context"
	"golang-pos-mongo-backend/backend/config"
	"golang-pos-mongo-backend/backend/dto"
	"golang-pos-mongo-backend/backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMemberClass(memberClass *models.MemberClass) (interface{}, error) {
	collection := config.GetCollection("member_class")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1, "className": memberClass.ClassName, "classCode": memberClass.ClassCode}

	count, err := collection.CountDocuments(ctx, filter)

	if err != nil {
		return gin.H{"error": "Create Member Class Failed"}, err
	}

	if count == 0 {
		ctxNew, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if memberClass.Status == 0 {
			memberClass.Status = 1
		}

		if memberClass.CreatedAt.IsZero() {
			memberClass.CreatedAt = time.Now()
		}

		if memberClass.UpdatedAt.IsZero() {
			memberClass.UpdatedAt = time.Now()
		}

		_, err := collection.InsertOne(ctxNew, memberClass)

		if err != nil {
			return gin.H{"error": "Create Member Class Failed"}, err
		}

		CreateActionRecord("MemberClass Create", "POST", "MemberClass", memberClass, "Success")
		return memberClass, nil

	} else {
		CreateActionRecord("MemberClass Create", "POST", "MemberClass", memberClass, "Failed")
		return gin.H{"error": "Member Class already exists"}, nil
	}
}

func GetOneMemberClassById(id string) (interface{}, error) {
	collection := config.GetCollection("member_class")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}

	var memberClass models.MemberClass
	err = collection.FindOne(ctx, filter).Decode(&memberClass)

	if err != nil {
		return nil, err
	}

	if memberClass.Status == 0 {
		return "Member Class is inactive", nil
	} else {
		return &memberClass, nil
	}
}

func VoidOneMemberClass(id string) (interface{}, error) {
	collection := config.GetCollection("member_class")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var memberClass models.MemberClass
	err = collection.FindOne(ctx, filter).Decode(&memberClass)
	if err != nil {
		return nil, err
	}

	if memberClass.Status == 1 {
		memberClass.Status = 0
		memberClass.UpdatedAt = time.Now()
		_, err := collection.UpdateOne(ctx, filter, bson.M{
			"$set": memberClass,
		})

		if err != nil {
			return nil, err
		}

		CreateActionRecord("Member Inactive", "DELETE", "Member Class", memberClass, "Success")
		return gin.H{"message": "Member deactivated successfully"}, nil
	} else {
		CreateActionRecord("Member Inactive", "DELETE", "Member Class", memberClass, "Failed")
		return gin.H{"error": "Member is already inactive"}, nil
	}
}

func ListMemberClass(pageDto *dto.MemberClassListDto) (interface{}, error) {
	collection := config.GetCollection("member_class")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if pageDto.Page < 1 {
		pageDto.Page = 1
	}
	if pageDto.Limit < 1 {
		pageDto.Limit = 10
	}

	skip := (pageDto.Page - 1) * pageDto.Limit
	limit := pageDto.Limit

	filter := bson.M{"status": 1}

	if pageDto.Name != "" {
		filter = bson.M{
			"status": 1,
			"$or": []bson.M{
				{"className": bson.M{"$regex": pageDto.Name, "$options": "i"}},
				{"classCode": bson.M{"$regex": pageDto.Name, "$options": "i"}},
			},
		}
	}

	count, errCount := collection.CountDocuments(ctx, filter)
	if errCount != nil {
		return nil, errCount
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	var memberClasses []models.MemberClass
	for cursor.Next(ctx) {
		var memberClass models.MemberClass
		if err := cursor.Decode(&memberClass); err != nil {
			return nil, err
		}
		memberClasses = append(memberClasses, memberClass)
	}

	return dto.MemberClassList{
		Lists: memberClasses,
		Total: count,
		Page:  pageDto.Page,
		Limit: pageDto.Limit,
	}, nil
}

func ListAllMemberClass() (interface{}, error) {
	collection := config.GetCollection("member_class")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var memberClasses []models.MemberClass
	for cursor.Next(ctx) {
		var memberClass models.MemberClass
		if err := cursor.Decode(&memberClass); err != nil {
			return nil, err
		}
		memberClasses = append(memberClasses, memberClass)
	}

	return memberClasses, nil
}

func UpdateMemberClass(updateDto *models.MemberClass) (interface{}, error) {
	collection := config.GetCollection("member_class")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": updateDto.ID}

	var existingMemberClass models.MemberClass
	err := collection.FindOne(ctx, filter).Decode(&existingMemberClass)
	if err != nil {
		return nil, err
	}

	if existingMemberClass.Status == 1 {
		updateDto.UpdatedAt = time.Now()

		result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": updateDto})
		if err != nil {
			return nil, err
		}

		CreateActionRecord("MemberClass Update", "POST", "MemberClass", updateDto, "Success")

		return result, nil
	} else {
		CreateActionRecord("MemberClass Update", "POST", "MemberClass", updateDto, "Failed")
		return gin.H{"message": "Member Class is already inactive"}, nil
	}

}
