// Package stupidqlite is a really dumn sqlite query generator
package stupidqlite

import (
	"fmt"
	"strings"
)

// Query struct that we attach our generate methods to
type Query struct {
	SQL string
}

// Conditional represents a sub-generator for conditionals such as in WHERE clauses
type Conditional func(...string) string

// CondMap is the type used for mapping field sets to conditionals
type CondMap map[string]Conditional

// Eq returns an = conditional clause, if 2 fields are specified they will be compared
// eg: f[0] = f[1] otherwise it assumed that you are comparing field[0] to a placeholder
func Eq(fields ...string) string {
	if len(fields) > 1 {
		return fmt.Sprintf("%s=%s", fields[0], fields[1])
	}
	return fmt.Sprintf("%s=?", fields[0])
}

// Like expects only one argument and always assumes you are comparing to a placeholder
func Like(fields ...string) string {
	return fmt.Sprintf("%s LIKE ?", fields[0])
}

// Between expects two fields to and returns a BETWEEN clause
func Between(fields ...string) string {
	if len(fields) != 2 {
		return ""
	}
	return fmt.Sprintf("BETWEEN %s AND %s", fields[0], fields[1])
}

// NewQuery creates a new query object
func NewQuery() *Query {
	return new(Query)
}

// Insert generates a basic insert statement
func (q *Query) Insert(table string, fields ...string) *Query {
	// Generate placeholders
	var p []string
	for range fields {
		p = append(p, "?")
	}

	q.SQL = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(fields, ","), strings.Join(p, ","))
	return q
}

// Update generates the first part of an UPDATE statement. a Where clause will be necessary
// to complete the SQL query
func (q *Query) Update(table string, fields ...string) *Query {
	q.SQL = fmt.Sprintf("UPDATE %s SET %s", table, strings.Join(fields, "=?, "))
	return q
}

// Select generates a SELECT that can then be chained with further functions
func (q *Query) Select(fields ...string) *Query {
	q.SQL = fmt.Sprintf("SELECT %s", strings.Join(fields, ", "))
	return q
}

// From adds a FROM clause
func (q *Query) From(table string) *Query {
	q.SQL = fmt.Sprintf("%s FROM %s", q.SQL, table)
	return q
}

// Where adds len(fields) WHERE field=?/AND clauses. Multiple field conditionals
// can be defined in the fields map as CondMap{"field1:field2": Condtional}
func (q *Query) Where(fields CondMap) *Query {
	first := true
	for field, cond := range fields {
		if first {
			q.SQL = fmt.Sprintf("%s WHERE %s", q.SQL, cond(strings.Split(field, ":")...))
			first = false
		} else {
			q.SQL = fmt.Sprintf("%s AND %s", q.SQL, cond(strings.Split(field, ":")...))
		}
	}
	return q
}

// Join adds a basic JOIN clause using fields to define pairs of matching join fields
func (q *Query) Join(table string, fields ...[]string) *Query {
	// Generate field pairs
	var j []string
	for _, v := range fields {
		j = append(j, fmt.Sprintf("%s=%s", v[0], v[1]))
	}

	q.SQL = fmt.Sprintf("%s JOIN %s ON %s", q.SQL, table, strings.Join(j, ","))
	return q
}

// Order appends an ORDER BY statement to a query
func (q *Query) Order(fields ...string) *Query {
	q.SQL = fmt.Sprintf("%s ORDER BY %s", q.SQL, strings.Join(fields, ","))
	return q
}
