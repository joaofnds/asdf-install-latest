package asdf

import (
	"bufio"
	"bytes"

	"github.com/coreos/go-semver/semver"
)

func parseResult(result []byte) ListResult {
	out := ListResult{}
	lang := ""

	s := bufio.NewScanner(bytes.NewReader(result))
	for s.Scan() {
		t := s.Text()
		if t[0] == ' ' {
			if v, err := semver.NewVersion(t[2:]); err == nil {
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
		if v, err := semver.NewVersion(s.Text()); err == nil {
			out = append(out, v)
		}
	}

	semver.Sort(out)
	return out
}

func parseCurrent(b []byte) *semver.Version {
	s := bufio.NewScanner(bytes.NewReader(b))
	s.Split(bufio.ScanWords)
	for s.Scan() {
		if v, err := semver.NewVersion(s.Text()); err == nil {
			return v
		}
	}

	return nil
}
