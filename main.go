package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
)

func main() {
	if err := Main(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

}

func Main(args []string) error {
	cmd := args[1]
	cmdArgs := args[2:]
	c := exec.Command(cmd, cmdArgs...)
	cmdLine := fmt.Sprintf("$ %s %s\n", cmd, strings.Join(cmdArgs, " "))
	r := NewRecorder(os.Stdout, os.Stderr, cmdLine)
	c.Stdin = os.Stdin
	c.Stdout = r.Stdout
	c.Stderr = r.Stderr

	err := c.Run()
	if err != nil {
		return err
	}
	return clipboard.WriteAll(r.String())
}

type Recorder struct {
	Stdout io.Writer
	Stderr io.Writer
	record *bytes.Buffer
}

func NewRecorder(stdout, stderr io.Writer, initial string) *Recorder {
	rec := bytes.NewBuffer([]byte(initial))
	r := &Recorder{
		record: rec,
		Stdout: NewPipe(rec, stdout),
		Stderr: NewPipe(rec, stderr),
	}

	return r
}

func (r *Recorder) String() string {
	return r.record.String()
}

type Pipe struct {
	record io.Writer
	out    io.Writer
}

func NewPipe(rec, out io.Writer) io.Writer {
	return &Pipe{
		record: rec,
		out:    out,
	}
}

func (p *Pipe) Write(b []byte) (int, error) {
	p.record.Write(b)
	return p.out.Write(b)
}
