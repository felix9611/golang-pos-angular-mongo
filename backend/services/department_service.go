package services

import (
	"context"
	"golang-fixedasset-mongo-backend/backend/config"
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// "go.mongodb.org/mongo-driver/mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "errors"
	// "log"
)

func BatchCreateDepartment(departments []models.Department) (interface{}, error) {
	for _, department := range departments {
		result, _ := CreateDepartment(&department)

		if result == nil {
			return "failed to create vendor", nil
		}
	}

	return "batch insert completed", nil
}

func CreateDepartment(department *models.Department) (interface{}, error) {

	filter := bson.M{"status": 1, "deptCode": department.DeptCode, "deptName": department.DeptName}

	// Check if a Department with the same name already exists
	collection := config.GetCollection("departments")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if department.Status == 0 {
			department.Status = 1
		}

		if department.CreatedAt.IsZero() {
			department.CreatedAt = time.Now()
		}

		if department.UpdatedAt.IsZero() {
			department.UpdatedAt = time.Now()
		}

		result, err := collection.InsertOne(ctx, department)
		if err != nil {
			return nil, err
		}

		CreateActionRecord("Department Create", "POST", "Department", department, "Success")

		return result, nil
	} else {
		CreateActionRecord("Department Create", "POST", "Department", department, "Failed")
		return "Department with the same name already exists", nil
	}
}

func GetAllDepartments() ([]*models.Department, error) {
	collection := config.GetCollection("departments")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var departments []*models.Department

	cursor, err := collection.Find(ctx, bson.M{"status": 1})
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var department models.Department
		err := cursor.Decode(&department)
		if err != nil {
			return nil, err
		}
		departments = append(departments, &department)
	}

	return departments, nil
}

func GetOneDepartment(id string) (interface{}, error) {
	collection := config.GetCollection("departments")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var department models.Department

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, filter).Decode(&department)
	if err != nil {
		return nil, err
	}

	if department.Status == 0 {
		return "Department is inactive", nil
	} else {
		return &department, nil
	}
}

func VoidDepartmentById(id string) (interface{}, error) {
	collection := config.GetCollection("departments")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var dept models.Department

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, filter).Decode(&dept)

	if err != nil {
		return nil, err
	}

	if dept.Status == 1 {

		dept.Status = 0
		dept.UpdatedAt = time.Now()

		result, err := collection.UpdateOne(ctx, filter, bson.M{
			"$set": dept,
		})

		if err != nil {
			return nil, err
		}

		CreateActionRecord("Department Void", "DELETE", "Department", dept, "Success")

		return result.ModifiedCount, nil
	} else {
		CreateActionRecord("Department Void", "DELETE", "Department", dept, "Failed")
		return "Department is already voided", nil
	}
}

func UpdateDeptById(id string, updateData *models.Department) (interface{}, error) {
	collection := config.GetCollection("departments")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID, "status": 1}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var checkDept models.Department
	err = collection.FindOne(ctx, filter).Decode(&checkDept)
	if err != nil {
		return nil, err
	}

	if checkDept.Status == 1 {
		updateData.UpdatedAt = time.Now()

		result, err := collection.UpdateOne(ctx, filter, bson.M{
			"$set": updateData,
		})
		if err != nil {
			return nil, err
		}

		CreateActionRecord("Department Update", "POST", "Department", updateData, "Success")

		return result, nil
	} else {
		CreateActionRecord("Department Update", "POST", "Department", updateData, "Failed")
		return "Department is already voided", nil
	}
}

func DepartmentListWithFilters(pageDto *dto.DepartmentPageDto) (interface{}, error) {
	collection := config.GetCollection("departments")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1}
	if pageDto.Name != "" {
		filter = bson.M{
			"status": 1,
			"$or": []bson.M{
				{"dept_name": bson.M{"$regex": pageDto.Name, "$options": "i"}},
				{"dept_code": bson.M{"$regex": pageDto.Name, "$options": "i"}},
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

	var results []models.Department
	for cursor.Next(ctx) {
		var department models.Department
		if err := cursor.Decode(&department); err != nil {
			return nil, err
		}
		results = append(results, department)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func DepartmentList(pageDto *dto.DepartmentPageDto) (interface{}, error) {

	if pageDto.Page < 1 {
		pageDto.Page = 1
	}

	if pageDto.Limit < 1 {
		pageDto.Limit = 10
	}

	collection := config.GetCollection("departments")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1}
	if pageDto.Name != "" {
		filter = bson.M{
			"status": 1,
			"$or": []bson.M{
				{"dept_name": bson.M{"$regex": pageDto.Name, "$options": "i"}},
				{"dept_code": bson.M{"$regex": pageDto.Name, "$options": "i"}},
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

	var results []models.Department
	for cursor.Next(ctx) {
		var department models.Department
		if err := cursor.Decode(&department); err != nil {
			return nil, err
		}
		results = append(results, department)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return gin.H{"lists": results, "total": count, "page": pageDto.Page, "limit": pageDto.Limit}, nil

}

func BatchCreateDepartments(departments []models.Department) (interface{}, error) {
	for _, department := range departments {
		result, _ := CreateDepartment(&department)
		if result == nil {
			return "failed to create department", nil
		}
	}

	return "batch insert completed", nil
}

func DepartmentDataFinder(deptCode string, deptName string) (models.Department, error) {
	collection := config.GetCollection("departments")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"status": 1, "deptCode": deptCode, "deptName": deptName}

	var department models.Department

	err := collection.FindOne(ctx, filter).Decode(&department)
	if err != nil {
		return models.Department{}, err
	}

	if department.Status == 0 {
		return models.Department{}, nil
	} else {
		return department, nil
	}
}
