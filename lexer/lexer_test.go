package lexer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/istsh/goimport-fmt/ast"
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

	//config.Setup()

	var strs []string
	for _, str := range strings.Split(input, "\n") {
		strs = append(strs, str)
	}

	if len(strs) == 1 {
		return
	}

	ids, err := ast.Analyze(strs)
	if err != nil {
		// TODO: error handling
	}

	fmt.Println("{")
	for _, id := range ids {
		fmt.Printf("\t{\n\t\tImportStr:   %s,\n\t\tAlias:       %s, \n\t\tPackageType: %s,\n\t},\n", id.ImportStr, id.Alias, id.PackageType)
	}
	fmt.Println("}")
}
