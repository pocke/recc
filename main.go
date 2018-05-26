package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/atotto/clipboard"
	"github.com/kr/pty"
)

func main() {
	exitStatus, err := Main(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(exitStatus)
	}
	os.Exit(exitStatus)
}

func Main(args []string) (int, error) {
	opt, err := OptionParse(args)
	if err != nil {
		return 1, err
	}

	cmd := opt.Args[0]
	cmdArgs := opt.Args[1:]
	var c *exec.Cmd
	if strings.Contains(cmd, " ") && len(cmdArgs) == 0 {
		c = exec.Command("bash", "-c", cmd)
	} else {
		c = exec.Command(cmd, cmdArgs...)
	}

	exitStatus := 0
	out, err := Exec(c, opt)
	if err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			if s, isWaitStatus := err.Sys().(syscall.WaitStatus); isWaitStatus {
				exitStatus = s.ExitStatus()
			} else {
				panic("Recc does not support this OS.")
			}
		case *exec.Error:
			if err.Err == exec.ErrNotFound {
				return 127, err
			} else {
				return 1, err
			}
		default:
			return 1, err
		}
	}

	cmdLine := fmt.Sprintf("$ %s %s\n", cmd, strings.Join(cmdArgs, " "))
	out = cmdLine + out

	if opt.Output != "" {
		return exitStatus, ioutil.WriteFile(opt.Output, []byte(out), 0644)
	} else {
		return exitStatus, clipboard.WriteAll(out)
	}
}

func Exec(c *exec.Cmd, opt *Option) (string, error) {
	if opt.Tty {
		return WithTty(c, opt)
	} else {
		return WithoutTty(c, opt)
	}
}

func WithTty(c *exec.Cmd, opt *Option) (string, error) {
	t, err := pty.Start(c)
	if err != nil {
		return "", err
	}
	defer t.Close()
	go func() {
		// io.Copy(t, os.Stdin)
		t.Write([]byte{4})
	}()
	if err := c.Wait(); err != nil {
		return "", err
	}
	b := bytes.NewBuffer([]byte{})
	io.Copy(b, t)
	str := b.String()
	fmt.Print(str)
	return str, nil
}

func WithoutTty(c *exec.Cmd, opt *Option) (string, error) {
	r := NewRecorder(os.Stdout, os.Stderr, "")
	c.Stdin = os.Stdin
	c.Stdout = r.Stdout
	if opt.Stderr {
		c.Stderr = r.Stderr
	} else {
		c.Stderr = os.Stderr
	}
	err := c.Run()

	return r.String(), err
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

func (r *Recorder) Bytes() []byte {
	return r.record.Bytes()
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
