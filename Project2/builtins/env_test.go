package builtins

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestEnvironmentVariables(t *testing.T) {
	envs := os.Environ()
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		setEnv  map[string]string
		wantOut string
		wantErr error
	}{
		{
			name: "success no args",
			setEnv: map[string]string{
				"ABCD": "FEG",
			},
			wantOut: fmt.Sprintln(strings.Join(append(envs, "ABCD=FEG"), "\n")),
		},
		{
			name: "bad args",
			args: args{
				args: []string{
					"-u",
				},
			},
			wantErr: ErrInvalidArgCount,
		},
		{
			name: "success no args",
			setEnv: map[string]string{
				"ABCD1": "FEG1",
				"ABCD2": "FEG2",
				"ABCD3": "FEG3",
			},
			args: args{
				args: []string{
					"-u", "ABCD1",
					"-u", "ABCD2",
				},
			},
			wantOut: fmt.Sprintln(strings.Join(append(envs, "ABCD3=FEG3"), "\n")),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// setup
			for k, v := range tt.setEnv {
				t.Setenv(k, v)
			}

			// test
			var out bytes.Buffer
			if err := EnvironmentVariables(&out, tt.args.args...); tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("EnvironmentVariables() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			} else if err != nil {
				t.Fatalf("EnvironmentVariables() unexpected error: %v", err)
			}
			if got := out.String(); got != tt.wantOut {
				t.Errorf("EnvironmentVariables() got = %v, want %v", got, tt.wantOut)
			}
		})
	}
}
