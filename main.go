package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var (
		remove = flag.Bool("remove", false, "remove note file")
		list   = flag.Bool("list", false, "list all note files")
	)

	flag.Parse()
	m, err := NewManager()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *list {
		ids, err := m.List()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, id := range ids {
			fmt.Println(id)
		}
		return
	}

	noteID := strings.Join(flag.Args(), " ")
	if *remove {
		exists, err := m.Exists(noteID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if !exists {
			fmt.Fprintf(os.Stderr, "note %q does not exist", noteID)
			os.Exit(1)
		}
		if err := m.Remove(noteID); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	exists, err := m.Exists(noteID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if !exists {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("note %q does not exists. Do you want to create it [Y/n]?", noteID)
		txt, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		txt = strings.ToLower(txt)
		if strings.HasPrefix(txt, "y") {
			if err := m.Create(noteID); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	}

	if err := m.Edit(noteID); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
