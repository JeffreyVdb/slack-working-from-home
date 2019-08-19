package util

func ContainsString(list []string, item string) bool {
	for _, currentItem := range list {
		if currentItem == item {
			return true
		}
	}

	return false
}
