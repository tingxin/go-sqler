package sql

// Wherer used to express sql where condition
type Wherer interface {
	Where(field, operator string, value interface{})
	AndWhere(field, operator string, value interface{})
	OrWhere(field, operator string, value interface{})
}
