package utils

func Contains[T string](list []T, target T) bool {

	for _, a := range list {
		if a == target {
			return true
		}
	}

	return false

}
