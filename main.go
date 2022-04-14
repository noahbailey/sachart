package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
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
	Network   Network
}

type Timestamp struct {
	Date     string
	Time     string
	Utc      int
	Interval int
}

type Cpu struct {
	Cpu    string
	User   float64
	Nice   float64
	System float64
	Iowait float64
	Steal  float64
	Idle   float64
}

type Memory struct {
	Memfree    int
	Avail      int
	Memused    int
	MemusedPct float64 `json:"memused-percent"`
	Buffers    int
	Cached     int
	Commit     int
	CommitPct  float64 `json:"commit-percent"`
	Active     int
	Inactive   int
	Dirty      int
}

type Queue struct {
	RunqSz  int     `json:"runq-sz"`
	PlistSz int     `json:"plist-sz"`
	Load1   float64 `json:"ldavg-1"`
	Load5   float64 `json:"ldavg-5"`
	Load15  float64 `json:"ldavg-15"`
	Blocked int
}

type Network struct {
	NetDev []NetDev `json:"net-dev"`
}

type NetDev struct {
	Iface   string
	Rxpck   float64
	Txpck   float64
	Rxkb    float64
	Txkb    float64
	Rxcmp   float64
	Txcmp   float64
	Rxmcst  float64
	UtilPct float64 `json:"ifutil-percent"`
}

func main() {
	//Command line flags:
	flagCpu := flag.Bool("cpu", true, "Show CPU/Memory/Load graph")
	flagNet := flag.Bool("net", false, "Show Network/IO graph")
	flagDays := flag.Int("days", 0, "Show data from previous days")
	flag.Parse()

	//Get data from sadf
	lines := getFile(*flagDays)

	//parse the JSON object
	system := parseJson(lines)

	//Output the formatted chart
	if *flagNet == true {
		drawNetChart(system)
	} else if *flagCpu == true {
		drawCpuMemLoadChart(system)
	}
}

func getFile(pastDays int) string {
	strPastDays := "-" + strconv.Itoa(pastDays)
	// Get CPU&Memory stats in JSON format:
	cmd := exec.Command("sadf", "-j", "--", "-r", "-u", "-q", "-n", "DEV", strPastDays)
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func numCores() int {
	cores := runtime.NumCPU()
	return cores
}

func parseJson(rawJson string) (system System) {
	json.Unmarshal([]byte(rawJson), &system)
	return
}

func parseDate(timeUtc string) string {
	t, err := time.ParseInLocation("15:04:05", timeUtc, time.UTC)
	if err != nil {
		log.Panic(err)
	}
	_, offset := time.Now().Zone()
	offsetDuration := time.Duration(offset) * time.Second
	localTime := t.Add(offsetDuration).Format("15:04:05")
	return localTime
}

//Draw a chart for CPU graph
func drawCpuMemLoadChart(sys System) {
	fmt.Println("TIME     | CPU                      | MEMORY                   | LOAD AVG")
	for i, val := range sys.Sysstat.Hosts[0].Statistics {
		// The first datapoint from "midnight" can contain strange or incorrect data, safer to skip it
		if i == 0 {
			continue
		}
		//Convert timestamp to localtime
		localTime := parseDate(val.Timestamp.Time)
		//draw the bar for all stats in this set
		cpuMemLoadBar := drawCpuMemLoadBar(val)
		fmt.Println(localTime + cpuMemLoadBar)
	}
}

// Draws out the individual row for the chart
// Could/should be cleaned up a little...
func drawCpuMemLoadBar(val Statistics) string {
	cpuSys := int(val.Cpu[0].System) / 4
	cpuUsr := int(val.Cpu[0].User) / 4
	cpuSpace := 25 - cpuSys - cpuUsr
	memUsd := int(val.Memory.MemusedPct) / 4
	memSpace := 25 - memUsd
	loadPerCore := val.Queue.Load5 / float64(numCores())
	load5 := int(loadPerCore * 10)
	barCpuSys := strings.Repeat("@", cpuSys)
	barCpuUsr := strings.Repeat("#", cpuUsr)
	barCpuSpace := strings.Repeat(" ", cpuSpace)
	barMem := strings.Repeat("*", memUsd)
	barMemSpace := strings.Repeat(" ", memSpace)
	barLoad5 := strings.Repeat("|", load5)
	output := " |\033[31m" + barCpuSys + "\033[32m" + barCpuUsr + "\033[0m" +
		barCpuSpace + " |" + "\033[33m" + barMem + "\033[0m" + barMemSpace + " |\033[34m" + barLoad5 + "\033[0m"
	return output
}

//Draw the network/Io row
func drawNetBar(val Statistics, highestTx float64, highestRx float64) string {
	//Determine the total throughput on all interfaces...
	var totalTx float64
	var totalRx float64
	for _, iface := range val.Network.NetDev {
		totalTx += iface.Txkb
		totalRx += iface.Rxkb
	}
	//Express the current value as a percent of the highest value:
	pctRx := int(totalRx / highestRx * 25)
	pctTx := int(totalTx / highestTx * 25)
	barRx := strings.Repeat("=", (pctRx))
	spacesRx := strings.Repeat(" ", (25 - pctRx))
	barTx := strings.Repeat("=", pctTx)
	spacesTx := strings.Repeat(" ", (25 - pctTx))
	barRq := strings.Repeat("-", val.Queue.RunqSz)
	barBk := strings.Repeat(">", val.Queue.Blocked)
	output := " |\033[34m" + barRx + spacesRx + " \033[0m|\033[35m" + barTx + spacesTx + " \033[0m|\033[36m" + barRq + "\033[31m" + barBk + "\033[0m"
	return output
}

func drawNetChart(sys System) {
	//Determine the "peak" values first:
	highestTx, highestRx := getHighestNetThroughput(sys)

	//Show a header with net throughput info:
	fmt.Println("Max TX (Kb/s): ", highestTx, " Max RX (Kb/s): ", highestRx)
	fmt.Println("TIME     | DOWNLOAD                 | UPLOAD                   | IO (RunQ + Blocked)")

	for i, val := range sys.Sysstat.Hosts[0].Statistics {
		// The first datapoint from "midnight" can contain strange or incorrect data, safer to skip it
		if i == 0 {
			continue
		}
		localTime := parseDate(val.Timestamp.Time)
		netBar := drawNetBar(val, highestTx, highestRx)
		fmt.Println(localTime + netBar)
	}
}

//For each time bucket, calculate the total throughput on all interfaces
//	The highestTx/Rx variables should be the same as the highest throughput seen
func getHighestNetThroughput(sys System) (highestTx float64, highestRx float64) {
	for _, val := range sys.Sysstat.Hosts[0].Statistics {
		var totalTx float64
		var totalRx float64
		for _, iface := range val.Network.NetDev {
			totalTx += iface.Txkb
			totalRx += iface.Rxkb
		}
		if totalTx > highestTx {
			highestTx = totalTx
		}
		if totalRx > highestRx {
			highestRx = totalRx
		}
	}
	return
}
