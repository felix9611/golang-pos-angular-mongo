package services

import (
	"context"
	"fmt"
	"golang-fixedasset-mongo-backend/backend/config"
	"golang-fixedasset-mongo-backend/backend/dto"
	"golang-fixedasset-mongo-backend/backend/models"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateSysMenu(data *models.SysMenus) (interface{}, error) {
	filter := bson.M{"status": 1, "name": data.Name, "path": data.Path}

	collection := config.GetCollection("sys_menus")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if data.Status == 0 {
			data.Status = 1
		}

		if data.CreatedAt.IsZero() {
			data.CreatedAt = time.Now()
		}

		if data.UpdatedAt.IsZero() {
			data.UpdatedAt = time.Now()
		}
		result, err := collection.InsertOne(ctx, data)
		if err != nil {
			return nil, err
		}

		CreateActionRecord("Menu Item Create", "POST", "System Menu", data, "Success")

		return result.InsertedID, nil
	} else {
		CreateActionRecord("Menu Item Create", "POST", "System Menu", data, "Failed")
		return "Menu with the same name already exists", nil
	}

}

func GetOneMenuItemById(id string) (interface{}, error) {
	var menuItem models.SysMenus
	collection := config.GetCollection("sys_menus")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return menuItem, err
	}

	filter := bson.M{"_id": objectID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, filter).Decode(&menuItem)
	if err != nil {
		return nil, err
	}

	if menuItem.Status == 0 {
		return "Menu item inactive", nil
	} else {
		return menuItem, nil
	}
}

func VoidMenuItemById(id string) (interface{}, error) {
	collection := config.GetCollection("sys_menus")

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var menuItem models.SysMenus

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, filter).Decode(&menuItem)

	if err != nil {
		return nil, err
	}

	if menuItem.Status == 0 {
		CreateActionRecord("Menu Item Void", "DELETE", "System Menu", menuItem, "Failed")
		return "Menu item already voided", nil
	} else {

		menuItem.Status = 0
		menuItem.UpdatedAt = time.Now()

		result, err := collection.UpdateOne(ctx, filter, bson.M{
			"$set": menuItem,
		})

		if err != nil {
			return nil, err
		}

		CreateActionRecord("Menu Item Void", "DELETE", "System Menu", menuItem, "Success")

		return result, nil
	}
}

func UpdateMenuItem(data *models.SysMenus) (interface{}, error) {
	collection := config.GetCollection("sys_menus")

	filter := bson.M{"_id": data.ID}

	var menuItem models.SysMenus
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&menuItem)
	if err != nil {
		return nil, err
	}

	if menuItem.Status == 0 {
		CreateActionRecord("Menu Item Update", "POST", "System Menu", data, "Failed")
		return "This Menu Item is inactive", nil
	} else {
		data.UpdatedAt = time.Now()

		result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": data})

		CreateActionRecord("Menu Item Update", "POST", "System Menu", data, "Success")

		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func ListAllMainIdMenu() (interface{}, error) {
	collection := config.GetCollection("sys_menus")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Aggregation pipeline
	pipeline := mongo.Pipeline{
		// $match: mainId = ""
		{{Key: "$match", Value: bson.M{"mainId": ""}}},
		// $group by mainId + name, keep first _id
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"mainId": "$mainId",
				"name":   "$name",
			},
			"_idValue": bson.M{"$first": "$_id"},
		}}},
		// $project final fields
		{{Key: "$project", Value: bson.M{
			"_id":    "$_idValue",
			"mainId": "$_id.mainId",
			"name":   "$_id.name",
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []dto.SysMenuMainId
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetAllMenuItems() (interface{}, error) {
	collection := config.GetCollection("sys_menus")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"status": 1}}},
		{{Key: "$sort", Value: bson.M{"sort": 1}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var menus []models.SysMenus
	if err := cursor.All(ctx, &menus); err != nil {
		return nil, err
	}

	return menus, nil

}

// ListAllMenu fetch menus by name and build hierarchical tree
func ListAllMenu(query *dto.SysMenuList) (interface{}, error) {
	collection := config.GetCollection("sys_menus")

	filter := bson.M{}
	if query.Name != "" {
		filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: fmt.Sprintf(".*%s.*", query.Name), Options: "i"}}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var menus []dto.SysMenuChildrens
	if err := cursor.All(ctx, &menus); err != nil {
		return nil, err
	}

	// Build hierarchical tree
	tree := BuildSortedTree(menus)
	return tree, nil
}

// BuildSortedTree builds a hierarchical and sorted tree from flat menu list
func BuildSortedTree(data []dto.SysMenuChildrens) []*dto.SysMenuChildrens {
	nodeMap := make(map[string]*dto.SysMenuChildrens, len(data))

	for i := range data {
		item := data[i]
		item.Childrens = []*dto.SysMenuChildrens{} // 初始化空 Slice
		nodeMap[item.ID.Hex()] = &item
	}

	var tree []*dto.SysMenuChildrens

	for _, item := range data {
		current := nodeMap[item.ID.Hex()]
		if item.MainId != "" {
			if parent, exists := nodeMap[item.MainId]; exists {
				parent.Childrens = append(parent.Childrens, current)
			} else {
				tree = append(tree, current)
			}
		} else {
			tree = append(tree, current)
		}
	}

	var sortTree func(nodes []*dto.SysMenuChildrens)
	sortTree = func(nodes []*dto.SysMenuChildrens) {
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Sort < nodes[j].Sort
		})

		for _, node := range nodes {
			if len(node.Childrens) > 0 {
				sortTree(node.Childrens)
			}
		}
	}

	sortTree(tree)
	return tree
}

