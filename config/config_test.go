package config

import "testing"

func TestGetGoRoot(t *testing.T) {
	want := "TEST_GOROOT"
	Set(want, "TEST_FILE_PATH", "TEST_OWN_PROJECT")

	got := GetGoRoot()
	if got != want {
		t.Fatalf("GetGoRoot()=%s; want: %s", got, want)
	}
}

func TestGetFilePath(t *testing.T) {
	want := "TEST_FILE_PATH"
	Set("TEST_GOROOT", want, "TEST_OWN_PROJECT")

	got := GetFilePath()
	if got != want {
		t.Fatalf("GetFilePath()=%s; want: %s", got, want)
	}
}

func TestGetOwnProject(t *testing.T) {
	want := "TEST_OWN_PROJECT"
	Set("TEST_GOROOT", "TEST_FILE_PATH", want)

	got := GetOwnProject()
	if got != want {
		t.Fatalf("GetOwnProject()=%s; want: %s", got, want)
	}
}
