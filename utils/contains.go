package utils

func ContainsString(list []string, value string) bool {
	for _, elem := range list {
		if elem == value {
			return true
		}
	}

	return false
}
