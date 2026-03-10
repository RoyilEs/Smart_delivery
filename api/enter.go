package api

import (
	"Smart_delivery_locker/api/settings_api"
	"Smart_delivery_locker/api/user_api"
)

type ApiGroup struct {
	UserApi     user_api.UserApi
	SettingsApi settings_api.SettingsApi
}

var Api = new(ApiGroup)
