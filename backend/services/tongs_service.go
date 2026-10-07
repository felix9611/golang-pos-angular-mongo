package services

import (
	"context"
	"golang-fixedasset-mongo-backend/backend/config"
	"golang-fixedasset-mongo-backend/backend/models"
    "golang-fixedasset-mongo-backend/backend/dto"
	"time"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
    "errors"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/gin-gonic/gin"
    "log"
)

func Foo() error {
    return errors.New("some error happened")
}

func CreateTongs(tongs *models.Tongs) (interface{}, error) {

    filter := bson.M{"status": 1, "name": tongs.Name}

    // Check if a Tong with the same name already exists
    collection := config.GetCollection("tongs")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    count, err := collection.CountDocuments(ctx, filter)
    if err != nil {
        return nil, err
    }

    if count == 0 {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        if tongs.Status == 0 {
            tongs.Status = 1
        }

        if tongs.CreatedAt.IsZero() {
            tongs.CreatedAt = time.Now()
        }

        if tongs.UpdatedAt.IsZero() {
            tongs.UpdatedAt = time.Now()
        }

        result, err := collection.InsertOne(ctx, tongs)
        if err != nil {
            return nil, err
        }
        return result.InsertedID, nil
    } else {
        return "Tong with the same name already exists", nil
    }

}

func GetOneTongsRun(id string) (*models.Tongs, error) {
	collection := config.GetCollection("tongs")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var tongs models.Tongs

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, filter).Decode(&tongs)
	if err != nil {
		return nil, err
	}

	return &tongs, nil
}

func VoidTongs(id string) (interface{}, error) {
    collection := config.GetCollection("tongs")

     objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var tongs models.Tongs

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, filter).Decode(&tongs)

    if err != nil {
        return nil, err
    }
    
    if tongs.Status == 1 {

        result, err := collection.UpdateOne(ctx, filter, bson.M{
            "$set": bson.M{
                "status": 0,
                "updated_at": time.Now(),
            },
        })

        if err != nil {
            return nil, err
        }
        if result.MatchedCount == 0 {
            return nil, mongo.ErrNoDocuments
        }
        return result.ModifiedCount, nil
    } else {
        return "Tongs is already voided", nil
    }
    
}

func UpdateTongsByID(id string, updateData *models.Tongs) (interface{}, error) {
	collection := config.GetCollection("tongs")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var checkTongs models.Tongs

	err = collection.FindOne(ctx, filter).Decode(&checkTongs)
	if err != nil {
		return nil, err
	}

	if checkTongs.Status == 1 {
		if updateData.UpdatedAt.IsZero() {
			updateData.UpdatedAt = time.Now()
		}

		result, err := collection.UpdateOne(ctx, filter, bson.M{
			"$set": updateData,
		})
		if err != nil {
			return nil, err
		}

		if result.MatchedCount == 0 {
			return nil, mongo.ErrNoDocuments
		}

		return result.ModifiedCount, nil
	} else {
		return "Tongs is already voided", nil
	}
}


func TongsListPage(pageDto *dto.TongsPageDto) (interface{}, error) {
    log.Println("TongsListPage called with pageDto:", pageDto)
    if pageDto.Page < 1 {
        pageDto.Page = 1
    }
    if pageDto.Limit < 1 {
        pageDto.Limit = 10
    }

    collection := config.GetCollection("tongs")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"status": 1}
    if pageDto.Name != "" {
        filter["name"] = bson.M{"$regex": pageDto.Name, "$options": "i"}
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

    var results []models.Tongs
    for cursor.Next(ctx) {
        var t models.Tongs
        if err := cursor.Decode(&t); err != nil {
            return nil, err
        }
        results = append(results, t)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return gin.H{"lists": results, "total": count, "page": pageDto.Page, "limit": pageDto.Limit }, nil
}
