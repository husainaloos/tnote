package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	testHomeDir = "/tmp/tnote_test"
)

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
			// clean the test folder
			if err := os.RemoveAll(testHomeDir); err != nil {
				t.Fatalf("cannot remove home dir: %v", err)
			}

			m, _ := NewManager()
			err := m.setHomeDir(testHomeDir)
			if err != nil {
				t.Fatalf("cannot set home dir: %v", err)
			}

			// create all files
			for _, fn := range test.has {
				fn = filepath.Join(testHomeDir, fn)
				f, err := os.Create(fn)
				if err != nil {
					t.Errorf("cannot create file: %v", err)
				}
				defer f.Close()
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
		toCreate string
		expected []string
	}{
		{
			name:     "testing creating simple file",
			toCreate: "file",
			expected: []string{"file.md"},
		},
		{
			name:     "testing simple file with .md",
			toCreate: "file.md",
			expected: []string{"file.md.md"},
		},
		{
			name:     "testing simple file with underscore",
			toCreate: "file_with_underscore",
			expected: []string{"file_with_underscore.md"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := os.RemoveAll(testHomeDir); err != nil {
				t.Fatalf("cannot remove home dir: %v", err)
			}

			m, _ := NewManager()
			err := m.setHomeDir(testHomeDir)
			if err != nil {
				t.Fatalf("cannot set home dir: %v", err)
			}

			if err := m.Create(test.toCreate); err != nil {
				t.Fatalf("failed: %v", err)
			}

			fis, err := ioutil.ReadDir(testHomeDir)
			if err != nil {
				t.Fatalf("cannot read %s: %v", testHomeDir, err)
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
