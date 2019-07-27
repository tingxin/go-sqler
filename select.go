package gosqler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tingxin/go-sqler/sql"
)

var (
	// ErrorEmptyParams occur when param is empty
	ErrorEmptyParams = errors.New("params should not be empty")
)

// Select used to express select sql
type Select interface {
	sql.Wherer
	sql.Haver
	Select(fields ...string)
	Choice(field string)
	From(tableNames ...string)
	Join(table string, conditions ...string)
	LeftJoin(table string, conditions ...string)
	RightJoin(table string, conditions ...string)
	FullJoin(table string, conditions ...string)
	Orderby(field string, desc bool)
	GroupBy(fields ...string)
	BeginWhere()
	EndWhere()
	BeginHaving()
	EndHaving()
	Limit(count int)
	Offset(offset int)
	String() string
}

// Where used to build where condition
type selecter struct {
	selectCache []string
	fromCache   []string
	groupCache  []string
	orderCache  []string
	joinCache   []string
	limit       int
	offset      int

	where  sql.Filter
	having sql.Filter
}

// NewSelect used to create where object
func NewSelect() Select {
	query := &selecter{}

	query.offset = 0
	query.limit = -1

	query.having = sql.NewFilter("Having")
	return query
}

func (p *selecter) Select(fields ...string) {
	if p.selectCache == nil {
		p.selectCache = fields
	} else {
		p.selectCache = append(p.selectCache, fields...)
	}
}

func (p *selecter) Choice(field string) {
	if p.selectCache == nil {
		p.selectCache = make([]string, 1, 1)
		p.selectCache[0] = field
	} else {
		p.selectCache = append(p.selectCache, field)
	}
}

func (p *selecter) From(tableNames ...string) {
	if p.fromCache == nil {
		p.fromCache = tableNames
	} else {
		p.fromCache = append(p.fromCache, tableNames...)
	}
}

func (p *selecter) Join(table string, conditions ...string) {
	p.joinBy("INNER", table, conditions...)
}

func (p *selecter) LeftJoin(table string, conditions ...string) {
	p.joinBy("LEFT", table, conditions...)
}

func (p *selecter) RightJoin(table string, conditions ...string) {
	p.joinBy("RIGHT", table, conditions...)
}

func (p *selecter) FullJoin(table string, conditions ...string) {
	p.joinBy("FULL", table, conditions...)
}

// Where used to add where condition
func (p *selecter) Where(field, operator string, value interface{}) {
	if p.where == nil {
		p.where = sql.NewFilter("WHERE")
	}

	p.where.And(field, operator, value)
}

// AndWhere used to add where condition
func (p *selecter) AndWhere(field, operator string, value interface{}) {
	p.where.And(field, operator, value)
}

// OrWhere used to add where condition
func (p *selecter) OrWhere(field, operator string, value interface{}) {
	p.where.OR(field, operator, value)
}

func (p *selecter) BeginWhere() {
	if p.where == nil {
		p.where = sql.NewFilter("WHERE")
	}
	p.where.BeginGroup()
}
func (p *selecter) EndWhere() {
	if p.where == nil {
		panic("Call BeginWhere First")
	}
	p.where.EndGroup()
}

func (p *selecter) Orderby(field string, desc bool) {
	if p.orderCache == nil {
		p.orderCache = make([]string, 1, 1)
		if desc {
			p.orderCache[0] = fmt.Sprintf("%s DESC", field)
		} else {
			p.orderCache[0] = fmt.Sprintf("%s ASC", field)
		}

	} else {
		if desc {
			p.orderCache = append(p.orderCache, fmt.Sprintf("%s DESC", field))
		} else {
			p.orderCache = append(p.orderCache, fmt.Sprintf("%s ASC", field))
		}
	}
}

func (p *selecter) GroupBy(fields ...string) {
	if p.groupCache == nil {
		p.groupCache = fields
	} else {
		p.groupCache = append(p.groupCache, fields...)
	}

}

func (p *selecter) Having(field, operator string, value interface{}) {
	if p.having == nil {
		p.having = sql.NewFilter("HAVING")
	}
	p.having.And(field, operator, value)
}
func (p *selecter) AndHaving(field, operator string, value interface{}) {
	p.having.And(field, operator, value)
}
func (p *selecter) OrHaving(field, operator string, value interface{}) {
	p.having.OR(field, operator, value)
}
func (p *selecter) BeginHaving() {
	if p.having == nil {
		p.having = sql.NewFilter("HAVING")
	}
	p.having.BeginGroup()
}
func (p *selecter) EndHaving() {
	if p.having == nil {
		panic("Call BeginHaving First")
	}
	p.having.EndGroup()
}
func (p *selecter) Limit(count int) {
	p.limit = count
}
func (p *selecter) Offset(count int) {
	p.offset = count
}

func (p *selecter) String() string {
	if p.selectCache == nil || p.fromCache == nil {
		return ""
	}
	cache := make([]string, 4, 4)
	selectStr := strings.Join(p.selectCache, ",")
	cache[0] = "SELECT"
	cache[1] = selectStr

	fromStr := strings.Join(p.fromCache, ",")
	cache[2] = "FROM"
	cache[3] = fromStr

	if len(p.joinCache) > 0 {
		joinStr := strings.Join(p.joinCache, " ")
		cache = append(cache, joinStr)
	}

	whereStr := p.where.String()
	if whereStr != "" {
		cache = append(cache, whereStr)
	}

	if p.groupCache != nil {
		groupStr := strings.Join(p.groupCache, ",")
		cache = append(cache, "GROUP BY")
		cache = append(cache, groupStr)
	}

	havingStr := p.having.String()
	if havingStr != "" {
		cache = append(cache, havingStr)
	}

	if p.orderCache != nil {
		orderStr := strings.Join(p.orderCache, ",")
		cache = append(cache, "ORDER BY")
		cache = append(cache, orderStr)
	}

	if p.limit > 0 {
		var limitStr string
		if p.offset > 0 {
			limitStr = fmt.Sprintf("LIMIT %d, %d", p.offset, p.limit)
		} else {
			limitStr = fmt.Sprintf("LIMIT %d", p.limit)
		}
		cache = append(cache, limitStr)
	}

	return strings.Join(cache, " ")
}

func (p *selecter) joinBy(typeStr, table string, conditions ...string) {
	if table == "" || len(conditions) == 0 {
		panic(ErrorEmptyParams)
	}
	head := fmt.Sprintf("%s JOIN %s ON", typeStr, table)
	body := strings.Join(conditions, " AND ")
	if p.joinCache == nil {
		p.joinCache = make([]string, 2, 2)
		p.joinCache[0] = head
		p.joinCache[1] = body
	} else {
		p.joinCache = append(p.joinCache, head)
		p.joinCache = append(p.joinCache, body)
	}
}
