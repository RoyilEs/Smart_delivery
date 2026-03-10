package service

import "Smart_delivery_locker/service/user_ser"

type ServiceGroup struct {
	UserService user_ser.UserService
}

var ServiceApp = new(ServiceGroup)
