package strings

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	sqlCharPattern *regexp.Regexp
	escapeRep      map[string]string

	operators map[string]bool
)

func init() {
	initSQLCharPattern()
	initOperators()

}

func initSQLCharPattern() {
	rep := make(map[string]string)
	rep["'"] = "‘"
	rep["\""] = "‘"
	rep[";"] = "；"
	rep[","] = "，"
	rep["<"] = "＜"
	rep[">"] = "＞"
	rep["("] = "（"
	rep[")"] = "）"
	rep["@"] = "＠"
	rep["="] = "＝"
	rep["+"] = "＋"
	rep["*"] = "＊"
	rep["&"] = "＆"
	rep["#"] = "＃"
	rep["$"] = "￥"

	escapeRep = make(map[string]string)
	keys := make([]string, len(rep), len(rep))
	var index = 0
	for key, v := range rep {
		escapeKey := regexp.QuoteMeta(key)
		escapeRep[escapeKey] = v
		keys[index] = escapeKey
		index++
	}

	patternStr := strings.Join(keys, "|")
	pattern := regexp.MustCompile(patternStr)
	sqlCharPattern = pattern
}

func initOperators() {
	operators = make(map[string]bool)
	operators[">"] = true
	operators["<"] = true
	operators[">="] = true
	operators["<="] = true
	operators["!"] = true
	operators["!>"] = true
	operators["LIKE"] = true
	operators["like"] = true

	operators["="] = true
	operators["!="] = true
	operators["IS"] = true
	operators["IS NOT"] = true
	operators["is"] = true
	operators["is not"] = true
	operators["in"] = true
	operators["IN"] = true
	operators["not in"] = true
	operators["NOT IN"] = true
}

// FilterSQL used to filter the unsupport string in sql
func FilterSQL(sql string) string {
	return sqlCharPattern.ReplaceAllStringFunc(sql, func(match string) string {
		return escapeRep[match]
	})
}

// GetSQLStr used to decorate sql condition value string
func GetSQLStr(v interface{}) string {
	if v == nil {
		return "null"
	}
	switch v.(type) {
	case string:
		strV := v.(string)
		part := FilterSQL(strV)
		part = fmt.Sprintf("\"%s\"", part)
		return part
	case int8, int, uint16, int32, int64, uint, uint8, bool, float32, float64:
		return fmt.Sprintf("%v", v)
	case []string:
		values := v.([]string)
		return handelStrValues(values)
	case []int:
		values := v.([]int)
		return handelIntValues(values)
	case []int32:
		values := v.([]int32)
		return handelInt32Values(values)
	case []int64:
		values := v.([]int64)
		return handelInt64Values(values)
	case []int8:
		values := v.([]int8)
		return handelInt8Values(values)
	case []bool:
		values := v.([]bool)
		return handelBoolValues(values)
	default:
		return ""

	}
}

// GetSQLOper used to decorate condition operator
func GetSQLOper(operator string, v interface{}) string {
	if _, ok := operators[operator]; !ok {
		panic(fmt.Sprintf("gosqler does not support operator %s", operator))
	}
	if v == nil {
		if operator == "=" {
			return "IS"
		} else if operator == "!=" {
			return "IS NOT"
		}
	}
	return operator
}

func handelStrValues(values []string) string {
	count := len(values)
	t := make([]string, count, count)
	for index, value := range values {
		t[index] = fmt.Sprintf("\"%v\"", value)
	}
	u := strings.Join(t, ",")
	return fmt.Sprintf("(%s)", u)
}

func handelIntValues(values []int) string {
	count := len(values)
	t := make([]string, count, count)
	for index, value := range values {
		t[index] = fmt.Sprintf("%d", value)
	}
	u := strings.Join(t, ",")
	return fmt.Sprintf("(%s)", u)
}

func handelInt32Values(values []int32) string {
	count := len(values)
	t := make([]string, count, count)
	for index, value := range values {
		t[index] = fmt.Sprintf("%d", value)
	}
	u := strings.Join(t, ",")
	return fmt.Sprintf("(%s)", u)
}

func handelInt64Values(values []int64) string {
	count := len(values)
	t := make([]string, count, count)
	for index, value := range values {
		t[index] = fmt.Sprintf("%d", value)
	}
	u := strings.Join(t, ",")
	return fmt.Sprintf("(%s)", u)
}

func handelInt8Values(values []int8) string {
	count := len(values)
	t := make([]string, count, count)
	for index, value := range values {
		t[index] = fmt.Sprintf("%d", value)
	}
	u := strings.Join(t, ",")
	return fmt.Sprintf("(%s)", u)
}

func handelBoolValues(values []bool) string {
	count := len(values)
	t := make([]string, count, count)
	for index, value := range values {
		t[index] = fmt.Sprintf("%v", value)
	}
	u := strings.Join(t, ",")
	return fmt.Sprintf("(%s)", u)
}
