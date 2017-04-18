package sqgenlite

import (
	"reflect"
	"testing"
)

func TestNewFilterSet(t *testing.T) {
	if !reflect.DeepEqual(NewFilterSet(), new(FilterSet)) {
		t.Fatal("Does not return expected value")
	}
}

func TestAdd(t *testing.T) {
	v := &FilterSet{
		filters: []Filter{
			Filter{
				Field: "field",
				Op:    Eq,
			},
		},
	}
	s := NewFilterSet().Add("field", Eq)
	if v.filters[0].Field != s.filters[0].Field {
		t.Fatal("Field set failed in Add")
	}

	if !reflect.DeepEqual(Conditional(v.filters[0].Op), Conditional(Eq)) {
		t.Fatal("Op set failed in Add")
	}
}
