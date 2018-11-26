package ui

import (
	"fmt"
	"io"
	"strings"
)

type UI struct {
	newline   bool
	writer    io.Writer
	reader    io.Reader
	noConfirm bool
}

func NewUI(writer io.Writer, reader io.Reader) *UI {
	return &UI{
		newline:   true,
		writer:    writer,
		reader:    reader,
		noConfirm: false,
	}
}

func (l *UI) clear() {
	if l.newline {
		return
	}

	l.writer.Write([]byte("\n"))
	l.newline = true
}

func (l *UI) Step(message string, a ...interface{}) {
	l.clear()
	fmt.Fprintf(l.writer, "step: %s\n", fmt.Sprintf(message, a...))
	l.newline = true
}

func (l *UI) Printf(message string, a ...interface{}) {
	l.clear()
	fmt.Fprintf(l.writer, "%s", fmt.Sprintf(message, a...))
}

func (l *UI) Println(message string) {
	l.clear()
	fmt.Fprintf(l.writer, "%s\n", message)
}

func (l *UI) NoConfirm() {
	l.noConfirm = true
}

func (l *UI) Prompt(message string) bool {
	if l.noConfirm {
		return true
	}

	l.clear()
	fmt.Fprintf(l.writer, "%s (y/N): ", message)
	l.newline = true

	var proceed string
	fmt.Fscanln(l.reader, &proceed)

	proceed = strings.ToLower(proceed)
	if proceed == "yes" || proceed == "y" {
		return true
	}
	return false
}

func (l *UI) GetInput() string {
	var proceed string
	_, err := fmt.Fscanln(l.reader, &proceed)
	if err != nil {
		// TODO This will procude  error like "unexpected newline", handle it later.
		// os.Exit(1)
	}
	return proceed
}
