package tests

import (
	"fmt"
	"testing"

	gosqler "github.com/tingxin/go-sqler"
)

func TestSelect1(t *testing.T) {
	query := gosqler.NewSelect()
	query.Select("name", "age", "sex", "birthday", "is_employee")
	query.From("users")
	query.Where("age", "=", 15)
	query.AndWhere("name", "!=", "barry")
	query.OrWhere("is_employee", "=", true)
	query.OrWhere("birthday", "=", nil)
	query.OrWhere("name", "like", "%go%")
	query.Orderby("birthday", true)
	query.Orderby("name", false)
	query.Select("address")
	query.Limit(10)
	query.Offset(36)
	query.Where("name", "in", []string{"edwin", "leo", "jack", "stacy"})
	query.Where("age", "not in", []int8{16, 24, 32})
	str := query.String()
	expected := `SELECT name,age,sex,birthday,is_employee,address FROM users WHERE age = 15 AND name != "barry" OR is_employee = true OR birthday IS null OR name like "%go%" AND name in ("edwin","leo","jack","stacy") AND age not in (16,24,32) ORDER BY birthday DESC,name ASC LIMIT 36, 10`
	if str != expected {
		t.Errorf("\nExpected:\n%s\nActually:\n%s\n", expected, str)
	}
}

func TestSelectWithJoin(t *testing.T) {
	query := gosqler.NewSelect()
	query.Select("city", "education", "AVG(age) as avg_age")
	query.From("people")
	query.Where("age", ">", 10)

	query.Join("orders", "orders.account = people.id",
		"orders.time = people.birthday")
	query.AndWhere("job", "like", "%it%")
	query.AndWhere("birthday", ">", "1988-09-12 12:12:12")
	query.AndWhere("address", "!=", nil)
	query.AndWhere("is_employee", "=", true)

	query.LeftJoin("vip", "vip.account = people.id")

	query.GroupBy("city", "education")
	query.Orderby("avg_age", true)
	query.Limit(10)
	query.Offset(8)

	expected := `SELECT city,education,AVG(age) as avg_age FROM people INNER JOIN orders ON orders.account = people.id AND orders.time = people.birthday LEFT JOIN vip ON vip.account = people.id WHERE age > 10 AND job like "%it%" AND birthday > "1988-09-12 12:12:12" AND address IS NOT null AND is_employee = true GROUP BY city,education ORDER BY avg_age DESC LIMIT 8, 10`
	queryStr := query.String()
	fmt.Printf(queryStr)

	if queryStr != expected {
		t.Errorf("\nExpected:\n%s\nActually:\n%s\n", expected, queryStr)
	}
}
