package ast

import (
	"bytes"
	"fmt"
	"os"

	"github.com/istsh/goimport-fmt/config"
)

// PackageType is type of import package
type PackageType uint

const (
	// Unknown means an unknown definition
	Unknown PackageType = iota
	// Standard means a standard package
	Standard
	// ThirdParty means a third party package
	ThirdParty
	// OwnProject means a own project
	OwnProject
)

// ImportDetail ImportDetail has details of each import path
type ImportDetail struct {
	Alias       []byte
	ImportPath  []byte
	PackageType PackageType
}

// AnalyzeIncludeAlias evaluates import path and returns ImportDetail with alias
func AnalyzeIncludeAlias(alias, ImportPath []byte) (*ImportDetail, error) {
	id, err := Analyze(ImportPath)
	if err != nil {
		return nil, err
	}

	id.Alias = alias
	return id, nil
}

// Analyze evaluates import path and returns ImportDetail
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
		// if it is neither Standard nor OwnProject, then ThirdParty
		packageType = ThirdParty
	}

	return &ImportDetail{
		ImportPath:  importPath,
		PackageType: packageType,
	}, nil
}

func isStandardPackage(importPath []byte) (bool, error) {
	p := fmt.Sprintf("%s/src/%s", config.GetGoRoot(), importPath)

	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func isOwnProjectPackage(importPath []byte) bool {
	return bytes.HasPrefix(importPath, []byte(config.GetOwnProject()))
}

// ImportDetails has some ImportDetail
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
