package lexer

import (
	"reflect"
	"strings"
	"testing"

	"github.com/istsh/goimport-fmt/ast"
)

func Test1(t *testing.T) {
	tests := []*ast.ImportDetail{
		{
			ast.NoAlias,
			"context",
			ast.Standard,
		},
		{
			ast.NoAlias,
			"math/rand",
			ast.Standard,
		},
		{
			ast.NoAlias,
			"net/http",
			ast.Standard,
		},
		{
			ast.NoAlias,
			"time",
			ast.Standard,
		},
		{
			ast.NoAlias,
			"github.com/go-sql-driver/mysql",
			ast.ThirdParty,
		},
		{
			ast.NoAlias,
			"github.com/jinzhu/gorm",
			ast.ThirdParty,
		},
		{
			ast.NoAlias,
			"github.com/labstack/echo",
			ast.ThirdParty,
		},
		{
			"sqltrace",
			`sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"`,
			ast.ThirdParty,
		},
		{
			"gormtrace",
			`gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jinzhu/gorm"`,
			ast.ThirdParty,
		},
		{
			ast.NoAlias,
			"github.com/istsh/own-project/app/impl/userserviceimpl",
			ast.OwnProject,
		},
		{
			ast.NoAlias,
			"github.com/istsh/own-project/app/impl/todoserviceimpl",
			ast.OwnProject,
		},
		{
			ast.NoAlias,
			"github.com/istsh/own-project/app/impl/repository/database",
			ast.OwnProject,
		},
		{
			"imiddleware",
			`imiddleware "github.com/istsh/own-project/app/middleware"`,
			ast.OwnProject,
		},
	}

	//config.Setup("github.com/istsh/own-project")

	var paths []string
	for _, tt := range tests {
		paths = append(paths, tt.ImportStr)
	}

	ids, _ := Lexer(nil)
	for i, tt := range tests {
		id := ids[i]

		if tt.Alias != ast.NoAlias {
			replaceStr := strings.Replace(tt.ImportStr, "\"", "", -1)
			splitStrs := strings.Split(replaceStr, " ")
			want := ast.ImportDetail{
				Alias:       tt.Alias,
				ImportStr:   splitStrs[1],
				PackageType: tt.PackageType,
			}
			if !reflect.DeepEqual(*id, want) {
				t.Errorf("Lexer(paths)\nindex:%d\ngot:  %#v\nwant: %#v", i, *id, want)
			}
		} else if !reflect.DeepEqual(*id, *tt) {
			t.Errorf("Lexer(paths)\nindex:%d\ngot:  %#v\nwant: %#v", i, *id, *tt)
		}
	}
}
