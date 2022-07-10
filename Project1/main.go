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

type Process struct {
	ProcessID     int64
	ArrivalTime   int64
	BurstDuration int64
	Priority      int64
}

// FCFSSchedule example output
// ----------------------------------------------
//            First-come, first-serve
// ----------------------------------------------
// Gantt schedule
// +-----+---+---+----+
// | PID | 1 | 2 |  3 |
// | T   | 0 | 5 | 14 |
// +-----+---+---+----+
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
		gantt           = make([][]string, len(processes))
	)
	for i := range processes {
		if processes[i].ArrivalTime > 0 {
			waitingTime = serviceTime - processes[i].ArrivalTime
		}
		totalWait += float64(waitingTime)

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
		gantt[i] = []string{
			fmt.Sprint(processes[i].ProcessID),
			fmt.Sprint(serviceTime),
		}

		serviceTime += processes[i].BurstDuration
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

func outputGantt(w io.Writer, gantt [][]string) {
	_, _ = fmt.Fprintln(w, "Gantt schedule")
	table := tablewriter.NewWriter(w)
	data := make([][]string, 2)
	data[0] = []string{"PID"}
	data[1] = []string{"T"}
	for i := range gantt {
		data[0] = append(data[0], gantt[i][0])
		data[1] = append(data[1], gantt[i][1])
	}
	table.AppendBulk(data)
	table.Render()
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
