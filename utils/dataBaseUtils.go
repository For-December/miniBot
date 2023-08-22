package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testbot/conf"
	"testbot/dao"
	"unsafe"
)

//// mysql 配置
//var mysqlConf struct {
//	databaseName string
//	username     string
//	password     string
//	host         string
//	port         string
//}

var db *sql.DB

func init() {
	databaseName := conf.Config.Mysql.DatabaseName
	username := conf.Config.Mysql.Username
	password := conf.Config.Mysql.Password
	host := conf.Config.Mysql.Host
	port := conf.Config.Mysql.Port

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

//func UpdateUser(users dao.Users) bool {
//
//}
//
//func UpdateCredentials(credentials dao.Credentials) bool {
//
//}
//func UpdateTasks(tasks dao.Tasks) bool {
//
//}

// 根据表名、两个指针 更新数据
// 试了一下午，没找到更好的方法，现在必须两个
func UpdateTableByUserId(tableName string, objPtr unsafe.Pointer, obj interface{}) bool {

	// 得到 foreignKey 指针
	foreignKey := (*dao.ForeignKey)(objPtr)
	if len(foreignKey.UserID) == 0 {
		Warning("未指定UserID！")
		return false
	}

	DebugF("foreignKey: %p", foreignKey)
	DebugF("obj       : %p", obj)
	DebugF("objPtr    : %p", objPtr)

	InfoF("[UserID: %v, table: %v] will be update...",
		len(foreignKey.UserID), tableName)

	// 删除 UserID
	foreignKey.UserID = ""

	return updateTable("users", dao.StructToMap(obj), "1", nil)

}
func updateTable(tableName string,
	kv map[string]string,
	userId string,
	resist map[string]string) bool {
	Info(kv)
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

func CreateUsers(userId string,
	userName string,
	email string,
	passwordHash string,
	otherInfo string) bool {
	if IsRegistered(userId) {
		WarningF("用户 %s 已注册！", userId)
		return false
	}
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
	if affectedRows <= 0 {
		Warning("新用户添加失败！")
		return false
	}

	InfoF("新用户 %s 添加成功！", userName)
	insertStatement = "insert into credentials (userid, username, passwordhash) values (?,?,?)"
	result, err = db.Exec(insertStatement, userId, userName, passwordHash)
	if err != nil {
		Error("出错了: ", err)
	}
	if affectedRows <= 0 {
		Error("用户登录凭据添加失败！")
		return false
	}
	Info("用户登录凭据添加成功！")
	return true

}
