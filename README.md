# tnote

A simple command line tool for managing notes. I created this tool because I needed a simple way to manage my notes.

## Setup

- Configure `$EDITOR` to be your favorite editor.

## Installation

Run `go install -u github.com/husainaloos/tnote`

## Usage

- `tnote mynotes` will edit `mynotes.md` file. If the file does not exists, it will create it.
- `tnote --remove mynotes` will remove the note file `mynotes.md`.
- `tnote --list` will list all notes.
- `tnote --help` for help.

## Integrating with FZF

Currently `tnote --list` outputs the list of notes. If you are trying to find a given note, this might not be the most helpful way. Instead, I suggest that you integrate this with FZF or any fuzzy finder. Just add 
```bash
alias fnote='tnote $(tnote --list | fzf)
``` 
to your shell startup file to find a note with the help of FZF. If the note you typed does not exist, tnote will create it.

## Future Features

I will be adding features as I need to.
