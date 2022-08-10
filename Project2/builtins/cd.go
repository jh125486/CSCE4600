package builtins

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrInvalidArgCount = errors.New("invalid argument count")
	HomeDir, _         = os.UserHomeDir()
)

func ChangeDirectory(args ...string) error {
	switch len(args) {
	case 0: // change to home directory if available
		if HomeDir == "" {
			return fmt.Errorf("%w: no homedir found, expected one argument (directory)", ErrInvalidArgCount)
		}
		return os.Chdir(HomeDir)
	case 1:
		return os.Chdir(args[0])
	default:
		return fmt.Errorf("%w: expected zero or one arguments (directory)", ErrInvalidArgCount)
	}
}
