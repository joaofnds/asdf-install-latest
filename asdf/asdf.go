package asdf

import (
	"os/exec"

	"github.com/coreos/go-semver/semver"
)

type ListResult = map[string]semver.Versions

func Reshim() error {
	return exec.Command("asdf", "reshim").Run()
}

func PluginList() ([]string, error) {
	b, err := exec.Command("asdf", "plugin", "list").Output()
	if err != nil {
		return nil, err
	}
	return parseWordList(b), nil
}

func List() (ListResult, error) {
	b, err := exec.Command("asdf", "list").Output()
	if err != nil {
		return nil, err
	}
	return parseResult(b), nil
}

func ListPkg(pkg string) (semver.Versions, error) {
	b, err := exec.Command("asdf", "list", pkg).Output()
	if err != nil {
		return nil, err
	}

	return parseVersionList(b), nil
}

func ListAllPkg(pkg string) (semver.Versions, error) {
	b, err := exec.Command("asdf", "list", "all", pkg).Output()
	if err != nil {
		return nil, err
	}

	return parseVersionList(b), nil
}

func Latest(pkg string) (*semver.Version, error) {
	vs, err := ListAllPkg(pkg)
	if err != nil {
		return nil, err
	}

	if len(vs) == 0 {
		return nil, NewNoVersionsErr(pkg)
	}

	return vs[len(vs)-1], nil
}

func Install(pkg string, v *semver.Version) error {
	return exec.Command("asdf", "install", pkg, v.String()).Run()
}

func SetGlobal(pkg string, v *semver.Version) error {
	return exec.Command("asdf", "global", pkg, v.String()).Run()
}

func Current(pkg string) (*semver.Version, error) {
	b, err := exec.Command("asdf", "current", pkg).Output()
	if err != nil {
		return nil, err
	}

	v := parseCurrent(b)

	if v == nil {
		return nil, NewNoVersionsErr(pkg)
	}

	return v, nil
}
