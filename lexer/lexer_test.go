package lexer

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/istsh/goimport-fmt/config"
)

type PackageType string

const (
	Unknown    PackageType = "unknown"
	Standard               = "standard"
	ThirdParty             = "third party"
	OwnProject             = "own project"
)

type ImportDetail struct {
	Alias       string
	ImportStr   string
	PackageType PackageType
}

//type Environment struct {
//	GOROOT      string
//	GOPATH      string
//	GO111MODULE string
//}
//
//var Env *Environment
//
//func setup() {
//	goroot := os.Getenv("GOROOT")
//	gopath := os.Getenv("GOPATH")
//	go111module := os.Getenv("GO111MODULE")
//
//	Env = &Environment{
//		GOROOT:      goroot,
//		GOPATH:      gopath,
//		GO111MODULE: go111module,
//	}
//}

func Test1(t *testing.T) {
	input := `import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jinzhu/gorm"

	"github.com/istsh/own-project/app/impl/userserviceimpl"
	"github.com/istsh/own-project/app/impl/todoserviceimpl"
	"github.com/istsh/own-project/app/impl/repository/database"
	imiddleware "github.com/istsh/own-project/app/middleware"
)`

	//tests := []ImportDetail{
	//	{
	//		ImportStr:   "context",
	//		Alias:       "<no alias>",
	//		PackageType: Standard,
	//	},
	//	{
	//		ImportStr:   "math/rand",
	//		Alias:       "<no alias>",
	//		PackageType: Standard,
	//	},
	//	{
	//		ImportStr:   "net/http",
	//		Alias:       "<no alias>",
	//		PackageType: Standard,
	//	},
	//	{
	//		ImportStr:   "time",
	//		Alias:       "<no alias>",
	//		PackageType: Standard,
	//	},
	//	{
	//		ImportStr:   "github.com/go-sql-driver/mysql",
	//		Alias:       "<no alias>",
	//		PackageType: ThirdParty,
	//	},
	//	{
	//		ImportStr:   "github.com/jinzhu/gorm",
	//		Alias:       "<no alias>",
	//		PackageType: ThirdParty,
	//	},
	//	{
	//		ImportStr:   "github.com/labstack/echo",
	//		Alias:       "<no alias>",
	//		PackageType: ThirdParty,
	//	},
	//	{
	//		ImportStr:   "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql",
	//		Alias:       "sqltrace",
	//		PackageType: ThirdParty,
	//	},
	//	{
	//		ImportStr:   "gopkg.in/DataDog/dd-trace-go.v1/contrib/jinzhu/gorm",
	//		Alias:       "gormtrace",
	//		PackageType: ThirdParty,
	//	},
	//	{
	//		ImportStr:   "github.com/istsh/own-project/app/impl/userserviceimpl",
	//		Alias:       "<no alias>",
	//		PackageType: OwnProject,
	//	},
	//	{
	//		ImportStr:   "github.com/istsh/own-project/app/impl/todoserviceimpl",
	//		Alias:       "<no alias>",
	//		PackageType: OwnProject,
	//	},
	//	{
	//		ImportStr:   "github.com/istsh/own-project/app/impl/repository/database",
	//		Alias:       "<no alias>",
	//		PackageType: OwnProject,
	//	},
	//	{
	//		ImportStr:   "github.com/istsh/own-project/app/middleware",
	//		Alias:       "imiddleware",
	//		PackageType: OwnProject,
	//	},
	//}

	// l := New(input)

	// replacedStr := strings.Replace(input, "\"", "", -1)
	// fmt.Printf(replacedStr)

	// setup()

	var importStrs []string
	for _, str := range strings.Split(input, "\n") {
		// fmt.Printf("L%d: %s\n", i, str)
		importStrs = append(importStrs, str)
	}

	if len(importStrs) == 1 {
		return
	}

	fmt.Println("{")
	for _, importStr := range importStrs {
		// ここからが実装

		id := &ImportDetail{
			PackageType: Unknown,
		}

		if importStr == "import (" || importStr == ")" {
			id.Alias = "<no alias>"
			id.ImportStr = importStr
			continue
		}
		if importStr == "" {
			continue
		}

		replacedStr := strings.Replace(importStr, "\"", "", -1)
		fmt.Printf("%s\n", replacedStr)

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
	fmt.Println("}")
}

func isStandardPackage(path string) (bool, error) {
	p := fmt.Sprintf("%s/src/%s", config.GetEnv().GetGoRoot(), path)
	//fmt.Printf("path: %s\n", p)

	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func isOwnProjectPackage(path string) (bool, error) {
	p := fmt.Sprintf("%s/src/%s", config.GetEnv().GetGoPath(), path)
	//fmt.Printf("path: %s\n", p)

	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
