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
	input := []byte(`
package main

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
	beforeBlock := []byte(`package main

import (
`)

	importBlock := [][]byte{
		[]byte("fmt\n"),
		[]byte("strings\n"),
		[]byte("\n"),
		[]byte("github.com/jinzhu/gorm\n"),
		[]byte("\n"),
		[]byte("github.com/istsh/goimport-fmt/ast\n"),
		[]byte("c \"github.com/istsh/goimport-fmt/config\"\n"),
	}

	afterBlock := []byte(`)

func main() {
	strs := strings.Split("Hello World", " ")
	fmt.Println(strs[0])
	fmt.Println(strs[1])
}
`)

	var input []byte
	var wantImport [][]byte
	input = append(input, beforeBlock...)
	for _, block := range importBlock {
		input = append(input, block...)
		if !bytes.Equal(block, []byte("\n")) {
			wantImport = append(wantImport, block)
		}
	}
	input = append(input, afterBlock...)

	want := &DividedSrc{
		BeforeImportDivision: beforeBlock,
		ImportDivision:       wantImport,
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
	beforeBlock := []byte(`package main

import (
`)

	importBlock := [][]byte{
		[]byte("fmt\n"),
		[]byte("strings\n"),
		[]byte("\n"),
		[]byte("github.com/jinzhu/gorm\n"),
		[]byte("\n"),
		[]byte("github.com/istsh/goimport-fmt/ast\n"),
		[]byte("c \"github.com/istsh/goimport-fmt/config\"\n"),
	}

	afterBlock := []byte(`)

func main() {
	strs := strings.Split("Hello World", " ")
	fmt.Println(strs[0])
	fmt.Println(strs[1])
}
`)

	var input []byte
	input = append(input, beforeBlock...)
	for _, block := range importBlock {
		input = append(input, block...)
	}
	input = append(input, afterBlock...)

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
		if !reflect.DeepEqual(*g, *w) {
			t.Errorf("GetImportDetails index: %d, g: %#v, w %#v", i, g, w)
		}
	}
}
