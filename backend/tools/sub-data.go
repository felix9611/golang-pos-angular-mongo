package tools

import (
	"fmt"
	"strconv"

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
func GetNestedValue(item interface{}, parentKey, childKey string) interface{} {
	parentVal := getFieldValue(item, parentKey)
	if parentVal == nil {
		return nil
	}

	return getFieldValue(parentVal, childKey)
}

// Helper to safely fetch a string from a nested map
func getFieldValue(item interface{}, key string) interface{} {
	switch m := item.(type) {
	case bson.M:
		return m[key]
	case map[string]any:
		return m[key]
	case primitive.D:
		for _, elem := range m {
			if elem.Key == key {
				return elem.Value
			}
		}
	}
	return nil
}

func GetFloat64(item interface{}, parentKey, childKey string) interface{} {
	parent := getFieldValue(item, parentKey)

	if parent == nil {
		return float64(0)
	}

	value := parent

	// If childKey is provided, get nested value
	if childKey != "" {
		value = getFieldValue(parent, childKey)
	}

	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case primitive.Decimal128:
		f, err := strconv.ParseFloat(v.String(), 64)
		if err == nil {
			return f
		}
	}

	return float64(0)
}

func ConvertStringsToObjectIDs(idStrings []string) ([]primitive.ObjectID, error) {
	objectIDs := make([]primitive.ObjectID, len(idStrings))

	for i, idStr := range idStrings {
		objID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid ObjectID at index %d (%s): %w", i, idStr, err)
		}
		objectIDs[i] = objID
	}

	return objectIDs, nil
}

func ClaimStrings(value interface{}) []primitive.ObjectID {
	result := []primitive.ObjectID{}

	values, ok := value.([]interface{})
	if !ok {
		return result
	}

	for _, v := range values {
		if s, ok := v.(string); ok {
			id, err := primitive.ObjectIDFromHex(s)
			if err == nil {
				result = append(result, id)
			}
		}
	}

	return result
}
