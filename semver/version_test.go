package semver

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name    string
		version *Version
		want    string
	}{
		{
			"just major",
			&Version{versions: []int{1}},
			"1",
		},
		{
			"major and minor",
			&Version{versions: []int{1, 2}},
			"1.2",
		},
		{
			"major, minor and patch",
			&Version{versions: []int{1, 2, 3}},
			"1.2.3",
		},
		{
			"to infinity and beyond",
			&Version{versions: []int{1, 2, 3, 4, 5}},
			"1.2.3.4.5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.version.String(); got != tt.want {
				t.Errorf("Version.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    *Version
		err     error
	}{
		{
			"just major",
			"10",
			&Version{versions: []int{10}},
			nil,
		},
		{
			"major and minor",
			"1.2",
			&Version{versions: []int{1, 2}},
			nil,
		},
		{
			"major, minor and patch",
			"1.2.3",
			&Version{versions: []int{1, 2, 3}},
			nil,
		},
		{
			"to infinity and beyond",
			"1.2.3.4.5",
			&Version{versions: []int{1, 2, 3, 4, 5}},
			nil,
		},
		{
			"something else with the version",
			"before1.2.3after",
			nil,
			errors.New("version contains extra stuff"),
		},
		{
			"errors when invalid format",
			"foo",
			nil,
			errors.New("could not match \"foo\" against version regexp"),
		},
		{
			"giant number",
			"1.22222222222222222222222222222222.3",
			nil,
			errors.New("could not convert \"22222222222222222222222222222222\" to number"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.version)
			if tt.err != nil {
				if err == nil {
					t.Error("want error, got nil")
					return
				}

				if err.Error() != tt.err.Error() {
					t.Errorf("New() error = %v, wantErr %v", err, tt.err)
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		v1   *Version
		v2   *Version
		want int
	}{
		{MustNew("1"), MustNew("1"), 0},
		{MustNew("1"), MustNew("2"), -1},
		{MustNew("2"), MustNew("1"), 1},

		{MustNew("1.1"), MustNew("1.1"), 0},
		{MustNew("1.1"), MustNew("1.2"), -1},
		{MustNew("1.2"), MustNew("1.1"), 1},

		{MustNew("1.1.3"), MustNew("1.1.3"), 0},
		{MustNew("1.1.2"), MustNew("1.1.3"), -1},
		{MustNew("1.1.3"), MustNew("1.1.2"), 1},

		{MustNew("1.1.1.3"), MustNew("1.1.1.3"), 0},
		{MustNew("1.1.1.2"), MustNew("1.1.1.3"), -1},
		{MustNew("1.1.1.3"), MustNew("1.1.1.2"), 1},

		{MustNew("1.0"), MustNew("1"), 1},
		{MustNew("1"), MustNew("1.0"), -1},

		{MustNew("1.2.3.4.5"), MustNew("1"), 1},
		{MustNew("1"), MustNew("1.2.3.4.5"), -1},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s <> %s = %d", tt.v1, tt.v2, tt.want), func(t *testing.T) {
			if got := tt.v1.Compare(tt.v2); got != tt.want {
				t.Errorf("Version.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
