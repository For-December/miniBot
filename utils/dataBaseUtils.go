package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"testbot/conf"
	"testbot/dao"
	"testbot/logger"
	"time"
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

// UpdateUsers 参数传入后内部属性不会发生变化，深拷贝
func UpdateUsers(users dao.Users) bool {
	// 这里 users 的内部属性会发生变化
	return updateTableByUserId("users", &users, nil)
}

// UpdateCredentials 参数传入后内部属性不会发生变化，深拷贝
func UpdateCredentials(credentials dao.Credentials) bool {
	return updateTableByUserId("credentials", &credentials, nil)
}

// UpdateTasks 参数传入后内部属性不会发生变化，深拷贝
func UpdateTasks(tasks dao.Tasks) bool {
	return updateTableByUserId("tasks", &tasks, nil)
}

// 根据表名、一个指针 更新数据
// 试了一下午，之前要两个指针，现在只要一个，完美！
func updateTableByUserId(tableName string, obj interface{}, rangeLimit map[string]string) bool {
	if reflect.ValueOf(obj).Kind() != reflect.Pointer {
		logger.WarningF("obj 必须传入指针类型，当前类型: %v", reflect.ValueOf(obj))
		return false
	}

	// 得到 foreignKey 指针
	foreignKey := (*dao.ForeignKey)(unsafe.Pointer(reflect.ValueOf(obj).Pointer()))
	if len(foreignKey.UserID) == 0 {
		logger.Warning("未指定UserID！")
		return false
	}

	logger.DebugF("foreignKey: %p", foreignKey)
	logger.DebugF("obj       : %p", obj)
	//DebugF("objPtr    : %p", objPtr)

	logger.InfoF("[UserID: %v, table: %v] will be update...",
		len(foreignKey.UserID), tableName)

	// 暂存并删除 UserID
	userId := foreignKey.UserID
	foreignKey.UserID = ""

	return updateTable("users", dao.StructToMap(obj), userId, rangeLimit)

}

