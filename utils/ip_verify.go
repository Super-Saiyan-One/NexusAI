package utils

func IsIPAllowed(clientIP string, allowedIPs []string, disallowedIPs []string) bool {
	// 如果存在黑名单，检查IP是否在黑名单中
	if len(disallowedIPs) > 0 {
		disallowedIPsMap := SliceToMap(disallowedIPs)
		if _, exists := disallowedIPsMap[clientIP]; exists {
			return false // IP在黑名单中
		}
		return true // IP不在黑名单中
	}

	// 如果存在白名单，检查IP是否在白名单中
	if len(allowedIPs) > 0 {
		allowedIPsMap := SliceToMap(allowedIPs)
		if _, exists := allowedIPsMap[clientIP]; exists {
			return true // IP在白名单中
		}
		return false // IP不在白名单中
	}

	return true // 既没有黑名单也没有白名单
}
