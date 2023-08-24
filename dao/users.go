package dao

import (
	"reflect"
)

// ForeignKey 放在其他结构体第一个元素，内存对齐实现强制类型转换
type ForeignKey struct {
	UserID string
}

type Users struct {
	UserID       string
	Username     string
	Email        string
	PasswordHash string
	OtherInfo    string
}

type Credentials struct {
	UserID       string
	CredentialID string
	Username     string
	PasswordHash string
}
type Tasks struct {
	UserID      string
	TaskID      string
	Username    string
	TaskNum     string
	Title       string
	Description string
	CreatedDate string
	UpdatedDate string
	DueDate     string
	Status      string
}

func StructToMap(obj interface{}) map[string]string {
	result := make(map[string]string)
	objValue := reflect.ValueOf(obj).Elem() // 获取指针指向的成员

	if objValue.Kind() == reflect.Struct {
		objType := objValue.Type()
		for i := 0; i < objValue.NumField(); i++ {
			field := objType.Field(i).Name
			var value string
			if objType.Field(i).Type == reflect.TypeOf(ForeignKey{}) {
				// 某字段为外键时，尝试获取内部的字符串
				value = objValue.Field(i).Field(0).String()
			} else {
				value = objValue.Field(i).String()
			}
			// 只更新非空的字段
			if !objValue.Field(i).IsZero() {
				result[field] = value
			}
		}
	}

	return result
}
