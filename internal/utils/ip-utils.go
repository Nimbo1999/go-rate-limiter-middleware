package utils

import "strings"

func FormattIp(ip string) string {
	if strings.Contains(ip, ":") {
		return strings.Split(ip, ":")[0]
	}
	return ip
}
