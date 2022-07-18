package builtins

import (
	"errors"
	"fmt"
	"os"
)

var ErrInvalidArgCount = errors.New("invalid argument count")

func ChangeDirectory(args ...string) error {
	switch len(args) {
	case 0: // change to home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		return os.Chdir(homeDir)
	case 1:
		return os.Chdir(args[0])
	default:
		return fmt.Errorf("%w: expected one argument (directory)", ErrInvalidArgCount)
	}
}
