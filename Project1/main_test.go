package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/google/go-cmp/cmp"
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
						ProcessID:     "P0",
						ArrivalTime:   0,
						BurstDuration: 5,
						Priority:      2,
					},
					{
						ProcessID:     "P1",
						ArrivalTime:   3,
						BurstDuration: 9,
						Priority:      1,
					},
					{
						ProcessID:     "P2",
						ArrivalTime:   6,
						BurstDuration: 6,
						Priority:      3,
					},
				},
				title: "First-come, first-serve",
			},
			wantOut: loadFixture(t, "fcfs_fixture.txt"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var w bytes.Buffer
			FCFSSchedule(&w, tt.args.title, tt.args.processes)
			if diff := cmp.Diff(w.String(), tt.wantOut); diff != "" {
				t.Errorf(diff)
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
				r: strings.NewReader(`ProcessID,Burst Duration,Arrival Time,Priority
P0,5,0,2
P1,9,3,1
P2,6,3,3`),
			},
			want: []Process{
				{
					ProcessID:     "P0",
					ArrivalTime:   0,
					BurstDuration: 5,
					Priority:      2,
				},
				{
					ProcessID:     "P1",
					ArrivalTime:   3,
					BurstDuration: 9,
					Priority:      1,
				},
				{
					ProcessID:     "P2",
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
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf(diff)
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

func Test_openProcessingFile1(t *testing.T) {
	tmpFile, tErr := os.CreateTemp(t.TempDir(), "")
	if tErr != nil {
		t.Fatal(tErr)
	}

	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    *os.File
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				args: []string{"binary_name", tmpFile.Name()},
			},
			want: tmpFile,
		},
		{
			name: "not enough args",
			args: args{
				args: []string{"binary_name"},
			},
			wantErr: true,
		},
		{
			name: "bad file",
			args: args{
				args: []string{"binary_name", "bad_file_name"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, closeFn, err := openProcessingFile(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("openProcessingFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			if got == nil {
				t.Fatal("file is unexpectedly nil")
			}
			if closeFn == nil {
				t.Fatal("closeFn is unexpectedly nil")
			}
			t.Cleanup(closeFn)

			f1, err := os.Stat(got.Name())
			if err != nil {
				t.Fatalf("Could not stat file: %v", got)
			}
			f2, err := os.Stat(tt.want.Name())
			if err != nil {
				t.Fatalf("Could not stat file: %v", tt.want)
			}

			if !os.SameFile(f1, f2) {
				t.Fatal("files are not the same")
			}
		})
	}
}

func Test_outputGantt(t *testing.T) {
	t.Parallel()
	type args struct {
		gantt []TimeSlice
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{
			name: "consecutive processes",
			args: args{
				gantt: []TimeSlice{
					{PID: "A", Start: 1, Stop: 2},
					{PID: "B", Start: 2, Stop: 4},
					{PID: "C", Start: 4, Stop: 7},
					{PID: "D", Start: 7, Stop: 11},
					{PID: "E", Start: 11, Stop: 16},
				},
			},
			wantW: `Gantt schedule
|  A  |  B  |  C  |  D  |  E  |
1     2     4     7     11    16

`,
		},
		{
			name: "nonconsecutive processes",
			args: args{
				gantt: []TimeSlice{
					{PID: "A", Start: 1, Stop: 2},
					{PID: "B", Start: 5, Stop: 6},
					{PID: "C", Start: 6, Stop: 7},
					{PID: "D", Start: 9, Stop: 11},
					{PID: "E", Start: 13, Stop: 16},
				},
			},
			wantW: `Gantt schedule
|  A  |  -  |  B  |  C  |  -  |  D  |  -  |  E  |
1     2     5     6     7     9     11    13    16

`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := &bytes.Buffer{}
			outputGantt(w, tt.args.gantt)
			if diff := cmp.Diff(tt.wantW, w.String()); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
