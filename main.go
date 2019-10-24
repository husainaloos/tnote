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
		listCmd   = flag.NewFlagSet("list", flag.ExitOnError)
		editCmd   = flag.NewFlagSet("edit", flag.ExitOnError)
		removeCmd = flag.NewFlagSet("remove", flag.ExitOnError)
	)
	flag.Parse()

	m, err := NewManager()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing command.")
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "list":
		err := listCmd.Parse(flag.Args()[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		ids, err := m.List()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, id := range ids {
			fmt.Println(id)
		}
		return
	case "edit":
		err := editCmd.Parse(flag.Args()[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(editCmd.Args())
		if editCmd.NArg() != 1 {
			fmt.Fprintln(os.Stderr, "missing file argument")
			os.Exit(1)
		}
		noteID := editCmd.Arg(0)
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
	case "remove":
		err := removeCmd.Parse(flag.Args()[1:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if removeCmd.NArg() != 1 {
			fmt.Fprintln(os.Stderr, "missing file argument")
			os.Exit(1)
		}
		noteID := editCmd.Arg(0)
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
}
