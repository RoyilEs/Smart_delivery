package api

import (
	"Smart_delivery_locker/api/grille_api"
	"Smart_delivery_locker/api/item_api"
	"Smart_delivery_locker/api/packages_api"
	"Smart_delivery_locker/api/settings_api"
	"Smart_delivery_locker/api/user_api"
)

type ApiGroup struct {
	UserApi     user_api.UserApi
	SettingsApi settings_api.SettingsApi
	ItemApi     item_api.ItemApi
	GrilleApi   grille_api.GrilleApi
	PackagesApi packages_api.PackageApi
}

var Api = new(ApiGroup)
