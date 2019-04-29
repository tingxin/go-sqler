
# go-sqler

Easy and safe to build sql
更简单更安全的创建sql 语句

## Build SELECT SQL 创建基本的select 语句

```golang
import (
    gosqler "github.com/tingxin/go-sqler"
    "fmt"
)

query := gosqler.NewSelect()
query.Select("name", "age", "sex", "birthday", "is_employee")
query.From("users")
query.Where("age", "=", 15)
quey.AndWhere("name", "!=", "barry")
query.OrWhere("is_employee", "=", true)
query.OrWhere("birthday", "=", nil)
query.Orderby("birthday", true)
query.Orderby("name", false)
query.Select("address")
query.Limit(10)
query.Offset(36)
query.Where("name", "in", []string{"edwin", "leo", "jack", "stacy"})
query.Where("age", "not in", []int8{16, 24, 32})
str := query.String()
fmt.Println(str)
```

output:

```sql
SELECT name,age,sex,birthday,is_employee,address
FROM users
WHERE age = 15 AND name != "barry" OR is_employee = true OR birthday IS null
AND name in ("edwin","leo","jack","stacy")
AND age not in (16,24,32) ORDER BY birthday DESC,name ASC LIMIT 36, 10
```

## Build INSERT SQL 创建基本的 insert 语句

```golang
import (
    gosqler "github.com/tingxin/go-sqler"
    "fmt"
)

query := gosqler.NewInsert("users")
query.AddColumns("name", "age", "sex", "birthday", "is_employee")
query.AddValues("barry", 30, "male", "1987-01-15", true)
query.AddValues("edwin", 35, "male", "1982-01-15", true)
query.AddValues("stacy", 32.5, "female", "1986-01-15", true)

str := query.String()
fmt.Println(str)
```

output:

```sql
INSERT INTO users ( name,age,sex,birthday,is_employee )
VALUES ("barry",30,"male","1987-01-15",true),
("edwin",35,"male","1982-01-15",true),
("stacy",32.5,"female","1986-01-15",true)
```