package asdf

type NoVersions string

func NewNoVersionsErr(pkg string) error {
	return NoVersions(pkg)
}

func (pkg NoVersions) Error() string {
	return "no versions found for package '" + string(pkg) + "'"
}
