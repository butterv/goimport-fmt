package ast

import (
	"fmt"
	"os"
	"strings"

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

func Analyze(strs []string) ([]*ImportDetail, error) {
	var ids []*ImportDetail

	if len(strs) == 0 {
		return ids, nil
	}

	if len(strs) == 1 {
		str := strs[0]
		trimStr := strings.Trim(str, "import ")
		replaceStr := strings.Replace(trimStr, "\"", "", -1)
		splitStrs := strings.Split(replaceStr, " ")
		return ids, nil
	}

	for _, str := range strs {
		if str == "" {
			continue
		}

		if str == "import (" || str == ")" {
			ids = append(ids, &ImportDetail{
				Alias:       NoAlias,
				ImportStr:   str,
				PackageType: Unknown,
			})
			continue
		}

		trimStr := strings.Trim(str, "\t")
		replaceStr := strings.Replace(trimStr, "\"", "", -1)
		splitStrs := strings.Split(replaceStr, " ")

		if len(splitStrs) == 2 {
			id, err := analyzeIncludeAlias(splitStrs[0], splitStrs[1])
			if err != nil {
				// TODO: error handling
			}
			ids = append(ids, id)
		} else {
			importStr := splitStrs[0]

			id := &ImportDetail{
				Alias:       NoAlias,
				ImportStr:   importStr,
				PackageType: Unknown,
			}

			isStandard, _ := isStandardPackage(importStr)
			if isStandard {
				id.PackageType = Standard
			}

			if id.PackageType == Unknown {
				isOwnProject, _ := isOwnProjectPackage(importStr)
				if isOwnProject {
					id.PackageType = OwnProject
				}
			}

			if id.PackageType == Unknown {
				// StandardでもOwnProjectでもなければThirdPartyとする
				id.PackageType = ThirdParty
			}

			ids = append(ids, id)
		}
	}

	return ids, nil
}

func analyzeIncludeAlias(alias, importStr string) (*ImportDetail, error) {
	id, err := analyze(importStr)
	if err != nil {
		return nil, err
	}

	id.Alias = alias
	return id, nil
}

func analyze(importStr string) (*ImportDetail, error) {
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

	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func isOwnProjectPackage(path string) (bool, error) {
	p := fmt.Sprintf("%s/src/%s", config.GetEnv().GetGoPath(), path)

	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
