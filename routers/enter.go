package routers

import (
	"Smart_delivery_locker/api"
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/middleware"
	"github.com/gin-gonic/gin"
)

type Group struct {
	*gin.RouterGroup
}

func InitRouter() *gin.Engine {
	gin.SetMode(global.Config.System.Env)
	router := gin.Default()

	router.GET("", func(c *gin.Context) {
		c.String(200, "XXX")
	})

	apiGroup := router.Group("api")

	routerGroup := Group{apiGroup}
	routerGroup.userRouter()
	routerGroup.settingsRouter()

	return router
}

// TODO userRouter 用户api的各种功能
func (router Group) userRouter() {
	userApi := api.Api.UserApi
	router.GET("users", middleware.JwtAuth(), userApi.UserListView)
	router.POST("user_login", userApi.LoginView)
	router.POST("user_create", userApi.UserCreateView)
	router.DELETE("user_remove", userApi.UserRemoveView)
	router.PUT("user_update_password", middleware.JwtAuth(), userApi.UserUpdatePasswordView)
}

func (router Group) settingsRouter() {
	settingsApi := api.Api.SettingsApi
	router.GET("settings/:name", settingsApi.SettingsInfoView)
	router.PUT("settings", settingsApi.SettingsInfoUpdateView)
	router.PUT("settings_jwt", settingsApi.SettingsJwtUpdateView)
	router.PUT("settings_admin", settingsApi.SettingsAdminUpdateView)
}
