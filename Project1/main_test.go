package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
)

func TestFCFSSchedule(t *testing.T) {
	t.Parallel()
	type args struct {
		processes []Process
		title     string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name: "default",
			args: args{
				processes: []Process{
					{
						ProcessID:     1,
						ArrivalTime:   0,
						BurstDuration: 5,
						Priority:      2,
					},
					{
						ProcessID:     2,
						ArrivalTime:   3,
						BurstDuration: 9,
						Priority:      1,
					},
					{
						ProcessID:     3,
						ArrivalTime:   6,
						BurstDuration: 6,
						Priority:      3,
					},
				},
				title: "First-come, First-serve",
			},
			wantOut: loadFixture(t, "fcfs_test.txt"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var w bytes.Buffer
			FCFSSchedule(&w, tt.args.title, tt.args.processes)
			if got := w.String(); got != tt.wantOut {
				t.Errorf("FCFSSchedule() = %v, want %v", got, tt.wantOut)
			}
		})
	}
}

func Test_loadProcesses(t *testing.T) {
	t.Parallel()
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []Process
		wantErr error
	}{
		{
			name: "bad CSV",
			args: args{
				r: iotest.ErrReader(io.ErrUnexpectedEOF),
			},
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name: "success",
			args: args{
				r: strings.NewReader(`1,5,0,2
2,9,3,1
3,6,3,3`),
			},
			want: []Process{
				{
					ProcessID:     1,
					ArrivalTime:   0,
					BurstDuration: 5,
					Priority:      2,
				},
				{
					ProcessID:     2,
					ArrivalTime:   3,
					BurstDuration: 9,
					Priority:      1,
				},
				{
					ProcessID:     3,
					ArrivalTime:   3,
					BurstDuration: 6,
					Priority:      3,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := loadProcesses(tt.args.r)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadProcesses() = %v, want %v", got, tt.want)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func loadFixture(t *testing.T, p ...string) string {
	b, err := os.ReadFile(path.Join(p...))
	if err != nil {
		t.Fail()
	}

	return string(b)
}
