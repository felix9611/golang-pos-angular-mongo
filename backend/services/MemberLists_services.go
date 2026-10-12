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
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateMember(member *models.MemberLists) (interface{}, error) {
	collection := config.GetCollection("member_lists")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1, "name": member.Name, "phone": member.Phone}

	count, err := collection.CountDocuments(ctx, filter)

	if err != nil {
		return nil, err
	}

	if count == 0 {
		ctxNew, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if member.MemberCode == "" {
			member.MemberCode, _ = CreateNewMemberNumber()
		}

		if member.Status == 0 {
			member.Status = 1
		}

		if member.CreatedAt.IsZero() {
			member.CreatedAt = time.Now()
		}

		if member.UpdatedAt.IsZero() {
			member.UpdatedAt = time.Now()
		}

		res, err := collection.InsertOne(ctxNew, member)
		if err != nil {
			return nil, err
		}
		memberId, ok := res.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, mongo.ErrNilDocument
		}

		if len(member.MemberSpecialDays) > 0 {
			var newData models.MemberSpecialDays

			for _, item := range member.MemberSpecialDays {
				newData = models.MemberSpecialDays{
					MemberId:  memberId.Hex(),
					Name:      item.Name,
					Date:      item.Date,
					Remark:    item.Remark,
					Status:    1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				_, err := CreateMemberSpecialDay(newData)

				if err != nil {
					return nil, err
				}
			}
		}

		CreateActionRecord("Member Create", "POST", "Member", member, "Success")

		return res, nil
	} else {
		CreateActionRecord("Member Create", "POST", "Member", member, "Failed")
		return gin.H{"message": "Member already exists"}, nil
	}
}

func CreateMemberSpecialDay(data models.MemberSpecialDays) (interface{}, error) {
	collection := config.GetCollection("member_special_days")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return gin.H{"message": "Failed to create member special day"}, err
	}

	return res, nil
}

func UpdateMemberSpecialDay(data models.MemberSpecialDays) (interface{}, error) {
	collection := config.GetCollection("member_special_days")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.UpdateOne(ctx, bson.M{"_id": data.Id}, bson.M{"$set": data})

	if err != nil {
		return gin.H{"message": "Failed to update member special day"}, err
	}

	return res, nil
}

func VoidMemberSpecialDay(id string) (interface{}, error) {
	collection := config.GetCollection("member_special_days")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var specialDay models.MemberSpecialDays

	err = collection.FindOne(ctx, filter).Decode(&specialDay)
	if err != nil {
		return nil, err
	}

	if specialDay.Status == 1 {

		specialDay.Status = 0
		specialDay.UpdatedAt = time.Now()

		_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": specialDay})

		if err != nil {
			return nil, err
		}

		CreateActionRecord("Member Special Day Remove", "POST", "Member", specialDay, "Success")
		return gin.H{"message": "This special day removed successfully"}, nil

	} else {
		CreateActionRecord("Member Special Day Remove", "POST", "Member", specialDay, "Failed")
		return gin.H{"message": "This special day is already inactive"}, nil
	}

}

func GetOneMemberById(id string) (interface{}, error) {
	collection := config.GetCollection("member_lists")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}

	var member models.MemberLists
	err = collection.FindOne(ctx, filter).Decode(&member)

	if err != nil {
		return nil, err
	}

	if member.Status == 0 {
		return gin.H{"message": "Member is already inactive"}, nil
	} else {
		speacialDayList, err := GetSpecialDayByMemberId(id)

		if len(speacialDayList) == 0 {
			member.MemberSpecialDays = []models.MemberSpecialDays{}
		} else {
			member.MemberSpecialDays = speacialDayList
		}

		if err != nil {
			return gin.H{"message": "Failed to get member special day"}, err
		}

		return &member, nil
	}
}

func GetSpecialDayByMemberId(id string) ([]models.MemberSpecialDays, error) {
	collection := config.GetCollection("member_special_days")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"memberId": id}

	var specialDay []models.MemberSpecialDays
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var special models.MemberSpecialDays
		if err := cursor.Decode(&special); err != nil {
			return nil, err
		}
		specialDay = append(specialDay, special)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return specialDay, nil
}

