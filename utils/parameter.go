package utils

func GetContextValue(contextValue interface{}, defaultValue any) any {
	value := contextValue
	if value == nil {
		return defaultValue
	}
	return value
}

func GetMaxInt(number1 int, number2 int) int { // 获取两个整数中的最大值
	if number1 > number2 {
		return number1
	}
	return number2
}

func SliceToMap(slice []string) map[string]struct{} { // 将切片转换为map
	m := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		m[s] = struct{}{}
	}
	return m
}
