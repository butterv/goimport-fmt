package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/istsh/goimport-fmt/config"
	"github.com/istsh/goimport-fmt/lexer"
)

func Test_createResult_OnlyOneImport(t *testing.T) {
	want := []byte(`package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}
`)

	beforeBlock := []byte(`package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}
`)
	ds := &lexer.DividedSrc{
		BeforeImportDivision: beforeBlock,
		ImportDivision:       nil,
		AfterImportDivision:  nil,
	}

	got, err := createResult(ds)
	if err != nil {
		t.Fatalf("createResult(ds)=_, %#v; want nil", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("createResult(ds)=%s, nil; want %s", got, want)
	}
}

func Test_createResult_OnlyOneImportWithBrackets(t *testing.T) {
	want := []byte(`package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World")
}
`)

	beforeBlock := []byte(`package main

import (
`)

	importBlock := [][]byte{
		[]byte("\t\"fmt\"\n"),
	}

	afterBlock := []byte(`)

func main() {
	fmt.Println("Hello World")
}
`)

	ds := &lexer.DividedSrc{
		BeforeImportDivision: beforeBlock,
		ImportDivision:       importBlock,
		AfterImportDivision:  afterBlock,
	}

	config.Set(os.Getenv("GOROOT"), "", "github.com/istsh/goimport-fmt")
	got, err := createResult(ds)
	if err != nil {
		t.Fatalf("createResult(ds)=_, %#v; want nil", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("createResult(ds)=%s, nil; want %s", got, want)
	}
}

func Test_createResult_SomeImport(t *testing.T) {
	want := []byte(`package main

import "fmt"
import "strings"
import g "github.com/jinzhu/gorm"
import "github.com/istsh/goimport-fmt/ast"
import c "github.com/istsh/goimport-fmt/config"

func main() {
	strs := strings.Split("Hello World", " ")
	fmt.Println(strs[0])
	fmt.Println(strs[1])
}
`)

	beforeBlock := []byte(`package main

import "fmt"
import "strings"
import g "github.com/jinzhu/gorm"
import "github.com/istsh/goimport-fmt/ast"
import c "github.com/istsh/goimport-fmt/config"

func main() {
	strs := strings.Split("Hello World", " ")
	fmt.Println(strs[0])
	fmt.Println(strs[1])
}
`)

	ds := &lexer.DividedSrc{
		BeforeImportDivision: beforeBlock,
		ImportDivision:       nil,
		AfterImportDivision:  nil,
	}

	got, err := createResult(ds)
	if err != nil {
		t.Fatalf("createResult(ds)=_, %#v; want nil", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("createResult(ds)=%s, nil; want %s", got, want)
	}
}

func Test_createResult_SomeImportWithBrackets(t *testing.T) {
	want := []byte(`package main

import (
	"fmt"
	"strings"

	g "github.com/jinzhu/gorm"

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
		[]byte("\tg \"github.com/jinzhu/gorm\"\n"),
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

	ds := &lexer.DividedSrc{
		BeforeImportDivision: beforeBlock,
		ImportDivision:       importBlock,
		AfterImportDivision:  afterBlock,
	}

	config.Set(os.Getenv("GOROOT"), "", "github.com/istsh/goimport-fmt")
	got, err := createResult(ds)
	if err != nil {
		t.Fatalf("createResult(ds)=_, %#v; want nil", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("createResult(ds)=%s, nil; want %s", got, want)
	}
}
