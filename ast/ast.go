package ast

import (
	"bytes"
	"fmt"
	"os"

	"github.com/istsh/goimport-fmt/config"
)

type PackageType uint

const (
	Unknown PackageType = iota
	Standard
	ThirdParty
	OwnProject
)

type ImportDetail struct {
	Alias       []byte
	ImportPath  []byte
	PackageType PackageType
}

func AnalyzeIncludeAlias(alias, ImportPath []byte) (*ImportDetail, error) {
	id, err := Analyze(ImportPath)
	if err != nil {
		return nil, err
	}

	id.Alias = alias
	return id, nil
}

func Analyze(importPath []byte) (*ImportDetail, error) {
	packageType := Unknown

	isStandard, err := isStandardPackage(importPath)
	if err != nil {
		return nil, err
	}
	if isStandard {
		packageType = Standard
	}

	if packageType == Unknown {
		isOwnProject := isOwnProjectPackage(importPath)
		if isOwnProject {
			packageType = OwnProject
		}
	}

	if packageType == Unknown {
		// StandardでもOwnProjectでもなければThirdPartyとする
		packageType = ThirdParty
	}

	return &ImportDetail{
		ImportPath:  importPath,
		PackageType: packageType,
	}, nil
}

func isStandardPackage(importPath []byte) (bool, error) {
	p := fmt.Sprintf("%s/src/%s", config.GetEnv().GetGoRoot(), importPath)

	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func isOwnProjectPackage(importPath []byte) bool {
	return bytes.HasPrefix(importPath, []byte(config.GetEnv().GetOwnProject()))
}

type ImportDetails []*ImportDetail

func (ids ImportDetails) Len() int {
	return len(ids)
}

func (ids ImportDetails) Less(i, j int) bool {
	if ids[i].PackageType < ids[j].PackageType {
		return true
	}
	if ids[i].PackageType > ids[j].PackageType {
		return false
	}
	return bytes.Compare(ids[i].ImportPath, ids[j].ImportPath) < 0
}

func (ids ImportDetails) Swap(i, j int) {
	ids[i], ids[j] = ids[j], ids[i]
}
