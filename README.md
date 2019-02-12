#tnote

A simple command line tool for managing notes.

## Setup

- Ensure that you have `$HOME/Documents` available.
- Configure `$EDITOR` to be your favorite editor.

## Installation

Run `go get -u github.com/husainaloos/tnote`

## Usage

- `tnote --init` to initialize the note directory. This creates the directory `$HOME/Documents/notes`
- `tnote --create homework` will create a note file `homework.md` in the note directory.
- `tnote --remove homework` will remove the note file `homework.md`.
- `tnote homework` will edit `homework.md` file.
- `tnote --list` will list all notes.

(you can use `--init` or `-init`. The number of dashes does not matter)

## FAQS

- Why creating this tool? Looks like just a wrapper.

True. It is just a wrapper. I created it so that I can pull up my notes quickly from anywhere. This is written in a language like go so that I can quickly install it without the hassle.
