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

	for _, str := range strs {
		// ここからが実装

		id := &ImportDetail{
			PackageType: Unknown,
		}

		if str == "import (" || str == ")" {
			ids = append(ids, &ImportDetail{
				Alias:       NoAlias,
				ImportStr:   str,
				PackageType: Unknown,
			})
			continue
		}
		if str == "" {
			continue
		}

		replacedStr := strings.Replace(str, "\"", "", -1)
		//fmt.Printf("%s\n", replacedStr)

		splitedStr := strings.Split(replacedStr, "\t")[1]
		splitedStrs := strings.Split(splitedStr, " ")

		if len(splitedStrs) == 2 {
			id.Alias = splitedStrs[0]

			str := splitedStrs[1]
			id.ImportStr = strings.Replace(str, " ", "", -1)

			isStandard, _ := isStandardPackage(str)
			if isStandard {
				id.PackageType = Standard
			}

			if id.PackageType == Unknown {
				isOwnProject, _ := isOwnProjectPackage(str)
				if isOwnProject {
					id.PackageType = OwnProject
				}
			}

			if id.PackageType == Unknown {
				// StandardでもOwnProjectでもなければThirdPartyとする
				id.PackageType = ThirdParty
			}
		} else {
			id.Alias = "<no alias>"
			id.ImportStr = splitedStrs[0]

			isStandard, _ := isStandardPackage(splitedStrs[0])
			if isStandard {
				id.PackageType = Standard
			}

			if id.PackageType == Unknown {
				isOwnProject, _ := isOwnProjectPackage(splitedStrs[0])
				if isOwnProject {
					id.PackageType = OwnProject
				}
			}

			if id.PackageType == Unknown {
				// StandardでもOwnProjectでもなければThirdPartyとする
				id.PackageType = ThirdParty
			}
		}

		fmt.Printf("\t{\n\t\tImportStr:   %s,\n\t\tAlias:       %s, \n\t\tPackageType: %s,\n\t},\n", id.ImportStr, id.Alias, id.PackageType)
	}

	return nil, nil
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
