package sqgenlite

import (
	"reflect"
	"testing"
)

var (
	table  = "table"
	fields = []string{"field1", "field2", "field3"}
	args   = []interface{}{"arg1", 2, float64(3)}
)

func TestNewQuery(t *testing.T) {
	q := NewQuery()
	if !reflect.DeepEqual(q, &Query{}) {
		t.Fatalf("NewQuery does not return a pointer to a Query object")
	}
}

func TestInsert(t *testing.T) {
	q := NewQuery().Insert(table, fields, args...)
	if q.SQL != "INSERT INTO table (field1, field2, field3) VALUES (?, ?, ?)" {
		t.Fatalf("Insert failed to create SQL successfully [%s]\n", q.SQL)
	}
	if !reflect.DeepEqual(q.Args, args) {
		t.Fatalf("Insert did not copy args to q.Args\n")
	}
}

func TestUpdate(t *testing.T) {
	q := NewQuery().Update(table, fields, args...)
	if q.SQL != "UPDATE table SET field1=?, field2=?, field3=?" {
		t.Fatalf("Update failed to create SQL successfully [%s]\n", q.SQL)
	}
	if !reflect.DeepEqual(q.Args, args) {
		t.Fatalf("Update did not copy args to q.Args\n")
	}
}

func TestSelect(t *testing.T) {
	q := NewQuery().Select(table, fields...)
	if q.SQL != "SELECT field1, field2, field3 FROM table" {
		t.Fatalf("Select failed to create SQL successfully [%s]\n", q.SQL)
	}
}

func TestWhere(t *testing.T) {
	q := NewQuery().Where("field1 = ?", args[:1])
	if q.SQL != " WHERE field1 = ?" {
		t.Fatalf("Where failed to create SQL successfully [%s]\n", q.SQL)
	}
	if len(q.Args) != 1 {
		t.Fatalf("Where did not update q.Args properly [%d]\n", len(q.Args))
	}
}

func TestDelete(t *testing.T) {
	q := NewQuery().Delete(table)
	if q.SQL != "DELETE FROM table" {
		t.Fatalf("Delete failed to create SQL successfully [%s]\n", q.SQL)
	}
}

func TestJoin(t *testing.T) {
	q := NewQuery().Join(table, []string{"t1.field1", "t2.field1"})
	if q.SQL != " JOIN table ON t1.field1=t2.field1" {
		t.Fatalf("Join failed to create SQL successfully [%s]\n", q.SQL)
	}
}

func TestOrder(t *testing.T) {
	q := NewQuery().Order(fields...)
	if q.SQL != " ORDER BY field1, field2, field3" {
		t.Fatalf("Order failed to create SQL successfully [%s]\n", q.SQL)
	}
}
