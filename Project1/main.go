package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func main() {
	// parse args.
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	scheduler, data, err := parseCLI(flagSet, os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stdout, err)
		flagSet.PrintDefaults()
		os.Exit(1)
	}

	// Load and parse processes.
	processes, err := loadProcesses(data)
	if err != nil {
		log.Fatal(err)
	}

	// Run the given scheduler.
	switch scheduler {
	case fcfs:
		FCFSSchedule(os.Stdout, "First-come, first-serve", processes)
	case sjf:
		SJFSchedule(os.Stdout, "Shortest-job-first", processes)
	case sjfp:
		SJFPrioritySchedule(os.Stdout, "Priority", processes)
	case rr:
		RRSchedule(os.Stdout, "Round-robin", processes)
	}
}

//go:generate stringer -type=Scheduler
type Scheduler uint

const (
	fcfs Scheduler = iota + 1
	sjf
	sjfp
	rr
)

func parseCLI(flagSet *flag.FlagSet, args []string) (cmd Scheduler, data io.Reader, err error) {
	fcfsFlag := flagSet.Bool(fcfs.String(), false, "First-come, first-serve scheduling")
	sjfFlag := flagSet.Bool(sjf.String(), false, "Shortest-job-first scheduling")
	sjfpFlag := flagSet.Bool(sjfp.String(), false, "Shortest-job-first with priority scheduling")
	rrFlag := flagSet.Bool(rr.String(), false, "Round-robin scheduling")
	if err := flagSet.Parse(args); err != nil {
		return 0, nil, err
	}
	// validate only one flag is set
	var count int
	if *fcfsFlag {
		count++
		cmd = fcfs
	}
	if *sjfFlag {
		count++
		cmd = sjf
	}
	if *sjfpFlag {
		count++
		cmd = sjfp
	}
	if *rrFlag {
		count++
		cmd = rr
	}
	switch count {
	case 0:
		return 0, nil, fmt.Errorf("one scheduler flag must be set")
	case 1:
		// validate that data file is piped in.
		if data, err := readData(os.Args[:2]); err != nil {
			return 0, nil, err
		} else {
			return cmd, data, nil
		}
	default:
		return 0, nil, fmt.Errorf("only one scheduler flag must be set")
	}
}

func readData(args []string) (io.Reader, error) {
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return os.Stdin, nil
	} else if len(args) == 1 {
		return nil, fmt.Errorf("scheduler data must be passed in or file given as last argument")
	}
	r, err := os.Open(os.Args[2])
	if err != nil {
		return nil, fmt.Errorf("%w: error opening data file", err)
	}

	return r, nil
}

func openProcessingFile(args ...string) (*os.File, func(), error) {
	if len(args) != 2 {
		return nil, nil, fmt.Errorf("%w: must give a scheduling file to process", ErrInvalidArgs)
	}
	// Read in CSV process CSV file
	f, err := os.Open(args[1])
	if err != nil {
		return nil, nil, fmt.Errorf("%v: error opening scheduling file", err)
	}
	closeFn := func() {
		if err := f.Close(); err != nil {
			log.Fatalf("%v: error closing scheduling file", err)
		}
	}

	return f, closeFn, nil
}

//region Output helpers

func outputTitle(w io.Writer, title string) {
	_, _ = fmt.Fprintln(w, strings.Repeat("-", len(title)*2))
	_, _ = fmt.Fprintln(w, strings.Repeat(" ", len(title)/2), title)
	_, _ = fmt.Fprintln(w, strings.Repeat("-", len(title)*2))
}

func outputGantt(w io.Writer, gantt []TimeSlice) {
	_, _ = fmt.Fprintln(w, "Gantt schedule")
	_, _ = fmt.Fprint(w, "|")
	for i := range gantt {
		pid := fmt.Sprint(gantt[i].PID)
		padding := strings.Repeat(" ", (8-len(pid))/2)
		_, _ = fmt.Fprint(w, padding, pid, padding, "|")
	}
	_, _ = fmt.Fprintln(w)
	for i := range gantt {
		_, _ = fmt.Fprint(w, fmt.Sprint(gantt[i].Start), "\t")
		if len(gantt)-1 == i {
			_, _ = fmt.Fprint(w, fmt.Sprint(gantt[i].Stop))
		}
	}
	_, _ = fmt.Fprintf(w, "\n\n")
}

func outputSchedule(w io.Writer, rows [][]string, wait, turnaround, throughput float64) {
	_, _ = fmt.Fprintln(w, "Schedule table")
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"ID", "Priority", "Burst", "Arrival", "Wait", "Turnaround", "Exit"})
	table.AppendBulk(rows)
	table.Render()
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintf(w, "Average wait: %.2f\n", wait)
	_, _ = fmt.Fprintf(w, "Average turnaround: %.2f\n", turnaround)
	_, _ = fmt.Fprintf(w, "Throughput: %.2f\n", throughput)
}

//endregion

//region Loading processes.

var ErrInvalidArgs = errors.New("invalid args")

func loadProcesses(r io.Reader) ([]Process, error) {
	rows, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w: reading CSV", err)
	}
	rows = rows[1:] // skip header row
	processes := make([]Process, len(rows))
	for i := range rows {
		processes[i].ProcessID = mustStrToInt(rows[i][0])
		processes[i].BurstDuration = mustStrToInt(rows[i][1])
		processes[i].ArrivalTime = mustStrToInt(rows[i][2])
		if len(rows[i]) == 4 {
			processes[i].Priority = mustStrToInt(rows[i][3])
		}
	}

	return processes, nil
}

func mustStrToInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	return i
}

//endregion
