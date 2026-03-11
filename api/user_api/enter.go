package user_api

import (
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models"
	"Smart_delivery_locker/models/ctype"
	"Smart_delivery_locker/models/res"
	CODE "Smart_delivery_locker/models/res/code"
	"Smart_delivery_locker/service/common"
	"Smart_delivery_locker/utils"
	"Smart_delivery_locker/utils/jwts"
	"Smart_delivery_locker/utils/pwd"
	"fmt"
	"github.com/gin-gonic/gin"
)

type UserApi struct{}

type UserResponse struct {
	models.User
	RoleID int `json:"role_id"`
}

type UserListRequest struct {
	models.PageInfo
	Permission int `json:"permission" form:"permission"`
}

func (UserApi) UserListView(c *gin.Context) {

	//TODO 正式使用jwt后 断言
	//_claims, _ := c.Get("claims")
	//claims := _claims.(*jwts.CustomClaims)

	var page UserListRequest
	if err := c.ShouldBind(&page); err != nil {
		res.ResultFailWithCode(CODE.ArgumentError, c)
		return
	}

	var users []UserResponse
	list, count, _ := common.ComList(models.User{Permission: ctype.Role(page.Permission)}, common.Option{
		PageInfo: page.PageInfo,
	})

	for _, user := range list {

		//if ctype.Role(claims.Role) != ctype.PermissionAdmin {
		//	//非管理员
		//	user.Username = ""
		//}
		// 脱敏
		user.Phone = utils.DesensitizationTel(user.Phone)
		user.Email = utils.DesensitizationEmail(user.Email)
		users = append(users, UserResponse{
			User:   user,
			RoleID: int(user.Permission),
		})
	}

	res.ResultOkWithList(users, count, c)
}

type LoginRequest struct {
	UserName string `json:"username" binding:"required" msg:"请输入用户名"`
	Password string `json:"password" binding:"required" msg:"请输入密码"`
}

func (UserApi) LoginView(c *gin.Context) {
	var cr LoginRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.ResultFailWithError(err, &cr, c)
		return
	}

	var userModel models.User
	err = global.DB.Take(&userModel, "username = ?", cr.UserName).Error
	if err != nil {
		global.Log.Warn("用户不存在")
		res.ResultFailWithMsg("用户不存在", c)
		return
	}

	// 密码验证
	ok := pwd.ComparePasswords(userModel.Password, cr.Password)
	if !ok {
		global.Log.Warn("密码错误")
		res.ResultFailWithMsg("密码错误", c)
		return
	}

	//登录成功生成token
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		Username: userModel.Username,
		Role:     int(userModel.Permission),
		UserID:   userModel.ID,
		Avatar:   userModel.Avatar,
	})
	if err != nil {
		global.Log.Error("token生成失败", err)
		res.ResultFailWithMsg("token生成失败", c)
		return
	}
	res.ResultOK(token, fmt.Sprintf("用户%s登录成功", userModel.Username), c)
}
