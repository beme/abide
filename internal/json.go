package internal

// UpdateKeyValuesInMap updates every instance of a key within an arbitrary
// `map[string]interface{}` with the given value.
func UpdateKeyValuesInMap(key string, value interface{}, m map[string]interface{}) map[string]interface{} {
	return updateMap(key, value, m)
}

// Recursively update the map to update the specified key.
func updateMap(key string, value interface{}, m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		switch v.(type) {
		// If slice, iterate through each entry and call updateMap
		// only if it's a map[string]interface{}.
		case []interface{}:
			for i := range v.([]interface{}) {
				switch v.([]interface{})[i].(type) {
				case map[string]interface{}:
					v.([]interface{})[i] = updateMap(key, value, v.([]interface{})[i].(map[string]interface{}))
				}
			}
		case map[string]interface{}:
			m[k] = updateMap(key, value, v.(map[string]interface{}))
		default:
			if k == key {
				m[k] = value
			}
			break
		}
	}

	return m
}
