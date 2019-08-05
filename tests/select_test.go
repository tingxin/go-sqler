package tests

import (
	"fmt"
	"strings"
	"testing"

	gosqler "github.com/tingxin/go-sqler"
)

func TestSelectBasc(t *testing.T) {
	query := gosqler.NewSelect()
	query.Select("name", "age", "sex", "birthday", "is_employee")
	query.From("users")
	expected := `SELECT name,age,sex,birthday,is_employee,address FROM users`
	queryStr := query.String()

	result := compareStr(expected, queryStr)

	if result {
		t.Errorf("\nExpected:\n%s\nActually:\n%s\n", expected, queryStr)
	}
}

func TestSelectBasc1(t *testing.T) {
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
	queryStr := query.String()
	expected := `SELECT name,age,sex,birthday,is_employee,address 
	FROM users 
	WHERE age = 15 AND name != "barry" 
	OR is_employee = true OR birthday IS null 
	OR name like "%go%" AND name in ("edwin","leo","jack","stacy") 
	AND age not in (16,24,32) ORDER BY birthday DESC,name ASC LIMIT 36, 10`

	result := compareStr(expected, queryStr)

	if result {
		t.Errorf("\nExpected:\n%s\nActually:\n%s\n", expected, queryStr)
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

	expected := `SELECT city,education,AVG(age) as avg_age FROM people 
	INNER JOIN orders ON orders.account = people.id AND orders.time = people.birthday 
	LEFT JOIN vip ON vip.account = people.id WHERE age > 10 AND job like "%it%" 
	AND birthday > "1988-09-12 12:12:12" AND address IS NOT null 
	AND is_employee = true GROUP BY city,education 
	ORDER BY avg_age DESC LIMIT 8, 10`
	queryStr := query.String()

	result := compareStr(expected, queryStr)

	if result {
		t.Errorf("\nExpected:\n%s\nActually:\n%s\n", expected, queryStr)
	}
}

func TestSelectWithGroupAndHaving(t *testing.T) {

	query := gosqler.NewSelect()
	df := "DATE_FORMAT(%s.%s,'%s') as dt"
	createTime := fmt.Sprintf(df, "c", "create_time", "%Y-%m")
	query.Choice(createTime)
	query.Choice("c.visibility as vi")
	query.Choice("SUM(c.image_count) as sum_images")
	query.From("comment as c")

	query.BeginWhere()

	query.BeginWhere()
	query.Where("c.id", ">", 441690)
	query.OrWhere("c.id", "<", 241690)
	query.EndWhere()

	query.BeginWhere()
	query.Where("c.like_count", ">", 100)
	query.OrWhere("c.like_count", "<", 10)
	query.EndWhere()

	query.EndWhere()

	query.OrWhere("c.visibility", "in", []string{"banned", "invisible"})

	query.GroupBy("dt", "vi")
	query.Having("sum_images", ">", 100)
	query.Orderby("c.id", true)

	expected := `select DATE_FORMAT(c.create_time, '%Y-%m') as dt, c.visibility as vi, 
	SUM(c.image_count) as sum_images 
	from   comment as c 
	where ((c.id > 441690 or c.id < 241690) 
	and (c.like_count > 100 
	or c.like_count < 10)) 
	or c.visibility in ('banned','invisible')  
	group by dt, vi
	having sum_images > 100
	order by c.id desc`

	queryStr := query.String()

	result := compareStr(expected, queryStr)

	if result {
		t.Errorf("\nExpected:\n%s\nActually:\n%s\n", expected, queryStr)
	}
}

func compareStr(expected, actually string) bool {
	v1 := strings.Replace(expected, "\n", "", -1)
	v1 = strings.Replace(v1, " ", "", -1)

	v1 = strings.ToLower(v1)

	v2 := strings.Replace(actually, "\n", "", -1)
	v2 = strings.Replace(v2, " ", "", -1)
	v2 = strings.ToLower(v2)
	return v1 == v2
}
