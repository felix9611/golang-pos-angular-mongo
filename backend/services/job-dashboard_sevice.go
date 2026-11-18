package services

import (
	"context"
	"golang-pos-mongo-backend/backend/config"
	"golang-pos-mongo-backend/backend/dto"
	"golang-pos-mongo-backend/backend/models"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetJobDashboardCollections() *mongo.Collection {
	return config.GetCollection("job_dashboards")
}

func CreateTasks(jobDashboard *models.JobDashboard) (interface{}, error) {
	collection := GetJobDashboardCollections()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"title": jobDashboard.Title,
	}

	var existingJobDashboard models.JobDashboard
	err := collection.FindOne(ctx, filter).Decode(&existingJobDashboard)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			newNumber, err := CreateJobNumber()
			if err != nil {
				return nil, err
			}

			jobDashboard.Number = newNumber
			jobDashboard.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
			jobDashboard.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
			_, errCr := collection.InsertOne(ctx, jobDashboard)
			if errCr != nil {
				return nil, errCr
			}

			CreateActionRecord("Job Dashboard Create", "POST", "Job Dashboard", jobDashboard, "Success")
			return gin.H{"message": "Job Dashboard created successfully"}, nil
		} else {
			CreateActionRecord("Job Dashboard Create", "POST", "Job Dashboard", jobDashboard, "Failed")
			return gin.H{"error": "Job Dashboard already exists"}, nil
		}
	}

	CreateActionRecord("Job Dashboard Create", "POST", "Job Dashboard", jobDashboard, "Failed")
	return gin.H{"error": "Job Dashboard already exists"}, nil
}

func CreateJobNumber() (string, error) {
	collection := GetJobDashboardCollections()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$addFields", Value: bson.D{
			{Key: "numberInt", Value: bson.D{
				{Key: "$convert", Value: bson.D{
					{Key: "input", Value: "$number"},
					{Key: "to", Value: "int"},
					{Key: "onError", Value: 0},
					{Key: "onNull", Value: 0},
				}},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "maxNumber", Value: bson.D{{Key: "$max", Value: "$numberInt"}}},
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
	return formatNumberNew(maxNumber, 8), nil
}

func ListJobDashboard(req dto.ListJobDashboardRequestDto, user dto.AuthUser) (gin.H, error) {

	userFilter := bson.M{}

	// Only restrict normal users
	if !user.SuperAdmin && !user.Manager {
		userFilter = bson.M{
			"$or": bson.A{
				bson.M{
					"userIds": bson.M{
						"$in": bson.A{user.ID},
					},
				},
				bson.M{
					"createdBy": user.Username,
				},
			},
			"jobValid": true,
		}
	}

	// Kanban statuses
	statuses := []struct {
		Key    string
		Status int
	}{
		{"toDoList", 1},
		{"inProgressList", 2},
		{"pendingReviewList", 3},
		{"inReviewList", 4},
		{"completedList", 5},
	}

	result := gin.H{}

	for _, item := range statuses {

		jobs, err := GetJobDashboardListQuery(
			req,
			item.Status,
			userFilter,
		)

		if err != nil {
			return nil, err
		}

		result[item.Key] = jobs
	}

	return result, nil
}

func GetJobDashboardListQuery(req dto.ListJobDashboardRequestDto, status int, userFilter bson.M) ([]bson.M, error) {
	collection := GetJobDashboardCollections()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// --------------------------
	// 1. Match Query
	// --------------------------

	allQuery := bson.M{
		"status": status,
	}

	if req.Name != "" {
		allQuery["title"] = bson.M{
			"$regex":   regexp.QuoteMeta(req.Name),
			"$options": "i",
		}
	}

	for key, value := range userFilter {
		allQuery[key] = value
	}

	// --------------------------
	// 2. Aggregation Pipeline
	// --------------------------

	pipeline := mongo.Pipeline{

		// Match
		bson.D{
			{Key: "$match", Value: allQuery},
		},

		// Lookup SysUsers
		bson.D{
			{Key: "$lookup", Value: bson.M{
				"from": "sys_users",
				"let": bson.M{
					"userIds": bson.M{
						"$map": bson.M{
							"input": "$userIds",
							"as":    "userId",
							"in": bson.M{
								"$toObjectId": "$$userId",
							},
						},
					},
				},

				"pipeline": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{
								"$in": bson.A{
									"$_id",
									"$$userIds",
								},
							},
						},
					},
					bson.M{
						"$project": bson.M{
							"_id":          0,
							"username":     1,
							"avatarBase64": 1,
						},
					},
				},

				"as": "users",
			}},
		},

		// Sort
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "sortOrder", Value: 1},
			}},
		},

		// Project
		bson.D{
			{Key: "$project", Value: bson.M{

				"_id": 1,

				// Comments Count
				"commentsCount": bson.M{
					"$size": bson.M{
						"$filter": bson.M{
							"input": bson.M{
								"$ifNull": bson.A{
									"$comments",
									bson.A{},
								},
							},
							"as": "comment",
							"cond": bson.M{
								"$eq": bson.A{
									"$$comment.status",
									1,
								},
							},
						},
					},
				},

				// Attachments Count
				"attachmentsCount": bson.M{
					"$size": bson.M{
						"$filter": bson.M{
							"input": bson.M{
								"$ifNull": bson.A{
									"$attachments",
									bson.A{},
								},
							},
							"as": "attachment",
							"cond": bson.M{
								"$eq": bson.A{
									"$$attachment.status",
									1,
								},
							},
						},
					},
				},

				// Finished Checklists Count
				"checklistsFinishedCount": bson.M{
					"$size": bson.M{
						"$filter": bson.M{
							"input": bson.M{
								"$ifNull": bson.A{
									"$checkLists",
									bson.A{},
								},
							},
							"as": "checkList",
							"cond": bson.M{
								"$and": bson.A{
									bson.M{
										"$eq": bson.A{
											"$$checkList.status",
											1,
										},
									},
									bson.M{
										"$eq": bson.A{
											"$$checkList.checked",
											true,
										},
									},
								},
							},
						},
					},
				},

				// Active Checklists Count
				"checklistsCount": bson.M{
					"$size": bson.M{
						"$filter": bson.M{
							"input": bson.M{
								"$ifNull": bson.A{
									"$checkLists",
									bson.A{},
								},
							},
							"as": "checkList",
							"cond": bson.M{
								"$eq": bson.A{
									"$$checkList.status",
									1,
								},
							},
						},
					},
				},

				// Job Fields
				"title":       1,
				"number":      1,
				"description": 1,
				"category":    1,
				"tags":        1,
				"labels":      1,
				"jobEndDate":  1,
				"createdAt":   1,
				"jobValid":    1,
				"status":      1,
				"sortOrder":   1,
				"userIds":     1,
				"users":       1,
			}},
		},
	}

	// --------------------------
	// 3. Execute Aggregation
	// --------------------------

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make([]bson.M, 0)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetTaskById(id string) (interface{}, error) {
	collection := GetJobDashboardCollections()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var task models.JobDashboard
	err = collection.FindOne(ctx, filter).Decode(&task)

	if err != nil {
		return gin.H{"error": "Not Found"}, nil
	}

	return &task, nil
}

