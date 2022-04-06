package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type System struct {
	Sysstat Sysstat
}

type Sysstat struct {
	Hosts []Host
}

type Host struct {
	Nodename     string
	Sysname      string
	Release      string
	Machine      string
	NumberOfCpus int    `json:"number-of-cpus"`
	FileDate     string `json:"file-date"`
	FileUtcTime  string `json:"file-utc-time"`
	Timezone     string
	Statistics   []Statistics
}

type Statistics struct {
	Timestamp Timestamp
	Cpu       []Cpu `json:"cpu-load"`
	Memory    Memory
	Queue     Queue
}

type Timestamp struct {
	Date     string
	Time     string
	Utc      int
	Interval int
}

type Cpu struct {
	Cpu    string
	User   float32
	Nice   float32
	System float32
	Iowait float32
	Steal  float32
	Idle   float32
}

type Memory struct {
	Memfree    int
	Avail      int
	Memused    int
	MemusedPct float32 `json:"memused-percent"`
	Buffers    int
	Cached     int
	Commit     int
	CommitPct  float32 `json:"commit-percent"`
	Active     int
	Inactive   int
	Dirty      int
}

type Queue struct {
	RunqSz  int     `json:"runq-sz"`
	PlistSz int     `json:"plist-sz"`
	Load1   float32 `json:"ldavg-1"`
	Load5   float32 `json:"ldavg-5"`
	Load15  float32 `json:"ldavg-15"`
	Blocked int
}

func main() {
	lines := getFile()

	//parse the JSON object
	system := parseJson(lines)

	//Output the formatted chart
	drawChart(system)
}

func getFile() string {
	// Get CPU&Memory stats in JSON format:
	cmd := exec.Command("sadf", "-j", "--", "-r", "-u", "-q")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func parseJson(rawJson string) (system System) {
	json.Unmarshal([]byte(rawJson), &system)
	return
}

//Draw a chart for CPU graph
func drawChart(sys System) {

	fmt.Println("TIME     | CPU                      | MEMORY                   | LOAD AVG")

	for _, val := range sys.Sysstat.Hosts[0].Statistics {
		cpuSys := int(val.Cpu[0].System) / 4
		cpuUsr := int(val.Cpu[0].User) / 4
		cpuSpace := 25 - cpuSys - cpuUsr

		memUsd := int(val.Memory.MemusedPct) / 4
		memSpace := 25 - memUsd

		load5 := int(val.Queue.Load5 * 10)

		barCpuSys := strings.Repeat("@", cpuSys)
		barCpuUsr := strings.Repeat("#", cpuUsr)
		barCpuSpace := strings.Repeat(" ", cpuSpace)
		barMem := strings.Repeat("*", memUsd)
		barMemSpace := strings.Repeat(" ", memSpace)

		barLoad5 := strings.Repeat("|", load5)

		fmt.Println(val.Timestamp.Time + " |\033[31m" + barCpuSys + "\033[32m" + barCpuUsr + "\033[0m" +
			barCpuSpace + " |" + "\033[33m" + barMem + "\033[0m" + barMemSpace + " | \033[34m" + barLoad5 + "\033[0m")
	}
}
