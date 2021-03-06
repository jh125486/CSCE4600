package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func main() {
	processes := loadProcesses(os.Args)

	FCFSSchedule(os.Stdout, processes, "First-come, first-serve")

	SJFSchedule(os.Stdout, processes, "Shortest-job-first")

	SJFPrioritySchedule(os.Stdout, processes, "Priority")

	RRSchedule(os.Stdout, processes, "Round-robin")
}

type (
	Process struct {
		ProcessID     int64
		ArrivalTime   int64
		BurstDuration int64
		Priority      int64
	}
	TimeSlice struct {
		PID   int64
		Start int64
		Stop  int64
	}
)

// FCFSSchedule example output
// ----------------------------------------------
// First-come, first-serve
// ----------------------------------------------
// Gantt schedule
// |   1   |   2   |   3   |
// 0       5       14      20
//
// Schedule table
// +----+----------+-------+---------+---------+------------+------------+
// | ID | PRIORITY | BURST | ARRIVAL |  WAIT   | TURNAROUND |    EXIT    |
// +----+----------+-------+---------+---------+------------+------------+
// |  1 |        2 |     5 |       0 |       0 |          5 |          5 |
// |  2 |        1 |     9 |       3 |       2 |         11 |         14 |
// |  3 |        3 |     6 |       6 |       8 |         14 |         20 |
// +----+----------+-------+---------+---------+------------+------------+
// |                                   AVERAGE |  AVERAGE   | THROUGHPUT |
// |                                    3.33   |   10.00    |   0.15/T   |
// +----+----------+-------+---------+---------+------------+------------+
func FCFSSchedule(w io.Writer, processes []Process, title string) {
	var (
		serviceTime     int64
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		waitingTime     int64
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)
	for i := range processes {
		if processes[i].ArrivalTime > 0 {
			waitingTime = serviceTime - processes[i].ArrivalTime
		}
		totalWait += float64(waitingTime)

		start := waitingTime + processes[i].ArrivalTime

		turnaround := processes[i].BurstDuration + waitingTime
		totalTurnaround += float64(turnaround)

		completion := processes[i].BurstDuration + processes[i].ArrivalTime + waitingTime
		lastCompletion = float64(completion)

		schedule[i] = []string{
			fmt.Sprint(processes[i].ProcessID),
			fmt.Sprint(processes[i].Priority),
			fmt.Sprint(processes[i].BurstDuration),
			fmt.Sprint(processes[i].ArrivalTime),
			fmt.Sprint(waitingTime),
			fmt.Sprint(turnaround),
			fmt.Sprint(completion),
		}
		serviceTime += processes[i].BurstDuration

		gantt = append(gantt, TimeSlice{
			PID:   processes[i].ProcessID,
			Start: start,
			Stop:  serviceTime,
		})
	}

	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

func SJFPrioritySchedule(w io.Writer, processes []Process, title string) {

}

func SJFSchedule(w io.Writer, processes []Process, title string) {

}

func RRSchedule(w io.Writer, processes []Process, title string) {

}

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
	_, _ = fmt.Fprintln(w, "\n")
}

func outputSchedule(w io.Writer, rows [][]string, wait, turnaround, throughput float64) {
	_, _ = fmt.Fprintln(w, "Schedule table")
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"ID", "Priority", "Burst", "Arrival", "Wait", "Turnaround", "Exit"})
	table.AppendBulk(rows)
	table.SetFooter([]string{"", "", "", "",
		fmt.Sprintf("Average\n%.2f", wait),
		fmt.Sprintf("Average\n%.2f", turnaround),
		fmt.Sprintf("Throughput\n%.2f/t", throughput)})
	table.Render()
}

func loadProcesses(args []string) []Process {
	if len(args) != 2 {
		exitWithMsg("Must give a scheduling file")
	}

	f, err := os.Open(args[1])
	if err != nil {
		exitWithMsg("Must give a scheduling file")
	}

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		exitWithMsg(err.Error())
	}

	processes := make([]Process, len(rows))
	for i := range rows {
		processes[i].ProcessID = mustStrToInt(rows[i][0])
		processes[i].BurstDuration = mustStrToInt(rows[i][1])
		processes[i].ArrivalTime = mustStrToInt(rows[i][2])
		if len(rows[i]) == 4 {
			processes[i].Priority = mustStrToInt(rows[i][3])
		}
	}

	return processes
}

func exitWithMsg(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func mustStrToInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		exitWithMsg(err.Error())
	}

	return i
}
