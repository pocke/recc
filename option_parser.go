package main

import "errors"

type Option struct {
	Output string
	Stderr bool
	Tty    bool

	Args []string
}

var errHelp = errors.New(`
Usage: reec [OPTION]... COMMAND [COMMAND_ARGS]...
       recc [OPTION]... 'COMMAND [COMMAND_ARGS]...'

Options:
	-o, --output FILE_NAME output file name
	--stderr               include Standard Error to result
	--tty                  TTY
	--help                 Display this message
`)

// --output FILE
// --stderr
// --help
func OptionParse(args []string) (*Option, error) {
	if len(args) == 1 {
		return nil, errHelp
	}

	o := &Option{}
	ptr := 1

	for {
		if ptr >= len(args) {
			return nil, errHelp
		}

		switch args[ptr] {
		case "-o", "--output":
			ptr++
			if ptr >= len(args) {
				return nil, errHelp
			}
			o.Output = args[ptr]
		case "--stderr":
			o.Stderr = true
		case "--tty", "-t":
			o.Tty = true
		case "--help":
			return nil, errHelp
		default:
			goto endLoop
		}
		ptr++
	}
endLoop:
	o.Args = args[ptr:]

	return o, nil
}
