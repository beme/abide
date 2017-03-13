package internal

// UpdateKeyValuesInMap updates every instance of a key within an arbitrary
// `map[string]interface{}` with the given value.
func UpdateKeyValuesInMap(key string, value interface{}, m map[string]interface{}) map[string]interface{} {
	return updateMap(key, value, m)
}

func updateMap(key string, value interface{}, m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		switch v.(type) {
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
