package khadijah

// M is a utility shortcut for a map
type M map[string]interface{}

func Contains(items []string, key string) bool {
	for _, s := range items {
		if s == key {
			return true
		}
	}

	return false
}
