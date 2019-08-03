package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
)

// Manager for notes that uses /Documents/notes directory
type Manager struct {
	homeDir string
	editor  string
	perm    os.FileMode
}

// NewManager creates an instance of the manager
func NewManager() (*Manager, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	home := fmt.Sprintf("%s/Documents/notes", usr.HomeDir)
	var perm os.FileMode = 0644
	if err := os.Mkdir(home, perm); !os.IsExist(err) {
		return nil, err
	}
	editor := os.Getenv("EDITOR")
	return &Manager{
		homeDir: home,
		editor:  editor,
		perm:    perm,
	}, nil
}

// setHomeDir for the manager. This is used only for testing
// and should not be used for any other reason
func (m *Manager) setHomeDir(dir string) error {
	err := os.MkdirAll(dir, 0744)
	if err != nil {
		return err
	}
	m.homeDir = dir
	return nil
}

// List the notes available
func (m *Manager) List() ([]string, error) {
	noteIDs := make([]string, 0)
	err := filepath.Walk(m.homeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".md") {
			noteID := strings.TrimSuffix(path, ".md")
			noteID = strings.TrimPrefix(noteID, m.homeDir+"/")
			noteIDs = append(noteIDs, noteID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(noteIDs)
	return noteIDs, nil
}

// Remove a note by note id
func (m *Manager) Remove(noteID string) error {
	p := m.getNotePath(noteID)
	return os.Remove(p)
}

// Create a new note by note id
func (m *Manager) Create(noteID string) error {
	ok, err := m.Exists(noteID)
	if err != nil {
		return err
	}
	if ok {
		return errors.New("cannot create duplicate file")
	}
	p := m.getNotePath(noteID)
	d := filepath.Dir(p)
	if err := os.MkdirAll(d, m.perm); err != nil {
		return err
	}
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	err = f.Chmod(m.perm)
	return err
}

// Exists checks if the note exits
func (m *Manager) Exists(noteID string) (bool, error) {
	p := m.getNotePath(noteID)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// Edit the note given a note id
func (m *Manager) Edit(noteID string) error {
	p := m.getNotePath(noteID)
	if _, err := os.Stat(p); err != nil {
		return err
	}
	cmd := exec.Command(m.editor, p)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (m *Manager) getNotePath(noteID string) string {
	return fmt.Sprintf("%s/%s.md", m.homeDir, noteID)
}
