package services

import (
	"golang-fixedasset-mongo-backend/backend/config"
	"golang-fixedasset-mongo-backend/backend/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func AddLoginRecord(username string, ip string, loginStatus string) (interface{}, error) {
	collection := config.GetCollection("login_records")

	saveData := bson.M{
		"username":  username,
		"ipAddress": ip,
		"loginStatus": loginStatus,
		"loginTime":  time.Now(),
	}

	res, err := collection.InsertOne(context.Background(), saveData)


	if err != nil {		
		return nil, err
	} else {
		return res, nil
	}
	
}

func GetLoginRecords(username string) ([]models.LoginRecords, error) {
	collection := config.GetCollection("login_records")

	filter := bson.M{"username": username}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.LoginRecords
	for cursor.Next(ctx) {
		var record models.LoginRecords
		if err := cursor.Decode(&record); err != nil {
			return nil, err
		}
		results = append(results, record)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}