func GetMenusByIds(query *dto.GetMenusByIds) (interface{}, error) {
	collection := config.GetCollection("sys_menus")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	initialIds := []interface{}{}

	oids := []primitive.ObjectID{}

	for _, id := range initialIds {
		strId, ok := id.(string) // interface{} -> string
		if !ok {
			continue
		}

		oid, err := primitive.ObjectIDFromHex(strId)
		if err != nil {
			continue
		}

		oids = append(oids, oid)
	}

	filter1 := bson.M{
		"status": 1,
		"$or": []bson.M{
			{"_id": bson.M{"$in": oids}}, // ✅ 使用 oids
			{"mainId": bson.M{"$in": query.IDS}},
		},
	}

	cursor1, err := collection.Find(ctx, filter1)
	if err != nil {
		return nil, err
	}
	var result1 []models.SysMenus
	if err := cursor1.All(ctx, &result1); err != nil {
		return nil, err
	}

	mainIdSet := make(map[string]struct{})
	for _, r := range result1 {
		if r.MainId != "" {
			mainIdSet[r.MainId] = struct{}{}
		}
	}

	mainIds := []string{}
	for k := range mainIdSet {
		mainIds = append(mainIds, k)
	}

	filter2 := bson.M{
		"status": 1,
		"_id": bson.M{"$in": func() []primitive.ObjectID {
			var oids2 []primitive.ObjectID
			for _, id := range mainIds {
				if oid, err := primitive.ObjectIDFromHex(id); err == nil {
					oids2 = append(oids2, oid)
				}
			}
			return oids2
		}()},
	}

	cursor2, err := collection.Find(ctx, filter2)
	if err != nil {
		return nil, err
	}
	var result2 []models.SysMenus
	if err := cursor2.All(ctx, &result2); err != nil {
		return nil, err
	}

	mergedMap := make(map[string]models.SysMenus)
	for _, r := range append(result1, result2...) {
		mergedMap[r.ID.Hex()] = r
	}

	var merged []dto.SysMenuChildrens
	for _, v := range mergedMap {
		merged = append(merged, dto.SysMenuChildrens{
			ID:                v.ID,
			MainId:            v.MainId,
			Name:              v.Name,
			Icon:              v.Icon,
			Path:              v.Path,
			Sort:              v.Sort,
			Type:              v.Type,
			ExcelFunctionCode: v.ExcelFunctionCode,
			ExcelFunctionName: v.ExcelFunctionName,
			Status:            v.Status,
			CreatedAt:         v.CreatedAt,
			UpdatedAt:         v.UpdatedAt,
			Childrens:         []*dto.SysMenuChildrens{},
		})
	}

	// Build tree
	finalTree := BuildSortedTree(merged)

	return finalTree, nil
}
