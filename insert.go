package gosqler

import (
	"fmt"
	"strings"

	str "github.com/tingxin/go-sqler/strings"
)

// Insert used to express insert sql
type Insert interface {
	AddColumns(fields ...string) error
	AddValues(values ...interface{}) error
	ClearValues() error
	String() string
}

// inserter used to build where condition
type inserter struct {
	tableName   string
	columns     []string
	pairs       [][]string
	columnCount int
}

// NewInsert used to create where object
func NewInsert(tableName string) Insert {
	p := &inserter{}
	p.tableName = tableName
	p.pairs = make([][]string, 0)
	p.columnCount = 0
	return p
}

func (p *inserter) AddColumns(fields ...string) error {
	p.columns = fields
	p.columnCount = len(p.columns)
	return nil
}

func (p *inserter) AddValues(values ...interface{}) error {
	count := len(values)
	if p.columnCount == 0 {
		p.columnCount = count
	}
	if p.columnCount > 0 && count != p.columnCount {
		err := fmt.Errorf("The count of values and columns are not match")
		panic(err)
	}

	parts := make([]string, count, count)
	for i, v := range values {
		parts[i] = str.GetSQLStr(v)
	}
	p.pairs = append(p.pairs, parts)
	return nil
}

func (p *inserter) ClearValues() error {
	p.pairs = make([][]string, 0)
	return nil
}

func (p *inserter) String() string {
	cache := make([]string, 1)
	cache[0] = fmt.Sprintf("INSERT INTO %s", p.tableName)

	cache = append(cache, "(")
	cache = append(cache, strings.Join(p.columns, ","))
	cache = append(cache, ") VALUES")

	count := len(p.pairs)
	if count == 0 {
		return ""
	}
	values := make([]string, count)
	for i, p := range p.pairs {
		t := strings.Join(p, ",")
		tquto := fmt.Sprintf("(%s)", t)
		values[i] = tquto
	}
	valueStr := strings.Join(values, ",")
	cache = append(cache, valueStr)
	return strings.Join(cache, " ")
}
