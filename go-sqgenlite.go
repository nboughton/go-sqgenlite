// Package sqgenlite is a really dumb sqlite query generator
package sqgenlite

import (
	"fmt"
	"strings"
)

// Query struct that we attach our generate methods to
type Query struct {
	SQL  string
	Args []interface{}
}

// NewQuery creates a new query object
func NewQuery() *Query {
	return new(Query)
}

// Insert generates a basic insert statement
func (q *Query) Insert(table string, fields []string, args ...interface{}) *Query {
	// Generate placeholders
	var p []string
	for range fields {
		p = append(p, "?")
	}

	q.SQL = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(fields, ","), strings.Join(p, ","))
	q.Args = append(q.Args, args...)
	return q
}

// Update generates the first part of an UPDATE statement. a Where clause will be necessary
// to complete the SQL query
func (q *Query) Update(table string, fields []string, args ...interface{}) *Query {
	q.SQL = fmt.Sprintf("UPDATE %s SET %s=?", table, strings.Join(fields, "=?, "))
	q.Args = append(q.Args, args...)
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
func (q *Query) Where(s string, args ...interface{}) *Query {
	q.SQL = fmt.Sprintf("%s WHERE %s", q.SQL, s)
	q.Args = append(q.Args, args...)
	return q
}

// Join adds a basic JOIN clause using fields to define pairs of matching join fields, this is fine in sqlite as it only
// recognises inner joins
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
