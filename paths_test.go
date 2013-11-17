package paths

import (
	"os"
	"path/filepath"
	"strings"
	. "testing"
)

func Test_FindVCSRoot(t *T) {
	if Short() {
		t.SkipNow()
	}

	var err error
	testdir := filepath.Join(os.TempDir(), "findvcsroottest")
	folder1 := filepath.Join(testdir, "project1", ".git")
	folder2 := filepath.Join(testdir, "project2", ".hg")
	subfolder := filepath.Join(testdir, "project2", "subfolder")
	defer os.RemoveAll(testdir)

	for _, dir := range []string{folder1, folder2, subfolder} {
		err = os.MkdirAll(dir, 0770)
		if err != nil {
			t.Error("Error creating dir:", dir, err)
		}
	}

	kind, p, err := FindVCSRoot(os.TempDir())
	if err != nil {
		t.Error("Unexpected error:", err)
	} else if kind != NONE || len(p) != 0 {
		t.Error("Expected empty return looking for root that doesn't exist.")
	}

	kind, p, err = FindVCSRoot(folder1)
	if err != nil {
		t.Error("Unexpected error:", err)
	} else if kind != GIT {
		t.Error("Expected it to find git, got:", kind)
	} else if p != WalkUpPath(folder1) {
		t.Error("Expected the correct path to be returned.")
	}

	kind, p, err = FindVCSRoot(subfolder)
	if err != nil {
		t.Error("Unexpected error:", err)
	} else if kind != HG {
		t.Error("Expected it to find hg, got:", kind)
	} else if p != WalkUpPath(folder2) {
		t.Error("Expected the correct path to be returned.")
	}
}

func Test_WalkUpPath(t *T) {
	var tests = []struct {
		In  string
		Out string
	}{
		{"", ""},
		{"/", ""},
		{"/home", "/"},
		{"/home/", "/"},
		{"/home/user", "/home"},
		{"/home/user/", "/home"},
		{"/home/user/file.txt", "/home/user"},
		{"home", ""},
	}
	for _, test := range tests {
		if result := WalkUpPath(test.In); result != test.Out {
			t.Errorf(`Expected: "%s" -> "%s" but got: "%s"`,
				test.In, test.Out, result)
		}
	}
}

func Test_DirAndFileExists(t *T) {
	if Short() {
		t.SkipNow()
	}
	var exist bool
	var err error
	testdir := filepath.Join(os.TempDir(), "dirandfileexists")
	testfile := filepath.Join(testdir, "testfile.txt")

	_, err = os.Stat(testdir)
	if err == nil || !os.IsNotExist(err) {
		t.Error("Expected the folder to not exist:", err)
	}

	if exist, err = DirExists(testdir); err != nil {
		t.Error("Unexpected error:", err)
	} else if exist {
		t.Error("Expected dir to not exist:", testdir)
	}

	if exist, err = FileExists(testfile); err != nil {
		t.Error("Unexpected error:", err)
	} else if exist {
		t.Error("Expected file to not exist:", testfile)
	}

	err = os.Mkdir(testdir, 0770)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	defer os.RemoveAll(testdir)

	if exist, err = DirExists(testdir); err != nil {
		t.Error("Unexpected error:", err)
	} else if !exist {
		t.Error("Expected an existing dir:", testdir)
	}

	f, err := os.Create(testfile)
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	err = f.Close()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if exist, err = FileExists(testfile); err != nil {
		t.Error("Unexpected error:", err)
	} else if !exist {
		t.Error("Expected file to exist:", testfile)
	}

	exist, err = DirExists(testfile)
	if err == nil || !strings.Contains(err.Error(), "dir, but found file") {
		t.Error("Expected an error due to not being dir, but got:", err)
	}

	exist, err = FileExists(testdir)
	if err == nil || !strings.Contains(err.Error(), "file, but found dir") {
		t.Error("Expected an error due to not being file, but got:", err)
	}
}

func Test_EnsureDirectory(t *T) {
	if Short() {
		t.SkipNow()
	}
	testdir := filepath.Join(os.TempDir(), "ensuredirectorytest")
	defer os.RemoveAll(testdir)

	created, err := EnsureDirectory(testdir)
	if err != nil {
		t.Error("Unexpected Error:", err)
	}
	if !created {
		t.Error("Expected the folder to be created.")
	}

	_, err = os.Stat(testdir)
	if os.IsNotExist(err) {
		t.Error("Expected the folder to be created.")
	}

	created, err = EnsureDirectory(testdir)
	if err != nil {
		t.Error("Unexpected Error:", err)
	}
	if created {
		t.Error("Expected the folder to exist.")
	}
}
