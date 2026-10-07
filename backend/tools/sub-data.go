package tools

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Helper to convert different map representations into a standard map[string]interface{}
func GetSubMap(val interface{}) map[string]interface{} {
	switch v := val.(type) {
	case bson.M:
		return map[string]interface{}(v)
	case map[string]interface{}:
		return v
	case primitive.D:
		return v.Map()
	}
	return nil
}

// Helper to safely fetch a string from a nested map
func GetNestedString(item interface{}, parentKey, childKey string) string {
	var parentMap map[string]interface{}

	switch m := item.(type) {
	case bson.M:
		parentMap = GetSubMap(m[parentKey])
	case map[string]interface{}:
		parentMap = GetSubMap(m[parentKey])
	}

	if parentMap != nil {
		if strVal, ok := parentMap[childKey].(string); ok {
			return strVal
		}
	}
	return ""
}
