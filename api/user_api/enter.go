package user_api

import (
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models"
	"Smart_delivery_locker/models/ctype"
	"Smart_delivery_locker/models/ctype/status"
	"Smart_delivery_locker/models/res"
	CODE "Smart_delivery_locker/models/res/code"
	"Smart_delivery_locker/service"
	"Smart_delivery_locker/service/common"
	"Smart_delivery_locker/service/user_ser"
	"Smart_delivery_locker/utils"
	"Smart_delivery_locker/utils/jwts"
	"Smart_delivery_locker/utils/pwd"
	"fmt"
	"time"

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

	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

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

		if ctype.Role(claims.Role) != ctype.PermissionAdmin {
			//非管理员
			user.Username = ""
		}
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

type LoginDataResponse struct {
	Token   string       `json:"token"`
	Profile LoginProfile `json:"profile"`
}

type LoginProfile struct {
	UserID     uint       `json:"user_id"`
	UserName   string     `json:"username"`
	NickName   string     `json:"nick_name"`
	Phone      string     `json:"phone"`
	Role       ctype.Role `json:"role"`
	Permission int        `json:"permission"`
	Status     string     `json:"status"`
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

	if err := global.DB.Model(&userModel).Update("token", token).Error; err != nil {
		global.Log.Error("更新Token失败", err)
	}

	profile := LoginProfile{
		UserID:     userModel.ID,
		UserName:   userModel.Username,
		NickName:   userModel.Username,
		Phone:      userModel.Phone,
		Role:       userModel.Permission,
		Permission: int(userModel.Permission),
		Status:     userModel.Status,
	}

	loginData := LoginDataResponse{
		Token:   token,
		Profile: profile,
	}

	now := time.Now()
	isoStr := utils.ToISO8601(now)
	if err := global.DB.Model(&userModel).Update("last_login_at", isoStr).Error; err != nil {
		global.Log.Error("更新用户登录时间失败", err)
	}

	if err != nil {
		global.Log.Error("token生成失败", err)
		res.ResultFailWithMsg("token生成失败", c)
		return
	}
	res.ResultOK(loginData, fmt.Sprintf("用户%s登录成功", userModel.Username), c)
}

type UserCreateRequest struct {
	Username   string     `json:"username" binding:"required" msg:"请输入用户名"`  //用户名
	Email      string     `json:"email"`                                     //邮箱
	Phone      string     `json:"phone"`                                     //手机号
	Password   string     `json:"password" binding:"required" msg:"请输入密码"`   //密码
	Permission ctype.Role `json:"permission" binding:"required" msg:"请选择权限"` //权限
}

// UserCreateView 控制台方面的用户建立
func (UserApi) UserCreateView(c *gin.Context) {
	var cr UserCreateRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.ResultFailWithError(err, &cr, c)
		return
	}
	err := user_ser.UserService{}.CreateUser(cr.Username, cr.Password, cr.Permission, cr.Email, cr.Phone)
	if err != nil {
		global.Log.Error(err)
		res.ResultFailWithMsg("用户创建失败!", c)
		return
	}

	res.ResultOkWithMsg(fmt.Sprintf("用户%s创建成功!", cr.Username), c)
	return
}

type UserCreateWebRequest struct {
	Username string            `json:"username"` //	是	账号
	Nickname string            `json:"nickname"` //	是	昵称
	Phone    string            `json:"phone"`    //是	手机号
	Email    string            `json:"email"`    //是	邮箱
	Role     ctype.Role        `json:"role"`     //是	admin / courier / user
	Status   status.UserStatus `json:"status"`   //是	enabled / disabled
	Password string            `json:"password"` //是	初始密码
	Avatar   string            `json:"avatar"`   //否	头像
}

