package pod

import (
	"archive/zip"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
)

type FileInfo struct {
	Name  string
	CRC32 uint32
	Dir   bool
}

// Reference for files and directories on mocks/ref/
var Reference = map[string]*FileInfo{
	"dir/": &FileInfo{
		Name:  "dir/",
		CRC32: 0,
		Dir:   true,
	},
	"dir/placeholder": &FileInfo{
		Name:  "dir/placeholder",
		CRC32: 258472330,
		Dir:   false,
	},
	"dir2/NotIgnored.md": &FileInfo{
		Name:  "dir2/NotIgnored.md",
		CRC32: 2286427926,
		Dir:   false,
	},
	"doc": &FileInfo{
		Name:  "doc",
		CRC32: 3402152999,
		Dir:   false,
	},
	"ignored": &FileInfo{
		Name:  "ignored",
		CRC32: 763821341,
		Dir:   false,
	},
}

func TestCompress(t *testing.T) {
	var ignoredList = []string{
		"ignored",
		"dir/foo/another_ignored_dir",
		"ignored_dir",
		"**/complex/*",
		"*Ignored*.md",
		"!NotIgnored.md",
	}

	var err = Compress(
		"mocks/res/compress.zip",
		"mocks/ref",
		ignoredList,
	)

	if err != nil {
		t.Errorf("Expected compress to end without errors, got %v error instead", err)
	}

	r, err := zip.OpenReader("mocks/res/compress.zip")
	defer r.Close()

	if err != nil {
		t.Errorf("Wanted no errors opening compressed file, got %v instead", err)
	}

	var found = map[string]*FileInfo{}

	for _, f := range r.File {
		found[f.Name] = &FileInfo{
			Name:  f.Name,
			CRC32: f.CRC32,
			Dir:   f.FileInfo().IsDir(),
		}
	}

	var assertNotIgnored = []string{
		"doc",
		"dir/",
		"dir/placeholder",
	}

	var refuteIgnored = []string{
		"ignored",
		"ignored_dir",
		"ignored_dir/placeholder",
		"dir/foo/another_ignored_dir",
		"dir/foo/another_ignored_dir/placeholder",
		"dir/sub/complex",
		"dir/sub/complex/placeholder",
		"dir/sub/complex/dir/placeholder",
	}

	for _, k := range assertNotIgnored {
		if found[k] == nil {
			t.Errorf("%v not found on zip", k)
		}

		if !reflect.DeepEqual(found[k], Reference[k]) {
			t.Errorf("Zipped %v headers doesn't match expected values", k)
		}
	}

	for _, k := range refuteIgnored {
		if found[k] != nil {
			t.Errorf("Expected file %v to be ignored.", k)
		}
	}
}

func TestCompressInvalidDestination(t *testing.T) {
	var invalid = fmt.Sprintf("mocks/res/invalid-dest-%d/foo.pod", rand.Int())
	var err = Compress(invalid, "mocks/ref", nil)

	if !os.IsNotExist(err) {
		t.Errorf("Wanted error te be due directory not found, got %v instead", err)
	}
}