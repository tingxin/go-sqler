package sql

import (
	"fmt"
	"strings"

	str "github.com/tingxin/go-sqler/strings"
)

// Filter used to express sql condition
type Filter interface {
	And(field, operator string, value interface{})
	OR(field, operator string, value interface{})
	BeginGroup()
	EndGroup()
	String() string
}

// filter used to build condition
type filter struct {
	key        string
	conditions []string
	groupIndex int16
	braceIndex int16
}

// NewFilter used to create a new Filter
func NewFilter(key string) Filter {
	f := filter{}
	f.key = key
	return &f
}

// And used to add  and condition
func (p *filter) And(key, operator string, value interface{}) {

	if p.conditions == nil {
		p.conditions = make([]string, 1, 1)
		p.conditions[0] = p.key
	} else {
		p.conditions = append(p.conditions, "AND")
	}

	var i int16
	for ; i < p.groupIndex; i++ {
		p.conditions = append(p.conditions, "(")
	}

	p.braceIndex += p.groupIndex
	p.groupIndex = 0
	v := str.GetSQLStr(value)
	op := str.GetSQLOper(operator, value)
	item := fmt.Sprintf("%s %s %s", key, op, v)
	p.conditions = append(p.conditions, item)

}

// OR used to add  OR condition
func (p *filter) OR(key, operator string, value interface{}) {

	if p.conditions == nil {
		p.conditions = make([]string, 1, 1)
		p.conditions[0] = "OR"
	} else {
		p.conditions = append(p.conditions, "OR")
	}

	var i int16
	for ; i < p.groupIndex; i++ {
		p.conditions = append(p.conditions, "(")
	}

	p.braceIndex += p.groupIndex
	p.groupIndex = 0
	v := str.GetSQLStr(value)
	op := str.GetSQLOper(operator, value)
	item := fmt.Sprintf("%s %s %s", key, op, v)
	p.conditions = append(p.conditions, item)

}

// BeginGroup used to begin a condition group
func (p *filter) BeginGroup() {
	p.groupIndex++
}

// EndGroup used to end a condition group
func (p *filter) EndGroup() {
	if p.braceIndex <= 0 {
		panic("Need call begin_group first")
	}

	p.conditions = append(p.conditions, ")")
	p.braceIndex--
}

// String
func (p *filter) String() string {
	if p.conditions != nil && len(p.conditions) > 0 {
		return strings.Join(p.conditions, " ")
	}
	return ""
}
