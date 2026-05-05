package routers

import (
	"Smart_delivery_locker/api"
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Group struct {
	*gin.RouterGroup
}

func InitRouter() *gin.Engine {
	gin.SetMode(global.Config.System.Env)
	router := gin.Default()
	router.Use(Cors())

	router.GET("", func(c *gin.Context) {
		c.String(200, "XXX")
	})

	apiGroup := router.Group("api")

	routerGroup := Group{apiGroup}
	routerGroup.userRouter()
	routerGroup.settingsRouter()
	routerGroup.itemRouter()
	routerGroup.grilleRouter()

	return router
}

// TODO userRouter 用户api的各种功能
func (router Group) userRouter() {
	userApi := api.Api.UserApi
	router.GET("users", middleware.JwtAuth(), userApi.UserListView)
	router.POST("user_login", userApi.LoginView)
	router.POST("user_create", userApi.UserCreateView)
	router.POST("users", userApi.UsersCreateFormWebView) // 区分命令行建立用户 此处检测Phone的差别
	router.DELETE("user_remove", userApi.UserRemoveView)
	router.DELETE("users/:id", middleware.JwtAdmin(), userApi.UserDeleteView)
	router.PUT("user_update_password", middleware.JwtAuth(), userApi.UserUpdatePasswordView)
	router.PUT("users/:id/reset_password", userApi.ResetPasswordView)
	router.PUT("users/:id", middleware.JwtAuth(), userApi.UserUpdateView)
}

func (router Group) settingsRouter() {
	settingsApi := api.Api.SettingsApi
	router.GET("settings/:name", middleware.JwtAdmin(), settingsApi.SettingsInfoView)
	router.PUT("settings", middleware.JwtAdmin(), settingsApi.SettingsInfoUpdateView)
	router.PUT("settings_jwt", middleware.JwtAdmin(), settingsApi.SettingsJwtUpdateView)
	router.PUT("settings_admin", middleware.JwtAdmin(), settingsApi.SettingsAdminUpdateView)
}

// TODO itemRouter 订单api的各种功能
func (router Group) itemRouter() {
	itemApi := api.Api.ItemApi
	router.GET("items/:name", itemApi.ItemListView)
	router.GET("user_items/:name", itemApi.ItemUserListView)
	router.POST("item_create", itemApi.ItemCreateView)
}

// TODO grilleRouter 格口api的各种功能
func (router Group) grilleRouter() {
	grilleApi := api.Api.GrilleApi
	router.POST("grille_form_item_create", grilleApi.GrilleFormItemCreateView)
	router.POST("grille_create", grilleApi.GrilleCreateView)
	router.POST("item_out_grille", grilleApi.ItemOutGrilleView)
	router.GET("grille_phone_get_item/:phone", grilleApi.PhoneGetItemsView)
}

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许前端localhost:5173
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS,PATCH")
		// 必须放行你所有自定义请求头！
		c.Header("Access-Control-Allow-Headers", "Content-Type,token,Authorization,X-User-Scene")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 预检OPTIONS请求直接放行200
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}
