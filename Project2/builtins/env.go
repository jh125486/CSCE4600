package builtins

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func EnvironmentVariables(w io.Writer, args ...string) error {
	toRemove := make([]string, 0)
	for i := 0; i < len(args); i++ {
		if args[i] == "-u" {
			if len(args) < i+2 {
				return fmt.Errorf("%w: -u requires an argument", ErrInvalidArgCount)
			}
			toRemove = append(toRemove, args[i+1])
			i++
		}
	}

	toShow := make([]string, 0)
	for _, env := range os.Environ() {
		show := true
		name := strings.Split(env, "=")[0]
		for _, v := range toRemove {
			if v == name {
				show = false
				break
			}
		}
		if show {
			toShow = append(toShow, env)
		}
	}

	_, err := fmt.Fprintln(w, strings.Join(toShow, "\n"))

	return err
}
