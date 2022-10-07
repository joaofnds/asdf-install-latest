package semver

import (
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
	versions := []*Version{
		MustNew("1.1.4"),
		MustNew("1.1.3"),
		MustNew("1.1.2"),
		MustNew("1.1.1"),
	}
	Sort(versions)

	expected := []*Version{
		MustNew("1.1.1"),
		MustNew("1.1.2"),
		MustNew("1.1.3"),
		MustNew("1.1.4"),
	}

	if !reflect.DeepEqual(versions, expected) {
		t.Errorf("expected %v, got %v", expected, versions)
	}
}
