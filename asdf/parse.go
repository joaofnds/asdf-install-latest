package asdf

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/coreos/go-semver/semver"
)

func parseResult(result []byte) ListResult {
	out := ListResult{}
	lang := ""

	s := bufio.NewScanner(bytes.NewReader(result))
	for s.Scan() {
		t := s.Text()
		if t[0] == ' ' {
			if v := parseSemVer(t[2:]); v != nil {
				out[lang] = append(out[lang], v)
			}
		} else {
			lang = t
		}
	}

	return out
}

func parseWordList(b []byte) []string {
	out := []string{}

	s := bufio.NewScanner(bytes.NewReader(b))
	s.Split(bufio.ScanWords)
	for s.Scan() {
		out = append(out, s.Text())
	}

	return out
}

func parseVersionList(b []byte) semver.Versions {
	out := semver.Versions{}

	s := bufio.NewScanner(bytes.NewReader(b))
	s.Split(bufio.ScanWords)
	for s.Scan() {
		if v := parseSemVer(s.Text()); v != nil {
			out = append(out, v)
		}
	}

	semver.Sort(out)
	return out
}

// parses output of `asdf current {plugin}`
// expected input: golang          1.19.2          /Users/joaofnds/.tool-versions
func parseCurrent(b []byte) *semver.Version {
	s := bufio.NewScanner(bytes.NewReader(b))
	s.Split(bufio.ScanWords)
	s.Scan() // skip plugin name
	s.Scan()
	return parseSemVer(s.Text())
}

func parseSemVer(str string) *semver.Version {
	v, err := semver.NewVersion(str)
	if err == nil && v.PreRelease == "" {
		return v
	}

	if err != nil && strings.Contains(err.Error(), "not in dotted-tri format") {
		str += ".0"
		v, err = semver.NewVersion(str)
		if err == nil && v.PreRelease == "" {
			return v
		}
	}

	return nil
}
