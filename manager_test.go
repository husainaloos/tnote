package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sync/atomic"
	"testing"
)

var (
	folder_id int64 = 0
)

func getManager(t *testing.T) (m *Manager, homeDir string) {
	m, err := NewManager()
	if err != nil {
		t.Fatalf("cannot create manager: %v", err)
	}

	homeDir = "/tmp/test_tnote_" + string(atomic.AddInt64(&folder_id, 1))
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, dir := getManager(t)
			defer func() { _ = os.RemoveAll(dir) }()

			// create all files
			for _, fn := range test.has {
				fn = filepath.Join(dir, fn)
				f, err := os.Create(fn)
				if err != nil {
					t.Errorf("cannot create file: %v", err)
				}
				f.Close()
			}

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

			for _, fn := range test.has {
				fn = filepath.Join(dir, fn)
				f, err := os.Create(fn)
				if err != nil {
					t.Fatalf("cannot create file %s: %v", fn, err)
				}
				f.Close()
			}

			if err := m.Create(test.noteid); (err != nil) != test.err {
				t.Fatalf("got err=%v, exptected err=%v", err, test.err)
			}

			fis, err := ioutil.ReadDir(dir)
			if err != nil {
				t.Fatalf("cannot read %s: %v", dir, err)
			}
			fns := make([]string, 0)
			for _, fi := range fis {
				fns = append(fns, fi.Name())
			}
			if !reflect.DeepEqual(fns, test.expected) {
				t.Fatalf("have %v in dir, but expected %v", fns, test.expected)
			}
		})
	}
}