func updateTable(
	tableName string,
	kv map[string]string,
	userId string,
	rangeLimit map[string]string) bool {
	logger.Info(kv)
	if !IsRegistered(userId) {
		logger.WarningF("用户 %s 未注册！", userId)
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

	insertStatement += "where" +
		" userid = ?"
	values = append(values, userId)

	// where 后面的限制条件，map 为空则不会遍历
	for key, value := range rangeLimit {
		insertStatement += " and " + key + " = ?"
		values = append(values, value)
	}

	logger.Info(insertStatement)
	result, err := db.Exec(insertStatement, values...)
	if err != nil {
		logger.Error("出错了: ", err)
	}
	// 获取插入操作的结果
	affectedRows, err := result.RowsAffected()
	if err != nil {
		logger.Error("出错了: ", err)
	}
	if affectedRows > 0 {
		logger.Info("数据更新成功！")
		return true
	}
	logger.Warning("数据更新失败，有可能新旧数据相同！")
	return false

}

func CreateUsers(
	userId string,
	userName string,
	email string,
	passwordHash string,
	otherInfo string) bool {
	if IsRegistered(userId) {
		logger.WarningF("用户 %s 已注册！", userId)
		return false
	}
	insertStatement := "insert into users values (?,?,?,?,?)"
	result, err := db.Exec(insertStatement, userId, userName, email, passwordHash, otherInfo)
	if err != nil {
		logger.Error("出错了: ", err)
	}
	// 获取插入操作的结果
	affectedRows, err := result.RowsAffected()
	if err != nil {
		logger.Error("出错了: ", err)
	}
	if affectedRows <= 0 {
		logger.Warning("新用户添加失败！")
		return false
	}

	logger.InfoF("新用户 %s 添加成功！", userName)
	insertStatement = "insert into credentials (userid, username, passwordhash) values (?,?,?)"
	result, err = db.Exec(insertStatement, userId, userName, passwordHash)
	if err != nil {
		logger.Error("出错了: ", err)
	}
	if affectedRows <= 0 {
		logger.Error("用户登录凭据添加失败！")
		return false
	}
	logger.Info("用户登录凭据添加成功！")
	return true
}

func CreateTasks(
	userId string,
	username string,
	title string,
	description string,
	date string,
	status string) (bool, string, *dao.Tasks) {

	if !IsRegistered(userId) {
		logger.WarningF("用户 [%v, %v] 未注册！", userId, username)
		return false, "用户未注册，待办事项设置失败...", nil
	}

	dateSet, _ := time.ParseInLocation(conf.Config.DateLayout, date, time.Local)
	if dateSet.Sub(time.Now()) < time.Minute*5 {
		logger.WarningF("用户 [%v, %v] 设置的待办事项过早...", userId, username)
		return false, "您的待办事项 时间参数 至少应当在 5 分钟后...", nil
	}

	rows, err := db.Query("select TaskNum from tasks where UserID = ?", userId)
	if err != nil {
		logger.Error("出错了: ", err)
	}

	var taskNum = 1
	for rows.Next() {
		var tempNum int
		err := rows.Scan(&tempNum)
		if err != nil {
			logger.Error("出错了: ", err)
		}
		if taskNum <= taskNum {
			taskNum = tempNum + 1
		}
	}
	// 得到接下来的 taskNum
	dueDate, _ := time.Parse(conf.Config.DateLayout, date)

	// CreatedDate, UpdatedDate 由 mysql 维护
	insertStatement := "insert " +
		"into tasks (UserID, Username, TaskNum, Title, Description, DueDate, Status) " +
		"values (?,?,?,?,?,?,?)"
	result, err := db.Exec(insertStatement, userId, username, taskNum, title, description, dueDate, status)
	if err != nil {
		logger.Error("出错了: ", err)
	}
	// 获取插入操作的结果
	affectedRows, err := result.RowsAffected()
	if err != nil {
		logger.Error("出错了: ", err)
	}
	if affectedRows <= 0 {
		logger.WarningF("用户 %v 的第 %v 个任务添加失败！", username, taskNum)
		return false, "任务添加失败！", nil
	}

	logger.InfoF("用户 %v 的第 %v 个任务添加成功！", username, taskNum)
	return true, "任务添加成功！", &dao.Tasks{
		UserID:      userId,
		Username:    username,
		TaskNum:     fmt.Sprint(taskNum),
		Title:       title,
		Description: description,
		DueDate:     date,
		Status:      status}

}

type limitParam struct {
	key   string
	value string
}

func GetTasks(params ...limitParam) (tasksArray []dao.Tasks) {

	queryStatement := "select " +
		"UserID, TaskNum, Username, Title, Description, CreatedDate, UpdatedDate, DueDate, Status " +
		"from tasks"
	values := make([]any, 0)
	if len(params) != 0 {
		queryStatement += " where"
		for i, param := range params {
			if i == 0 {
				queryStatement += " " + param.key + " = ?"
				values = append(values, param.value)
			} else {
				queryStatement += " and " + param.key + " = ?"
				values = append(values, param.value)
			}
		}
	}
	logger.Debug("query tasks: ", queryStatement)
	rows, err := db.Query(queryStatement, values...)
	if err != nil {
		logger.Error("出错了: ", err)
	}

	tasksArray = make([]dao.Tasks, 0)
	for rows.Next() {
		var tasks dao.Tasks
		err := rows.Scan(&tasks.UserID,
			&tasks.TaskNum,
			&tasks.Username,
			&tasks.Title,
			&tasks.Description,
			&tasks.CreatedDate, &tasks.UpdatedDate, &tasks.DueDate,
			&tasks.Status)

		formatDate := func(originDate string) string {
			if originDate == "" {
				return ""
			}
			datetime, err := time.Parse("2006-01-02 15:04:05", originDate)
			if err != nil {
				logger.Error("出错了：", err)
			}
			return datetime.Format(conf.Config.DateLayout)
		}

		tasks.CreatedDate = formatDate(tasks.CreatedDate)
		tasks.UpdatedDate = formatDate(tasks.UpdatedDate)
		tasks.DueDate = formatDate(tasks.DueDate)

		tasksArray = append(tasksArray, tasks)
		if err != nil {
			logger.Error("出错了: ", err)
		}

	}
	return

}

func GetTasksById(userId string) (tasksArray []dao.Tasks) {
	return GetTasks(limitParam{key: "UserID", value: userId})

}
