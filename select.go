package gosqler

import (
	"fmt"
	"strings"

	"github.com/tingxin/go-sqler/sql"
)

// Select used to express select sql
type Select interface {
	sql.Where
	Select(fields ...string)
	From(tableNames ...string)
	Orderby(field string, desc bool)
	GroupBy(fields ...string)
	Limit(count int)
	Offset(count int)
	String() string
}

// Where used to build where condition
type selecter struct {
	sql.Wherer
	selectCache []string
	fromCache   []string
	groupCache  []string
	orderCache  []string
	limit       int
	offset      int
}

// NewSelect used to create where object
func NewSelect() Select {
	query := &selecter{}

	query.offset = 0
	query.limit = -1
	return query
}

func (p *selecter) Select(fields ...string) {
	if p.selectCache == nil {
		p.selectCache = fields
	} else {
		p.selectCache = append(p.selectCache, fields...)
	}
}

func (p *selecter) From(tableNames ...string) {
	if p.fromCache == nil {
		p.fromCache = tableNames
	} else {
		p.fromCache = append(p.fromCache, tableNames...)
	}
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

	whereStr := p.Wherer.String()
	if whereStr != "" {
		cache = append(cache, whereStr)
	}

	if p.groupCache != nil {
		groupStr := strings.Join(p.groupCache, ",")
		cache = append(cache, "GROUP BY")
		cache = append(cache, groupStr)
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
