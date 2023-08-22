package dao

import "reflect"

type Users struct {
	UserID       string
	Username     string
	Email        string
	PasswordHash string
	OtherInfo    string
}
type Credentials struct {
	CredentialID string
	UserID       string
	Username     string
	PasswordHash string
}
type Tasks struct {
	TaskID      string
	UserID      string
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
	objValue := reflect.ValueOf(obj)

	if objValue.Kind() == reflect.Struct {
		objType := objValue.Type()
		for i := 0; i < objValue.NumField(); i++ {
			field := objType.Field(i)
			value := objValue.Field(i).String()

			// 只更新非空的字段
			if !objValue.Field(i).IsZero() {
				result[field.Name] = value
			}
		}
	}

	return result
}
