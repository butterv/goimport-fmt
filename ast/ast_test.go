package ast

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/istsh/goimport-fmt/config"
)

func TestAnalyze(t *testing.T) {
	tests := []struct {
		in      []byte
		want    *ImportDetail
		wantErr error
	}{
		{
			in: []byte("context"),
			want: &ImportDetail{
				ImportPath:  []byte("context"),
				PackageType: Standard,
			},
		},
		{
			in: []byte("crypto/rand"),
			want: &ImportDetail{
				ImportPath:  []byte("crypto/rand"),
				PackageType: Standard,
			},
		},
		{
			in: []byte("github.com/jinzhu/gorm"),
			want: &ImportDetail{
				ImportPath:  []byte("github.com/jinzhu/gorm"),
				PackageType: ThirdParty,
			},
		},
		{
			in: []byte("github.com/istsh/goimport-fmt/ast"),
			want: &ImportDetail{
				ImportPath:  []byte("github.com/istsh/goimport-fmt/ast"),
				PackageType: OwnProject,
			},
		},
		{
			in: []byte("github.com/istsh/goimport-fmt/config"),
			want: &ImportDetail{
				ImportPath:  []byte("github.com/istsh/goimport-fmt/config"),
				PackageType: OwnProject,
			},
		},
		{
			in: []byte("github.com/istsh/goimport-fmt/lexer"),
			want: &ImportDetail{
				ImportPath:  []byte("github.com/istsh/goimport-fmt/lexer"),
				PackageType: OwnProject,
			},
		},
	}

	config.Set(os.Getenv("GOROOT"), "", "github.com/istsh/goimport-fmt")
	for _, tt := range tests {
		got, err := Analyze(tt.in)
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("Analyze(%s)=_, %#v; wantErr: %#v", tt.in, err, tt.wantErr)
		}
		if !bytes.Equal(got.Alias, tt.want.Alias) {
			t.Errorf("Analyze Alias got: %s, want: %s", got.Alias, tt.want.Alias)
		}
		if !bytes.Equal(got.ImportPath, tt.want.ImportPath) {
			t.Errorf("Analyze ImportPath got: %s, want: %s", got.ImportPath, tt.want.ImportPath)
		}
		if got.PackageType != tt.want.PackageType {
			t.Errorf("Analyze PackageType got: %d, want: %d", got.PackageType, tt.want.PackageType)
		}
	}
}

func TestAnalyzeIncludeAlias(t *testing.T) {
	type args struct {
		alias []byte
		path  []byte
	}
	tests := []struct {
		in      args
		want    *ImportDetail
		wantErr error
	}{
		{
			in: args{
				path: []byte("context"),
			},
			want: &ImportDetail{
				ImportPath:  []byte("context"),
				PackageType: Standard,
			},
		},
		{
			in: args{
				path: []byte("crypto/rand"),
			},
			want: &ImportDetail{
				ImportPath:  []byte("crypto/rand"),
				PackageType: Standard,
			},
		},
		{
			in: args{
				alias: []byte("rand"),
				path:  []byte("crypto/rand"),
			},
			want: &ImportDetail{
				Alias:       []byte("rand"),
				ImportPath:  []byte("crypto/rand"),
				PackageType: Standard,
			},
		},
		{
			in: args{
				alias: []byte("_"),
				path:  []byte("github.com/go-sql-driver/mysql"),
			},
			want: &ImportDetail{
				Alias:       []byte("_"),
				ImportPath:  []byte("github.com/go-sql-driver/mysql"),
				PackageType: ThirdParty,
			},
		},
		{
			in: args{
				path: []byte("github.com/jinzhu/gorm"),
			},
			want: &ImportDetail{
				ImportPath:  []byte("github.com/jinzhu/gorm"),
				PackageType: ThirdParty,
			},
		},
		{
			in: args{
				path: []byte("github.com/istsh/goimport-fmt/ast"),
			},
			want: &ImportDetail{
				ImportPath:  []byte("github.com/istsh/goimport-fmt/ast"),
				PackageType: OwnProject,
			},
		},
		{
			in: args{
				path: []byte("github.com/istsh/goimport-fmt/config"),
			},
			want: &ImportDetail{
				ImportPath:  []byte("github.com/istsh/goimport-fmt/config"),
				PackageType: OwnProject,
			},
		},
		{
			in: args{
				alias: []byte("l"),
				path:  []byte("github.com/istsh/goimport-fmt/lexer"),
			},
			want: &ImportDetail{
				Alias:       []byte("l"),
				ImportPath:  []byte("github.com/istsh/goimport-fmt/lexer"),
				PackageType: OwnProject,
			},
		},
	}

	config.Set(os.Getenv("GOROOT"), "", "github.com/istsh/goimport-fmt")
	for _, tt := range tests {
		got, err := AnalyzeIncludeAlias(tt.in.alias, tt.in.path)
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("Analyze(%s)=_, %#v; wantErr: %#v", tt.in, err, tt.wantErr)
		}
		if !bytes.Equal(got.Alias, tt.want.Alias) {
			t.Errorf("Analyze Alias got: %s, want: %s", got.Alias, tt.want.Alias)
		}
		if !bytes.Equal(got.ImportPath, tt.want.ImportPath) {
			t.Errorf("Analyze ImportPath got: %s, want: %s", got.ImportPath, tt.want.ImportPath)
		}
		if got.PackageType != tt.want.PackageType {
			t.Errorf("Analyze PackageType got: %d, want: %d", got.PackageType, tt.want.PackageType)
		}
	}
}
