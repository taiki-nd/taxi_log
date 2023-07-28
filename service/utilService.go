package service

func IsEmptyString(s string) bool {
	return len(s) == 0
}

func IsValidHour(hour int64) bool {
	return 0 <= hour && hour <= 24
}

func IsValidRate(rate float64) bool {
	return 0 <= rate && rate <= 100
}
