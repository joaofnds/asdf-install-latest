package semver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var reg = regexp.MustCompile(`\d+(.\d+)*`)

type Version struct {
	versions []int
}

func (v *Version) Len() int {
	return len(v.versions)
}

func (v *Version) LessThan(other *Version) bool {
	return v.Compare(other) < 0
}

func (v *Version) Compare(other *Version) int {
	vLen := v.Len()
	otherLen := other.Len()
	for i := 0; i < max(vLen, otherLen); i++ {
		switch {
		case i >= vLen:
			return -1
		case i >= otherLen:
			return 1
		case v.versions[i] > other.versions[i]:
			return 1
		case v.versions[i] < other.versions[i]:
			return -1
		}
	}
	return 0
}

func (v *Version) String() string {
	versions := make([]string, len(v.versions))
	for i, n := range v.versions {
		versions[i] = strconv.Itoa(n)
	}
	return strings.Join(versions, ".")
}

func MustNew(s string) *Version {
	v, err := New(s)
	if err != nil {
		panic(err)
	}
	return v
}

func New(v string) (*Version, error) {
	match := reg.FindString(v)
	if match == "" {
		return nil, fmt.Errorf("could not match %q against version regexp", v)
	}

	if match != v {
		return nil, errors.New("version contains extra stuff")
	}

	versions := strings.Split(match, ".")
	version := &Version{versions: make([]int, len(versions))}
	for i, v := range versions {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("could not convert %q to number", v)
		}

		version.versions[i] = n
	}

	return version, nil
}
