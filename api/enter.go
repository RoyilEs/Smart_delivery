package api

import (
	"Smart_delivery_locker/api/item_api"
	"Smart_delivery_locker/api/settings_api"
	"Smart_delivery_locker/api/user_api"
)

type ApiGroup struct {
	UserApi     user_api.UserApi
	SettingsApi settings_api.SettingsApi
	ItemApi     item_api.ItemApi
}

var Api = new(ApiGroup)
