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