// UsersCreateFormWebView Web方面的用户建立
func (UserApi) UsersCreateFormWebView(c *gin.Context) {
	var cr UserCreateWebRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.ResultFailWithError(err, &cr, c)
	}

	// 寻找是否存在手机号相同用户
	var userModel models.User
	err := global.DB.Take(&userModel, "phone = ?", cr.Phone).Error
	if err != nil {
		global.Log.Warn("用户已存在，请重新输入")
		res.ResultFailWithCode(CODE.ArgumentError, c)
		return
	}
	hashPassword := pwd.HashPassword(cr.Password)

	err = global.DB.Create(&models.User{
		Username:   cr.Username,
		Nickname:   cr.Nickname,
		Phone:      cr.Phone,
		Email:      cr.Email,
		Password:   hashPassword,
		Avatar:     cr.Avatar,
		Status:     cr.Status.String(),
		Permission: cr.Role,
	}).Error
	if err != nil {
		global.Log.Error(err)
		res.ResultFailWithMsg("角色创建失败", c)
		return
	}

	err = global.DB.Find(&userModel, "phone = ?", cr.Phone).Error
	if err != nil {
		global.Log.Error(err)
		res.ResultFailWithCode(CODE.ArgumentError, c)
		return
	}
	userModel.Password = "" // 不返回密码
	res.ResultOkWithData(userModel, c)
}

func (UserApi) UserRemoveView(c *gin.Context) {
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.ResultFailWithCode(CODE.ArgumentError, c)
		return
	}
	var userList []models.User
	count := global.DB.Find(&userList, cr.IDList).RowsAffected
	if count == 0 {
		res.ResultFailWithMsg("用户不存在", c)
		return
	}
	err = global.DB.Delete(&userList).Error
	if err != nil {
		global.Log.Error(err)
		res.ResultFailWithMsg("用户删除失败", c)
		return
	}
	res.ResultOkWithMsg(fmt.Sprintf("成功删除%d个用户", count), c)
}

type UpdateUserUri struct {
	ID uint `json:"id"`
}

func (UserApi) UserUpdateView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	var cr UpdateUserUri
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.ResultFailWithCode(CODE.ArgumentError, c)
		return
	}
	var userModel models.User
	// 本人操作
	err = global.DB.Take(&userModel, claims.UserID).Error
	if err != nil {
		res.ResultFailWithMsg("用户不存在", c)
		return
	}
	// 同ID 或者 管理员操作
	if cr.ID == claims.UserID || ctype.PermissionAdmin == ctype.Role(claims.Role) {
		err = global.DB.Model(userModel).Where("id = ?", cr.ID).Updates(userModel).Error
		if err != nil {
			global.Log.Error(err)
			res.ResultFailWithMsg("修改失败", c)
			return
		}
		global.DB.Find(&userModel, "id = ?", cr.ID)
		userModel.Password = ""
		res.ResultOkWithData(userModel, c)
	}
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" msg:"请输入密码"`  //旧密码
	NewPassword string `json:"new_password" binding:"required" msg:"请输入新密码"` //新密码
}

// UserUpdatePasswordView 修改登录人的密码
func (UserApi) UserUpdatePasswordView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)
	fmt.Println("claims", claims)
	// 参数绑定
	var cr UpdatePasswordRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.ResultFailWithError(err, &cr, c)
		return
	}

	var user models.User
	err := global.DB.Take(&user, claims.UserID).Error
	if err != nil {
		res.ResultFailWithMsg("用户不存在", c)
		return
	}
	// 判断密码是否一致
	if !pwd.ComparePasswords(user.Password, cr.OldPassword) {
		res.ResultOkWithMsg("密码错误", c)
		return
	}
	hashPwd := pwd.HashPassword(cr.NewPassword)
	err = global.DB.Model(&user).Update("password", hashPwd).Error
	if err != nil {
		global.Log.Error(err)
		res.ResultFailWithMsg("密码修改失败", c)
		return
	}
	res.ResultOkWithMsg("密码修改成功", c)
	return
}

// LogoutView 登出
func (UserApi) LogoutView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	token := c.Request.Header.Get("token")

	err := service.ServiceApp.UserService.Logout(claims, token)

	if err != nil {
		global.Log.Error(err)
		res.ResultFailWithMsg("注销失败", c)
		return
	}

	res.ResultOkWithMsg("注销成功", c)
}
