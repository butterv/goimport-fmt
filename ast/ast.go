package ast

import (
	"fmt"
	"os"
	"strings"

	"github.com/istsh/goimport-fmt/config"
)

type PackageType uint

const (
	Unknown PackageType = iota
	Standard
	ThirdParty
	OwnProject
)

const NoAlias = "<no alias>"

type ImportDetail struct {
	Alias       string
	ImportStr   string
	PackageType PackageType
}

func AnalyzeIncludeAlias(alias, importStr string) (*ImportDetail, error) {
	id, err := Analyze(importStr)
	if err != nil {
		return nil, err
	}

	id.Alias = alias
	return id, nil
}

func Analyze(importStr string) (*ImportDetail, error) {
	packageType := Unknown

	isStandard, err := isStandardPackage(importStr)
	if err != nil {
		return nil, err
	}
	if isStandard {
		packageType = Standard
	}

	if packageType == Unknown {
		isOwnProject := isOwnProjectPackage(importStr)
		if isOwnProject {
			packageType = OwnProject
		}
	}

	if packageType == Unknown {
		// StandardでもOwnProjectでもなければThirdPartyとする
		packageType = ThirdParty
	}

	return &ImportDetail{
		Alias:       NoAlias,
		ImportStr:   importStr,
		PackageType: packageType,
	}, nil
}

func isStandardPackage(path string) (bool, error) {
	p := fmt.Sprintf("%s/src/%s", config.GetEnv().GetGoRoot(), path)

	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func isOwnProjectPackage(path string) bool {
	return strings.HasPrefix(path, config.GetEnv().GetOwnProject())
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
	return ids[i].ImportStr < ids[j].ImportStr
}

func (ids ImportDetails) Swap(i, j int) {
	ids[i], ids[j] = ids[j], ids[i]
}
