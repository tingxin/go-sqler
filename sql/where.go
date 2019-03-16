package sql

import (
	"fmt"
	"strings"

	str "github.com/tingxin/go-sqler/strings"
)

// Where used to express sql where condition
type Where interface {
	Where(field, operator string, value interface{})
	AndWhere(field, operator string, value interface{})
	OrWhere(field, operator string, value interface{})
}

// Wherer used to build where condition
type Wherer struct {
	cache []string
}

// Where used to add where condition
func (p *Wherer) Where(field, operator string, value interface{}) {
	v := str.GetSQLStr(value)
	op := str.GetSQLOper(operator, value)
	if p.cache == nil {
		condition := fmt.Sprintf("%s %s %s", field, op, v)
		p.cache = make([]string, 2, 2)
		p.cache[0] = "WHERE"
		p.cache[1] = condition
	} else {
		condition := fmt.Sprintf("AND %s %s %s", field, op, v)
		p.cache = append(p.cache, condition)
	}
}

// AndWhere used to add where condition
func (p *Wherer) AndWhere(field, operator string, value interface{}) {
	p.add("AND", field, operator, value)
}

// OrWhere used to add where condition
func (p *Wherer) OrWhere(field, operator string, value interface{}) {
	p.add("OR", field, operator, value)
}

func (p *Wherer) add(prefix, field, operator string, value interface{}) {
	if p.cache == nil {
		panic("need call Where method first")
	}
	v := str.GetSQLStr(value)
	op := str.GetSQLOper(operator, value)
	condition := fmt.Sprintf("%s %s %s %s", prefix, field, op, v)
	p.cache = append(p.cache, condition)
}

func (p *Wherer) String() string {
	whereStr := strings.Join(p.cache, " ")
	return whereStr
}
