package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func main() {
	var (
		init   = flag.Bool("init", false, "create note folder")
		create = flag.Bool("create", false, "create note file")
		remove = flag.Bool("remove", false, "remove note file")
		list   = flag.Bool("list", false, "list all note files")
	)

	flag.Parse()
	if *init {
		if err := createFolder(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	if *create {
		filepath, err := getFileName()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := createFile(filepath); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	if *remove {
		filepath, err := getFileName()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := removeFile(filepath); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	if *list {
		files, err := listFiles()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(files)
		return
	}

	filepath, err := getFileName()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := editFile(filepath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func listFiles() ([]string, error) {
	path := fmt.Sprintf("%s/Documents/notes/", os.Getenv("HOME"))
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0)
	for _, fi := range fis {
		files = append(files, fi.Name())
	}
	return files, nil
}

func createFolder() error {
	path := fmt.Sprintf("%s/Documents/notes", os.Getenv("HOME"))
	return os.Mkdir(path, 0755)
}

func editFile(filepath string) error {
	editor := os.Getenv("EDITOR")
	cmd := exec.Command(editor, filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func createFile(filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func removeFile(filepath string) error {
	return os.Remove(filepath)
}

func getFileName() (string, error) {
	if flag.NArg() != 1 {
		return "", fmt.Errorf("received %d file names, but expected one file name", flag.NArg())
	}

	filepath := fmt.Sprintf("%s/Documents/notes/%s.md", os.Getenv("HOME"), flag.Arg(0))
	return filepath, nil

}
