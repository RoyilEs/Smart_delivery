package jwts

import (
	"Smart_delivery_locker/global"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"strings"
	"time"
)

// JwtPayLoad jwt中包含的payload数据
type JwtPayLoad struct {
	Username string `json:"username"`
	Role     int    `json:"role"`
	UserID   uint   `json:"userID"`
	Avatar   string `json:"avatar"`
}

type CustomClaims struct {
	JwtPayLoad
	jwt.StandardClaims
}

// GenToken 创建 token
func GenToken(user JwtPayLoad) (string, error) {
	var MySecret = []byte(global.Config.Jwt.Secret)
	claim := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * time.Duration(global.Config.Jwt.Expires))),
			Issuer:    global.Config.Jwt.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(MySecret)
}

// ParseToken 解析 token
func ParseToken(tokenStr string) (*CustomClaims, error) {
	var MySecret = []byte(global.Config.Jwt.Secret)
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		global.Log.Error(fmt.Sprintf("token parse error: %s", err.Error()))
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// ParseToken2 解析 token（修复版）
func ParseToken2(tokenStr string) (*CustomClaims, error) {
	// 1. 提前校验 token 格式（空值/非 Bearer 格式等）
	if strings.TrimSpace(tokenStr) == "" {
		global.Log.Error("token parse error: token string is empty")
		return nil, errors.New("token string is empty")
	}

	// 2. 从全局配置读取密钥（避免函数内重复定义）
	MySecret := []byte(global.Config.Jwt.Secret)
	if len(MySecret) == 0 {
		global.Log.Error("token parse error: jwt secret is empty")
		return nil, errors.New("jwt secret is not configured")
	}

	// 3. 解析 token 并验证签名算法
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 关键：验证签名算法是否为预期的 HMAC 算法（比如 HS256）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			errMsg := fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
			global.Log.Error(fmt.Sprintf("token parse error: %s", errMsg))
			return nil, errors.New(errMsg)
		}
		return MySecret, nil
	})

	// 4. 细分解析错误类型
	if err != nil {
		errMsg := fmt.Sprintf("token parse error: %s", err.Error())
		global.Log.Error(errMsg)
		// 区分具体错误类型（过期、签名错误等）
	}

	// 5. 验证 claims 类型和 token 有效性
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	global.Log.Error("token parse error: invalid token")
	return nil, errors.New("invalid token")
}
