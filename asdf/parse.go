package asdf

import (
	"ail/semver"
	"bufio"
	"bytes"
	"strings"
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
	var out []string

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
		if v := parseSemVer(removeCurrentMarkFromVersion(s.Text())); v != nil {
			out = append(out, v)
		}
	}

	semver.Sort(out)
	return out
}

// removeCurrentMarkFromVersion removes the '*' char that asdf now uses on `asdf list`
// to indicate which of the installed versions is currently being used
func removeCurrentMarkFromVersion(s string) string {
	if s[0] == '*' {
		return s[1:]
	}
	return s
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
	v, err := semver.New(str)
	if err == nil {
		return v
	}

	if err != nil && strings.Contains(err.Error(), "not in dotted-tri format") {
		str += ".0"
		v, err = semver.New(str)
		if err == nil {
			return v
		}
	}

	return nil
}
