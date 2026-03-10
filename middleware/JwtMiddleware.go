package middleware

import (
	"Goblog/models/ctype"
	"Goblog/models/res"
	"Goblog/service/redis_ser"
	"Smart_delivery_locker/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			res.ResultFailWithMsg("未携带token", c)
			c.Abort() //拦截
			return
		}
		claims, err := jwts.ParseToken(token)
		if err != nil {
			res.ResultFailWithMsg("token错误", c)
			c.Abort()
			return
		}
		fmt.Println(claims)
		// 判断是否在redis中
		if redis_ser.CheckLogout(token) {
			res.ResultFailWithMsg("token失效", c)
			c.Abort()
			return
		}
		// 登录的用户
		c.Set("claims", claims)
	}
}

func JwtAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			res.ResultFailWithMsg("未携带token", c)
			c.Abort() //拦截
			return
		}
		claims, err := jwts.ParseToken(token)
		if err != nil {
			res.ResultFailWithMsg("token错误", c)
			c.Abort()
			return
		}
		fmt.Println(claims)
		// 判断是否在redis中
		if redis_ser.CheckLogout(token) {
			res.ResultFailWithMsg("token失效", c)
			c.Abort()
			return
		}
		// 登录的用户
		if ctype.Role(claims.Role) != ctype.PermissionAdmin {
			res.ResultFailWithMsg("权限错误", c)
			c.Abort()
			return
		}
		c.Set("claims", claims)
	}
}
