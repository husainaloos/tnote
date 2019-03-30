package main

import (
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