func VoidOneMember(id string) (interface{}, error) {
	collection := config.GetCollection("member_lists")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var member models.MemberLists
	err = collection.FindOne(ctx, filter).Decode(&member)
	if err != nil {
		return nil, err
	}

	if member.Status == 1 {
		member.Status = 0
		member.UpdatedAt = time.Now()
		_, err := collection.UpdateOne(ctx, filter, bson.M{
			"$set": member,
		})

		if err != nil {
			return nil, err
		}

		CreateActionRecord("Member Inactive", "DELETE", "Member", member, "Success")
		return gin.H{"message": "Member deactivated successfully"}, nil
	} else {
		CreateActionRecord("Member Inactive", "DELETE", "Member", member, "Failed")
		return gin.H{"message": "Member is already inactive"}, nil
	}
}

func UpdateMember(updateDto *models.MemberLists) (interface{}, error) {
	collection := config.GetCollection("member_lists")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": updateDto.Id}

	var existingMember models.MemberLists
	err := collection.FindOne(ctx, filter).Decode(&existingMember)
	if err != nil {
		return nil, err
	}

	if existingMember.Status == 1 {
		updateDto.UpdatedAt = time.Now()

		result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": updateDto})
		if err != nil {
			return nil, err
		}

		if len(updateDto.MemberSpecialDays) > 0 {
			for _, specialDay := range updateDto.MemberSpecialDays {
				if specialDay.Id.IsZero() {
					var newData models.MemberSpecialDays
					newData = models.MemberSpecialDays{
						MemberId:  updateDto.Id.Hex(),
						Name:      specialDay.Name,
						Date:      specialDay.Date,
						Remark:    specialDay.Remark,
						Status:    1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}

					_, err := CreateMemberSpecialDay(newData)

					if err != nil {
						return nil, err
					}
				} else {
					_, err := UpdateMemberSpecialDay(specialDay)

					if err != nil {
						return nil, err
					}
				}
			}
		}

		CreateActionRecord("Member Update", "POST", "Member", updateDto, "Success")

		return result, nil
	} else {
		CreateActionRecord("Member Update", "POST", "Member", updateDto, "Failed")
		return gin.H{"message": "This member is already inactive"}, nil
	}
}

func CreateNewMemberNumber() (string, error) {
	collection := config.GetCollection("member_lists")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$addFields", Value: bson.D{
			{Key: "memberCodeInt", Value: bson.D{
				{Key: "$convert", Value: bson.D{
					{Key: "input", Value: "$memberCode"},
					{Key: "to", Value: "int"},
					{Key: "onError", Value: 0},
					{Key: "onNull", Value: 0},
				}},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "maxNumber", Value: bson.D{{Key: "$max", Value: "$memberCodeInt"}}},
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

func ListMemberLists(pageDto dto.MemberListsDto) (dto.MemberList, error) {
	collection := config.GetCollection("member_lists")
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
				{"name": bson.M{"$regex": pageDto.Name, "$options": "i"}},
				{"phone": bson.M{"$regex": pageDto.Name, "$options": "i"}},
			},
		}
	}

	count, errCount := collection.CountDocuments(ctx, filter)
	if errCount != nil {
		return dto.MemberList{}, errCount
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "member_class"},
				{Key: "let", Value: bson.D{
					{Key: "classIdStr", Value: bson.D{
						{Key: "$convert", Value: bson.D{
							{Key: "input", Value: "$classId"},
							{Key: "to", Value: "objectId"},
							{Key: "onError", Value: nil},
							{Key: "onNull", Value: nil},
						}},
					}},
				}},
				{Key: "pipeline", Value: mongo.Pipeline{
					{{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$eq", Value: bson.A{"$_id", "$$classIdStr"}},
						}},
					}}},
				}},
				{Key: "as", Value: "memberClass"},
			},
		}},
		// unwind
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$memberClass"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return dto.MemberList{}, err
	}
	defer cursor.Close(ctx)

	var results []dto.MemberListDetail
	if err = cursor.All(ctx, &results); err != nil {
		return dto.MemberList{}, err
	}

	return dto.MemberList{
		Lists: results,
		Total: count,
		Page:  pageDto.Page,
		Limit: pageDto.Limit,
	}, nil
}

func FindMembersByName(data dto.MemberListsByNameDto) ([]models.MemberLists, error) {
	collection := config.GetCollection("member_lists")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1}

	if data.Name != "" {
		filter = bson.M{
			"status": 1,
			"$or": []bson.M{
				{"name": bson.M{"$regex": data.Name, "$options": "i"}},
				{"phone": bson.M{"$regex": data.Name, "$options": "i"}},
			},
		}
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.MemberLists
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
