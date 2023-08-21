package utils

import (
	"regexp"
)

const taskPattern = "(date:\\s*\\S+\\s+" +
	"(email:\\s*\\S+@\\S+\\s+)?" + // 邮箱可选，问号表示前面括号内的内容可选
	"title:\\s*\\S+\\s+" +
	"context:\\s*```[\\S\\s]+?```)+" // 匹配整个任务一次或多次，实现同时设置多个任务
const datePattern = "date:\\s*(\\S+)\\s+"
const emailPattern = "email:\\s*(\\S+@\\S+)\\s+"
const titlePattern = "title:\\s*(\\S+)\\s+"
const contextPattern = "context:\\s*```([\\S\\s]+?)```"

func IsTaskTxt(txt string) bool {

	Info(txt)
	isMatch, err := regexp.MatchString(taskPattern, txt)
	if err != nil {
		return false
	}
	return isMatch
}

func GetTaskParam(txt string) (params map[string][]string) {

	getParam := func(pattern string) (res []string) {
		matches := regexp.
			MustCompile(pattern).
			FindAllStringSubmatch(txt, -1)
		for _, elem := range matches {
			Warning(elem[1])
			res = append(res, elem[1])
		}
		return
	}

	// map必须初始化，无法在返回值直接插值
	params = map[string][]string{
		"date":    getParam(datePattern),
		"email":   getParam(emailPattern),
		"title":   getParam(titlePattern),
		"context": getParam(contextPattern),
	}

	// 遍历匹配结果
	return
}
