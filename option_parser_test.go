package main

import (
	"reflect"
	"testing"
)

func TestOptionParse(t *testing.T) {
	cases := []struct {
		args     []string
		expected *Option
		err      error
	}{
		{
			[]string{"recc", "ls"},
			&Option{Args: []string{"ls"}},
			nil,
		},
		{
			[]string{"./recc", "ls"},
			&Option{Args: []string{"ls"}},
			nil,
		},
		{
			[]string{"recc", "--output", "ls"},
			nil,
			errHelp,
		},
		{
			[]string{"recc", "--output", "hoge.out", "ls"},
			&Option{Args: []string{"ls"}, Output: "hoge.out"},
			nil,
		},
	}

	for _, c := range cases {
		option, err := OptionParse(c.args)
		if err != c.err {
			t.Errorf("Expected: %v, but got %v", c.err, err)
		}
		if !reflect.DeepEqual(option, c.expected) {
			t.Errorf("Expected: %v, but got %v", c.expected, option)
		}
	}
}
