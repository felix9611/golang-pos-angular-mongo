package services

import (
	"context"
	"golang-pos-mongo-backend/backend/config"
	"golang-pos-mongo-backend/backend/dto"
	"golang-pos-mongo-backend/backend/models"
	"golang-pos-mongo-backend/backend/tools"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetSysUsersCollection() (*mongo.Collection, context.Context) {
	collection := config.GetCollection("sys_users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return collection, ctx
}

func CreateSysUser(user *models.SysUsers) (interface{}, error) {
	collection, ctx := GetSysUsersCollection()

	filter := bson.M{"status": 1, "username": user.Username}

	count, err := collection.CountDocuments(ctx, filter)

	if err != nil {
		return gin.H{"error": "Failed to create user"}, err
	}

	if count == 0 {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if user.Status == 0 {
			user.Status = 1
		}

		if user.CreatedAt.IsZero() {
			user.CreatedAt = time.Now()
		}

		if user.UpdatedAt.IsZero() {
			user.UpdatedAt = time.Now()
		}

		if user.Password == "" {
			user.Password = tools.HashPassword("888888", tools.Salt)
		}

		if len(user.Roles) == 0 {
			user.Roles = []primitive.ObjectID{}
		}

		if len(user.PlaceIds) == 0 {
			user.PlaceIds = []primitive.ObjectID{}
		}

		result, err := collection.InsertOne(ctx, user)
		if err != nil {
			return gin.H{"error": "Failed to create user"}, err
		}

		return result, nil
	} else {
		return gin.H{"error": "User with the same username already exists"}, nil
	}
}

func UpdateSysUserByID(updateData *models.SysUsers) (interface{}, error) {
	collection, ctx := GetSysUsersCollection()

	filter := bson.M{"_id": updateData.ID, "status": 1}

	var checkSysUser models.SysUsers
	err := collection.FindOne(ctx, filter).Decode(&checkSysUser)
	if err != nil {
		return nil, err
	}

	if checkSysUser.Status == 1 {

		updateData.UpdatedAt = time.Now()

		result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": updateData})
		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		return gin.H{"error": "User not found or already inactive"}, nil
	}
}

func InactiveUserByID(id string) (interface{}, error) {
	collection, ctx := GetSysUsersCollection()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID, "status": 1}

	var user models.SysUsers

	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	if user.Status == 1 {
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
		return gin.H{"error": "User not found or already inactive"}, nil
	}
}

func GetOneSysUserById(id string) (interface{}, error) {
	collection, ctx := GetSysUsersCollection()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID, "status": 1}

	var user models.SysUsers
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	if user.Status == 0 {
		return gin.H{"error": "User not found or already inactive"}, nil
	} else {
		return &user, nil
	}
}

func GetUserByUsername(username string) (*models.SysUsers, error) {
	collection, ctx := GetSysUsersCollection()

	filter := bson.M{"username": username, "status": 1}

	var user models.SysUsers

	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserDetailByUsername(username string) (interface{}, error) {
	collection, ctx := GetSysUsersCollection()

	filter := bson.M{"username": username, "status": 1}

	var user models.SysUsers

	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	var roleLists []models.SysRoles
	roleLists, err = findRoleListByUser(user.Roles)

	var locationLists []models.Locations
	locationLists, err = GetLocationsByIds(user.PlaceIds)

	departments, err := GetDepartmentsByIds(user.DeptIds)
	var loginRecords []models.LoginRecords
	loginRecords, err = GetLoginRecords(user.Username)

	if err != nil {
		return nil, err
	}

	last := gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"avatarBase64": user.AvatarBase64,
		"roles":        user.Roles,
		"deptIds":      user.DeptIds,
		"placeIds":     user.PlaceIds,
		"roleLists":    roleLists,
		"loginRecords": loginRecords,
		"superAdmin":   user.SuperAdmin,
		"locations":    locationLists,
		"departments":  departments,
	}

	return last, nil
}

func UpdateUserPassword(Username string, NewPassword string) (interface{}, error) {
	collection, ctx := GetSysUsersCollection()

	filter := bson.M{"username": Username, "status": 1}

	var user models.SysUsers
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	if user.Status == 1 {
		result, err := collection.UpdateOne(ctx, filter, bson.M{
			"$set": bson.M{
				"password":  tools.HashPassword(NewPassword, tools.Salt),
				"updatedAt": time.Now(),
			},
		})

		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		return gin.H{"error": "User not found or already inactive"}, nil
	}
}

func SysUserList(pageDto *dto.SysUserPageDto) (interface{}, error) {
	collection, ctx := GetSysUsersCollection()

	if pageDto.Page < 1 {
		pageDto.Page = 1
	}

	if pageDto.Limit < 1 {
		pageDto.Limit = 10
	}

	filter := bson.M{"status": 1}

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

	var results []models.SysUsers
	for cursor.Next(ctx) {
		var user models.SysUsers
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		results = append(results, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return gin.H{"lists": results, "total": count, "page": pageDto.Page, "limit": pageDto.Limit}, nil

}

func UserUpdateAvatar(Username string, PhotoBase string) (interface{}, error) {
	collection, ctx := GetSysUsersCollection()

	filter := bson.M{"username": Username, "status": 1}

	var user models.SysUsers
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	if user.Status == 1 {
		result, err := collection.UpdateOne(ctx, filter, bson.M{
			"$set": bson.M{
				"avatarBase64": PhotoBase,
				"createdAt":    time.Now(),
			},
		})

		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		return gin.H{"error": "This user has been invalidated! Please contact admin!"}, nil
	}
}
