# tnote

A simple command line tool for managing notes. I created this tool because I needed a simple way to manage my notes.

## Setup

- Configure `$EDITOR` to be your favorite editor.

## Installation

Run `go install -u github.com/husainaloos/tnote`

## Usage

- `tnote homework` will edit `homework.md` file. If the file does not exists, it will create it.
- `tnote --remove homework` will remove the note file `homework.md`.
- `tnote --list` will list all notes.
- `tnote --help` for help.

(you can use `--init` or `-init`. The number of dashes does not matter)
