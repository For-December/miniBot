package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log"
)

var databaseName string
var username string
var password string
var host string
var port string

var db *sql.DB

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file:", err)
		return // 自动退出
	}

	databaseName = viper.GetString("mysql.databaseName")
	username = viper.GetString("mysql.username")
	password = viper.GetString("mysql.password")
	host = viper.GetString("mysql.host")
	port = viper.GetString("mysql.port")

	dataSourceName := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + databaseName
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
}

func IsRegistered(userId string) bool {
	query, err := db.Query("select * from users where userid = ?", userId)
	if err != nil {
		log.Fatal(err)
		return false
	}
	if query.Next() {
		return true
	}
	return false
}

func UpdateTable(tableName string, kv map[string]string, userId string, resist map[string]string) bool {
	if !IsRegistered(userId) {
		WarningF("用户 %s 未注册！", userId)
		return false
	}

	insertStatement := "update " + tableName + " set "
	values := make([]any, 0)
	// 拼接要修改的字段和值
	for key, value := range kv {
		insertStatement += key + " = ? ,"
		values = append(values, value)
	}
	// 切片前闭后开，这里去掉逗号
	insertStatement = insertStatement[0 : len(insertStatement)-1]
	insertStatement += "where userid = ?"
	values = append(values, userId)

	// where 后面的限制条件，map 为空则不会遍历
	for key, value := range resist {
		insertStatement += " and " + key + " = ?"
		values = append(values, value)
	}

	Info(insertStatement)
	result, err := db.Exec(insertStatement, values...)
	if err != nil {
		Error("出错了: ", err)
	}
	// 获取插入操作的结果
	affectedRows, err := result.RowsAffected()
	if err != nil {
		Error("出错了: ", err)
	}
	if affectedRows > 0 {
		Info("数据更新成功！")
		return true
	}
	Warning("数据更新失败，有可能新旧数据相同！")
	return false

}

func CreateUsers(userId string, userName string, email string, passwordHash string, otherInfo string) bool {
	insertStatement := "insert into users values (?,?,?,?,?)"
	result, err := db.Exec(insertStatement, userId, userName, email, passwordHash, otherInfo)
	if err != nil {
		Error("出错了: ", err)
	}
	// 获取插入操作的结果
	affectedRows, err := result.RowsAffected()
	if err != nil {
		Error("出错了: ", err)
	}
	if affectedRows > 0 {
		Info("新用户添加成功！")
		return true
	}
	Warning("新用户添加失败！")
	return false
}

//func UpdateEmail(userId string, value string) bool {
//	return updateUser(userId, "Email", value)
//}
//
//func UpdateUsername(userId string, value string) bool {
//	return updateUser(userId, "Username", value)
//}
//
//func UpdatePasswordHash(userId string, value string) bool {
//	return updateUser(userId, "PasswordHash", value)
//}
//func UpdateOtherInfo(userId string, value string) bool {
//	return updateUser(userId, "OtherInfo", value)
//}

func Test() {

	// 查询数据
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// 遍历查询结果
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", tableName)
	}

	// 检查查询错误
	if err = rows.Err(); err != nil {
		panic(err)
	}
}
