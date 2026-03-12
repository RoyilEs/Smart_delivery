package utils

import (
	"Smart_delivery_locker/global"
	"crypto/md5"
	"encoding/hex"
	"github.com/go-playground/validator/v10"
	"math/rand"
	"reflect"
	"strings"
)

// InList 是否存在列表里面
func InList(key string, list []string) bool {
	for _, s := range list {
		if key == s {
			return true
		}
	}
	return false
}

// Md5 加密
func Md5(str []byte) string {
	m := md5.New()
	m.Write(str)
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

// GetValidMsg 返回结构体中中的msg参数
func GetValidMsg(err error, obj any) string {
	//使用的时候 需要传Obj的指针
	getObj := reflect.TypeOf(obj)
	//将err接口断言为具体类型
	if errs, ok := err.(validator.ValidationErrors); ok {
		//断言成功
		for _, e := range errs {
			//循环每一个错误信息
			//根据报错字段名 获得结构体的具体字段
			if f, exits := getObj.Elem().FieldByName(e.Field()); exits {
				msg := f.Tag.Get("msg")
				return msg
			}
		}
	}
	return err.Error()
}

// DesensitizationTel 手机号脱敏
func DesensitizationTel(tel string) string {
	if len(tel) != 11 {
		return ""
	}
	return tel[:3] + "****" + tel[7:]
}

// DesensitizationEmail  邮箱脱敏
func DesensitizationEmail(email string) string {
	emailList := strings.Split(email, "@")
	if len(emailList) != 2 {
		return ""
	}
	return emailList[0][:2] + "******" + emailList[1]
}

func Random(min, max int) int {
	return rand.Intn(max-min) + min
}

// DeleteByIndex 删除指定索引的元素（高效，O(1)，但会改变原切片顺序）
func DeleteByIndex[T any](slice []T, index int) []T {
	// 索引合法性检查
	if index < 0 || index >= len(slice) {
		global.Log.Warn("索引越界，返回原切片")
		return slice
	}
	// 用最后一个元素覆盖目标索引，截断末尾
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

// DeleteByIndexKeepOrder 删除指定索引的元素（保持顺序，O(n)）
func DeleteByIndexKeepOrder[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		global.Log.Warn("索引越界，返回原切片")
		return slice
	}
	// 拼接前后切片，保持顺序
	return append(slice[:index], slice[index+1:]...)
}

// DeleteByValue 删除指定值的元素（删除所有匹配项）
func DeleteByValue[T comparable](slice []T, value T) []T {
	newSlice := make([]T, 0, len(slice)) // 预分配容量
	for _, v := range slice {
		if v != value {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

// DeleteByCondition 按条件删除元素（通用方法）
func DeleteByCondition[T any](slice []T, condition func(T) bool) []T {
	newSlice := make([]T, 0, len(slice))
	for _, v := range slice {
		if !condition(v) { // 保留不满足删除条件的元素
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}
