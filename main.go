package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	var (
		create = flag.Bool("create", false, "create a note file")
		remove = flag.Bool("remove", false, "remove a note file")
	)
	flag.Parse()

	if *create {
		fname, err := getFileName(os.Args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		log.Println(fname)
		if err := createFile(fname); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if *remove {
		fname, err := getFileName(os.Args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := removeFile(fname); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fname, err := getFileName(os.Args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := editFile(fname); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func editFile(fname string) error {
	path := fmt.Sprintf("%s/Documents/notes/%s.md", os.Getenv("HOME"), fname)
	editor := os.Getenv("EDITOR")
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func createFile(fname string) error {
	path := fmt.Sprintf("%s/Documents/notes", os.Getenv("HOME"))
	if err := os.Mkdir(path, 0755); err != nil {
		if os.IsExist(err) {
		} else {
			return err
		}
	}

	path = fmt.Sprintf("%s/%s.md", path, fname)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func removeFile(fname string) error {
	path := fmt.Sprintf("%s/Documents/notes/%s.md", os.Getenv("HOME"), fname)
	return os.Remove(path)
}

func getFileName(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("too few arguments")
	}

	if len(args) == 2 {
		return args[1], nil
	} else {
		return args[2], nil
	}
}
