package Messages

import (
	"strings"
)

// вырезать токен из строки, отправленной пользователем
func CutAccessTokenAndUserId(str string) (string, string) {
	runes := []rune(str)
	tokenStr := ""
	userIdStr := ""
	startToken := strings.Index(str, "access_token=")
	endToken := strings.Index(str, "&expires_in")
	startUserId := strings.Index(str, "user_id=")
	if startToken != -1 && endToken != -1 && startUserId != -1 {
		tokenStr = string(runes[startToken+len("access_token=") : endToken])
		userIdStr = string(runes[startUserId+len("user_id="):])
	}
	return tokenStr, userIdStr
}
