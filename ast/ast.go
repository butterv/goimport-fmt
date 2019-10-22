package ast

import (
	"fmt"
	"os"

	"github.com/istsh/goimport-fmt/config"
)

type PackageType string

const (
	Unknown    PackageType = "unknown"
	Standard               = "standard"
	ThirdParty             = "third party"
	OwnProject             = "own project"
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
		isOwnProject, _ := isOwnProjectPackage(importStr)
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
		PackageType: Unknown,
	}, nil
}

func isStandardPackage(path string) (bool, error) {
	p := fmt.Sprintf("%s/src/%s", config.GetEnv().GetGoRoot(), path)
	fmt.Printf("p: %s\n", p)

	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func isOwnProjectPackage(path string) (bool, error) {
	// TODO: ここをどう判定するか。できればmodulesの設定によって処理を分けたくない。

	p := fmt.Sprintf("%s/src/%s", config.GetEnv().GetGoPath(), path)

	fmt.Printf("p: %s\n", p)

	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
