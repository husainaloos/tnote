package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sync/atomic"
	"testing"
)

var (
	folderID int64
)

// createFilesForTest creates files in a given directory
// used for testing
func createFilesForTest(t *testing.T, dir string, fs []string) {
	for _, fn := range fs {
		fn = filepath.Join(dir, fn)
		parent := path.Dir(fn)
		_ = os.Mkdir(parent, os.ModePerm)
		f, err := os.Create(fn)
		if err != nil {
			t.Fatalf("cannot create file %s: %v", fn, err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("cannot close file %s: %v", fn, err)
		}
	}
}

// getFilesInDir get the list of files in a given directory
// used for testing
func getFilesInDir(t *testing.T, dir string) []string {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatalf("cannot read %s: %v", dir, err)
	}
	fns := make([]string, 0)
	for _, fi := range fis {
		fns = append(fns, fi.Name())
	}

	return fns
}

func getManager(t *testing.T) (m *Manager, homeDir string) {
	m, err := NewManager()
	if err != nil {
		t.Fatalf("cannot create manager: %v", err)
	}

	homeDir = "/tmp/test_tnote_" + string(atomic.AddInt64(&folderID, 1))
	if err := m.setHomeDir(homeDir); err != nil {
		t.Fatalf("cannot set homeDir %s: %v", homeDir, err)
	}

	return
}

func TestCreatingManager(t *testing.T) {
	manager, err := NewManager()
	if err != nil {
		t.Fatalf("expected no error, but found %v", err)
	}
	if manager == nil {
		t.Fatalf("expected non-nil value, but received nil")
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name     string
		has      []string
		expected []string
	}{
		{
			name:     "testing simple files",
			has:      []string{"file1.md", "file2.md"},
			expected: []string{"file1", "file2"},
		},
		{
			name:     "testing files with dashes",
			has:      []string{"file-with-dash.md"},
			expected: []string{"file-with-dash"},
		},
		{
			name:     "testing files with underscore",
			has:      []string{"file_with_underscore.md"},
			expected: []string{"file_with_underscore"},
		},
		{
			name:     "testing files with two dots",
			has:      []string{"file_with_two_mds.md.md"},
			expected: []string{"file_with_two_mds.md"},
		},
		{
			name:     "testing files in folders",
			has:      []string{"folder/file1.md", "file2.md"},
			expected: []string{"file2", "folder/file1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, dir := getManager(t)
			defer func() { _ = os.RemoveAll(dir) }()
			createFilesForTest(t, dir, test.has)
			got, err := m.List()
			if err != nil {
				t.Errorf("failed to list files: %v", err)
			}
			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("got %v, but expected %v", got, test.expected)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name     string
		has      []string
		noteid   string
		expected []string
		err      bool
	}{
		{
			name:     "testing creating simple file",
			has:      []string{},
			noteid:   "file",
			expected: []string{"file.md"},
			err:      false,
		},
		{
			name:     "testing simple file with .md",
			has:      []string{},
			noteid:   "file.md",
			expected: []string{"file.md.md"},
			err:      false,
		},
		{
			name:     "testing simple file with underscore",
			has:      []string{},
			noteid:   "file_with_underscore",
			expected: []string{"file_with_underscore.md"},
			err:      false,
		},
		{
			name:     "fail to create duplicate file",
			has:      []string{"file.md"},
			noteid:   "file",
			expected: []string{"file.md"},
			err:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			m, dir := getManager(t)
			defer func() { _ = os.RemoveAll(dir) }()
			createFilesForTest(t, dir, test.has)

			if err := m.Create(test.noteid); (err != nil) != test.err {
				t.Fatalf("got err=%v, exptected err=%v", err, test.err)
			}
			if test.err {
				return
			}

			fns := getFilesInDir(t, dir)
			if !reflect.DeepEqual(fns, test.expected) {
				t.Fatalf("have %v in dir, but expected %v", fns, test.expected)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name     string
		has      []string
		noteid   string
		expected []string
		err      bool
	}{
		{
			name:     "should remove simple file",
			has:      []string{"file.md", "extrafile.md"},
			noteid:   "file",
			expected: []string{"extrafile.md"},
			err:      false,
		},
		{
			name:     "should remove simple file even if there is a similarly named file",
			has:      []string{"file.md", "file_.md", "extrafile.md"},
			noteid:   "file",
			expected: []string{"extrafile.md", "file_.md"},
			err:      false,
		},
		{
			name:     "should fails to remove non-existing file",
			has:      []string{"extrafile.md"},
			noteid:   "file",
			expected: nil,
			err:      true,
		},
	}

	for _, test := range tests {
		m, dir := getManager(t)
		defer func() { _ = os.RemoveAll(dir) }()
		createFilesForTest(t, dir, test.has)
		if err := m.Remove(test.noteid); (err != nil) != test.err {
			t.Fatalf("got err=%v, expected err=%v", err, test.err)
		}
		if test.err {
			return
		}
		fns := getFilesInDir(t, dir)
		if !reflect.DeepEqual(fns, test.expected) {
			t.Fatalf("got %v, but expected %v", fns, test.expected)
		}
	}
}

func TestExists(t *testing.T) {
	tests := []struct {
		name   string
		has    []string
		noteid string
		expect bool
		err    bool
	}{
		{
			name:   "file should exists",
			has:    []string{"file.md"},
			noteid: "file",
			expect: true,
			err:    false,
		},
		{
			name:   "file should not exists",
			has:    []string{"file.md"},
			noteid: "none_existing_note",
			expect: false,
			err:    false,
		},
		{
			name:   "file should not exists event if extension is provided",
			has:    []string{"file.md"},
			noteid: "file.md",
			expect: false,
			err:    false,
		},
	}

	for _, test := range tests {
		m, dir := getManager(t)
		defer func() { _ = os.RemoveAll(dir) }()
		createFilesForTest(t, dir, test.has)
		got, err := m.Exists(test.noteid)
		if err != nil {
			if !test.err {
				t.Fatalf("got err: %v", err)
			}
			return
		}
		if test.err {
			t.Fatalf("got no error, but expected err")
		}
		if got != test.expect {
			t.Fatalf("got %v, expected %v", got, test.expect)
		}
	}
}
