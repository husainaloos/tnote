package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func getListCmd() *flag.FlagSet {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listCmd.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintln(w, "Usage: tnote list")
		fmt.Fprintln(w, "list lists all the note IDs the system has. list considers note IDs that live in nested folders")
	}
	return listCmd
}

func getEditCmd() *flag.FlagSet {
	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	editCmd.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintln(w, "Usage: tnote edit <noteID>")
		fmt.Fprintln(w, "edit edits a given noteID. noteID can be a relative path (e.g. foo/bar). noteID cannot have space character.")
	}
	return editCmd
}
func getRemoveCmd() *flag.FlagSet {
	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	removeCmd.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintln(w, "Usage: tnote remove <noteID>")
		fmt.Fprintln(w, "remove removes a given noteID. noteID can be a relative path (e.g. foo/bar). noteID cannot have space character.")
	}
	return removeCmd
}

func main() {
	var (
		listCmd   = getListCmd()
		editCmd   = getEditCmd()
		removeCmd = getRemoveCmd()
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
