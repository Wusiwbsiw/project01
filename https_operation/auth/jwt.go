package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// !!! 警告：这是一个示例密钥。在生产环境中，
// !!! 这个值必须从安全的地方（如环境变量）加载，并且应该更长、更复杂。
var jwtKey = []byte("my_super_secret_key_that_is_long_and_secure")

// AppClaims 是我们自定义的 JWT Payload 结构
type AppClaims struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}

// GenerateToken 为指定的用户 ID 生成一个新的 JWT
func GenerateToken(userID int64) (string, error) {
	// 设置 Token 的过期时间，例如 24 小时
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建我们的自定义 Claims
	claims := &AppClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置过期时间
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			// 设置签发时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// 设置签发方
			Issuer: "https://my-awesome-app.com",
		},
	}

	// 使用 HS256 签名方法创建一个新的 Token 对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用我们定义的密钥来为 Token 签名，并获取完整的 Token 字符串
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证一个 JWT 字符串的有效性
// 如果 Token 有效，它会返回解析出的 AppClaims；否则返回错误
func ValidateToken(tokenString string) (*AppClaims, error) {
	claims := &AppClaims{}

	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法是我们期望的 HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// 返回我们的密钥
		return jwtKey, nil
	})

	// 检查解析过程中是否发生错误 (例如格式错误、签名不匹配)
	if err != nil {
		return nil, err
	}

	// 检查 Token 是否有效 (例如，是否已过期)
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// 一切正常，返回解析出的 Claims
	return claims, nil
}
