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

// Get generates a SELECT that can then be chained with further functions
func (q *Query) Get(fields ...string) *Query {
	q.SQL = fmt.Sprintf("SELECT %s", strings.Join(fields, ", "))
	return q
}

// From adds a FROM clause
func (q *Query) From(table string) *Query {
	q.SQL = fmt.Sprintf("%s FROM %s", q.SQL, table)
	return q
}

// Where adds len(fields) WHERE/AND clauses
func (q *Query) Where(fields ...string) *Query {
	q.SQL = fmt.Sprintf("%s WHERE %s=?", q.SQL, fields[0])
	if len(fields) > 1 {
		for _, v := range fields[1:] {
			q.SQL = fmt.Sprintf("%s AND %s=?", q.SQL, v)
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
