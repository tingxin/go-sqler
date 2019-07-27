package sql

// Haver used to express sql where condition
type Haver interface {
	Having(field, operator string, value interface{})
	AndHaving(field, operator string, value interface{})
	OrHaving(field, operator string, value interface{})
}
