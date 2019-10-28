package lexer

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/istsh/goimport-fmt/ast"
	"github.com/istsh/goimport-fmt/config"
)

func TestLexer_OnlyOneImport(t *testing.T) {
	input := []byte(`package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}
`)

	want := &DividedSrc{
		BeforeImportDivision: input,
		ImportDivision:       nil,
		AfterImportDivision:  nil,
	}

	got := Lexer(input)
	if !reflect.DeepEqual(got.BeforeImportDivision, want.BeforeImportDivision) {
		t.Errorf("Lexer BeforeImportDivision got: %s; want: %s", got.BeforeImportDivision, want.BeforeImportDivision)
	}
	if !reflect.DeepEqual(got.ImportDivision, want.ImportDivision) {
		t.Errorf("Lexer ImportDivision got: %s; want: %s", got.ImportDivision, want.ImportDivision)
	}
	if !reflect.DeepEqual(got.AfterImportDivision, want.AfterImportDivision) {
		t.Errorf("Lexer AfterImportDivision got: %s; want: %s", got.AfterImportDivision, want.AfterImportDivision)
	}
}

func TestLexer_SomeImport(t *testing.T) {
	input := []byte(`package main

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/istsh/goimport-fmt/ast"
	c "github.com/istsh/goimport-fmt/config"
)

func main() {
	strs := strings.Split("Hello World", " ")
	fmt.Println(strs[0])
	fmt.Println(strs[1])
}
`)

	beforeBlock := []byte(`package main

import (
`)

	importBlock := [][]byte{
		[]byte("\t\"fmt\"\n"),
		[]byte("\t\"strings\"\n"),
		[]byte("\t\"github.com/jinzhu/gorm\"\n"),
		[]byte("\t\"github.com/istsh/goimport-fmt/ast\"\n"),
		[]byte("\tc \"github.com/istsh/goimport-fmt/config\"\n"),
	}

	afterBlock := []byte(`)

func main() {
	strs := strings.Split("Hello World", " ")
	fmt.Println(strs[0])
	fmt.Println(strs[1])
}
`)

	want := &DividedSrc{
		BeforeImportDivision: beforeBlock,
		ImportDivision:       importBlock,
		AfterImportDivision:  afterBlock,
	}

	config.Set(os.Getenv("GOROOT"), "", "github.com/istsh/goimport-fmt")
	got := Lexer(input)
	if !reflect.DeepEqual(got.BeforeImportDivision, want.BeforeImportDivision) {
		t.Errorf("Lexer BeforeImportDivision got: %s; want: %s", got.BeforeImportDivision, want.BeforeImportDivision)
	}
	if !reflect.DeepEqual(got.ImportDivision, want.ImportDivision) {
		t.Errorf("Lexer ImportDivision got: %s; want: %s", got.ImportDivision, want.ImportDivision)
	}
	if !reflect.DeepEqual(got.AfterImportDivision, want.AfterImportDivision) {
		t.Errorf("Lexer AfterImportDivision got: %s; want: %s", got.AfterImportDivision, want.AfterImportDivision)
	}
}

func TestGetImportDetails(t *testing.T) {
	input := []byte(`package main

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/istsh/goimport-fmt/ast"
	c "github.com/istsh/goimport-fmt/config"
)

func main() {
	strs := strings.Split("Hello World", " ")
	fmt.Println(strs[0])
	fmt.Println(strs[1])
}
`)

	want := ast.ImportDetails{
		&ast.ImportDetail{
			Alias:       []byte(""),
			ImportPath:  []byte("fmt"),
			PackageType: ast.Standard,
		},
		&ast.ImportDetail{
			Alias:       []byte(""),
			ImportPath:  []byte("strings"),
			PackageType: ast.Standard,
		},
		&ast.ImportDetail{
			Alias:       []byte(""),
			ImportPath:  []byte("github.com/jinzhu/gorm"),
			PackageType: ast.ThirdParty,
		},
		&ast.ImportDetail{
			Alias:       []byte(""),
			ImportPath:  []byte("github.com/istsh/goimport-fmt/ast"),
			PackageType: ast.OwnProject,
		},
		&ast.ImportDetail{
			Alias:       []byte("c"),
			ImportPath:  []byte("github.com/istsh/goimport-fmt/config"),
			PackageType: ast.OwnProject,
		},
	}

	config.Set(os.Getenv("GOROOT"), "", "github.com/istsh/goimport-fmt")
	ds := Lexer(input)
	got, err := ds.GetImportDetails()
	if err != nil {
		t.Fatalf("GetImportDetails()=_, %#v; want nil", err)
	}
	for i := 0; i < len(got); i++ {
		g := got[i]
		w := want[i]
		if !bytes.Equal(g.Alias, w.Alias) {
			t.Errorf("GetImportDetails Alias got: %s, want: %s\n", g.Alias, w.Alias)
		}
		if !bytes.Equal(g.ImportPath, w.ImportPath) {
			t.Errorf("GetImportDetails ImportPath got: %s, want: %s\n", g.ImportPath, w.ImportPath)
		}
		if g.PackageType != w.PackageType {
			t.Errorf("GetImportDetails PackageType got: %d, want: %d\n", g.PackageType, w.PackageType)
		}
	}
}
