package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

// Manager for notes that uses /Documents/notes directory
type Manager struct {
	homeDir string
	editor  string
}

// NewManager creates an instance of the manager
func NewManager() (*Manager, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	home := fmt.Sprintf("%s/Documents/notes", usr.HomeDir)
	if err := os.Mkdir(home, 0755); !os.IsExist(err) {
		return nil, err
	}
	editor := os.Getenv("EDITOR")
	return &Manager{
		homeDir: home,
		editor:  editor,
	}, nil
}

// List the notes available
func (m *Manager) List() ([]string, error) {
	fis, err := ioutil.ReadDir(m.homeDir)
	if err != nil {
		return nil, err
	}
	noteIDs := make([]string, 0)
	for _, fi := range fis {
		noteID := strings.TrimSuffix(fi.Name(), ".md")
		noteIDs = append(noteIDs, noteID)
	}
	return noteIDs, nil
}

// Remove a note by note id
func (m *Manager) Remove(noteID string) error {
	p := m.getNotePath(noteID)
	return os.Remove(p)
}

// Create a new note by note id
func (m *Manager) Create(noteID string) error {
	p := m.getNotePath(noteID)
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte("#"))
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
	p := fmt.Sprintf("%s/%s.md", m.homeDir, noteID)
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