func UpdateTask(task *models.JobDashboard) (interface{}, error) {
	collection := GetJobDashboardCollections()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": task.ID}

	var existingTask models.JobDashboard
	err := collection.FindOne(ctx, filter).Decode(&existingTask)
	if err != nil {
		return nil, err
	}

	if existingTask.Status == 0 {
		CreateActionRecord("Job Dashboard Update", "POST", "Job Dashboard", task, "Failed")
		return gin.H{"error": "This task is voided!"}, nil
	} else {
		task.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": task})
		if err != nil {
			return nil, err
		}

		CreateActionRecord("Job Dashboard Update", "POST", "Job Dashboard", task, "Success")

		return result, nil
	}
}

func UpdateJobDashboardStatus(task *dto.UpdateJobDashboardStatusDto) (interface{}, error) {
	collection := GetJobDashboardCollections()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": task.ID}

	var existingTask models.JobDashboard
	err := collection.FindOne(ctx, filter).Decode(&existingTask)
	if err != nil {
		return nil, err
	}

	if existingTask.Status == 0 {
		CreateActionRecord("Job Dashboard Update", "POST", "Job Dashboard", task, "Failed")
		return gin.H{"error": "This task is voided!"}, nil
	} else {

		finalData := bson.M{
			"$set": bson.M{
				"status":    task.Status,
				"updatedAt": time.Now().Format("2006-01-02 15:04:05"),
				"sortOrder": task.SortOrder,
			},
		}

		result, err := collection.UpdateOne(ctx, filter, finalData)
		if err != nil {
			return nil, err
		}

		CreateActionRecord("Job Dashboard Update", "POST", "Job Dashboard", task, "Success")

		return result, nil
	}
}

func VoidTaskById(id string) (interface{}, error) {
	collection := GetJobDashboardCollections()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var task models.JobDashboard
	err = collection.FindOne(ctx, filter).Decode(&task)

	if err != nil {
		return gin.H{"error": "Not Found"}, nil
	}

	if task.Status == 0 {
		CreateActionRecord("Job Dashboard Void", "POST", "Job Dashboard", task, "Failed")
		return gin.H{"error": "This task is voided!"}, nil
	} else {
		task.Status = 0
		task.JobValid = false
		task.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": task})
		if err != nil {
			return nil, err
		}

		CreateActionRecord("Job Dashboard Void", "POST", "Job Dashboard", task, "Success")

		return result, nil
	}
}

func SyncTasks(user *dto.AuthUser) (interface{}, error) {
	collection := GetJobDashboardCollections()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"jobValid": true}

	if !(user.SuperAdmin || user.Manager) {
		filter["$or"] = bson.A{
			bson.M{
				"userIds": bson.M{
					"$in": bson.A{user.ID},
				},
			},
			bson.M{
				"createdBy": user.Username,
			},
		}
	}

	pipeline := mongo.Pipeline{
		// Match
		bson.D{
			{Key: "$match", Value: filter},
		},

		//project
		bson.D{
			{
				Key: "$project",
				Value: bson.M{
					"_id":    1,
					"title":  1,
					"status": 1,
					"name":   1,
				},
			},
		},

		//sort
		bson.D{
			{Key: "$sort", Value: bson.M{"sortOrder": 1}},
		},

		//limit
		bson.D{
			{Key: "$limit", Value: 10},
		},
	}

	// --------------------------
	// Execute Aggregation
	// --------------------------

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make([]bson.M, 0)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
