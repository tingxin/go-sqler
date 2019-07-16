package tests

import (
	"testing"

	gosqler "github.com/tingxin/go-sqler"
)

func TestInsert1(t *testing.T) {
	query := gosqler.NewInsert("users")
	query.AddColumns("name", "age", "sex", "birthday", "is_employee")
	query.AddValues("barry", 30, "male", "1987-01-15", true)
	query.AddValues("edwin", 35, "male", "1982-01-15", true)
	query.AddValues("stacy", 32.5, "female", nil, true)
	str := query.String()
	expected := `INSERT INTO users ( name,age,sex,birthday,is_employee ) VALUES ("barry",30,"male","1987-01-15",true),("edwin",35,"male","1982-01-15",true),("stacy",32.5,"female",null,true)`
	if str != expected {
		t.Errorf("\nExpected:\n%s\nActually:\n%s\n", expected, str)
	}
}
