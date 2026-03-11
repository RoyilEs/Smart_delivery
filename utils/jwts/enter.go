package jwts

import (
	"Smart_delivery_locker/global"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
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
