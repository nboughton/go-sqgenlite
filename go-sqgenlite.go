// Package sqgenlite is a really dumb sqlite query generator
package sqgenlite

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

// Filter is utilised by FilterSet so that an ordered Array of conditions can be
// matched to appropriate values in db.Query/Exec etc
type Filter struct {
	Op    Conditional
	Field string
}

// FilterSet is the type used for mapping field sets to conditionals
type FilterSet []Filter

// NewFilterSet returns a FilterSet pointer as a more intuitive way of creating filters
func NewFilterSet() *FilterSet {
	return new(FilterSet)
}

// Add adds a condition to the FilterSet as a shorthand to improve readability, returns the
// FilterSet so it can be chained. For operations that use multiple fields use ':' as a separator
// e.g for a BETWEEN statement using a DATE operator you would use c.Add(Between, "date:DATE")
func (c *FilterSet) Add(field string, op Conditional) *FilterSet {
	*c = append(*c, Filter{Op: op, Field: field})
	return c
}

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

// Between expects a field and potentially a function i.e DATE(), SUM() etc
// requires the comparison values to be passed in during execution. Returns a BETWEEN clause
func Between(fields ...string) string {
	if len(fields) > 1 {
		return fmt.Sprintf("%s(%s) BETWEEN %s(?) AND %s(?)", fields[1], fields[0], fields[1], fields[1])
	}
	return fmt.Sprintf("%s BETWEEN ? AND ?", fields[0])
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
	q.SQL = fmt.Sprintf("UPDATE %s SET %s=?", table, strings.Join(fields, "=?, "))
	return q
}

// Select generates a SELECT that can then be chained with further functions
func (q *Query) Select(fields ...string) *Query {
	q.SQL = fmt.Sprintf("SELECT %s", strings.Join(fields, ", "))
	return q
}

// Delete begins a Delete query.
func (q *Query) Delete() *Query {
	q.SQL = "DELETE"
	return q
}

// From adds a FROM clause
func (q *Query) From(table string) *Query {
	q.SQL = fmt.Sprintf("%s FROM %s", q.SQL, table)
	return q
}

// Where adds len(fields) WHERE fields is a filter set created by
// NewFilterSet and added to with the .Add function.
func (q *Query) Where(fields *FilterSet) *Query {
	for i, c := range *fields {
		if i == 0 {
			q.SQL = fmt.Sprintf("%s WHERE %s", q.SQL, c.Op(strings.Split(c.Field, ":")...))
		} else {
			q.SQL = fmt.Sprintf("%s AND %s", q.SQL, c.Op(strings.Split(c.Field, ":")...))
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

// Group appends an ORDER BY statement to a query
func (q *Query) Group(fields ...string) *Query {
	q.SQL = fmt.Sprintf("%s GROUP BY %s", q.SQL, strings.Join(fields, ","))
	return q
}

// Append is the cop-out method to just string stuff together.
func (q *Query) Append(s string) *Query {
	q.SQL = fmt.Sprintf("%s %s", q.SQL, s)
	return q
}
