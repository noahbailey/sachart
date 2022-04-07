package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
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
	flag.Parse()

	//Get data from sadf
	lines := getFile()

	//parse the JSON object
	system := parseJson(lines)

	//Output the formatted chart
	if *flagNet == true {
		drawNetChart(system)
	} else if *flagCpu == true {
		drawCpuMemLoadChart(system)
	}
}

func getFile() string {
	// Get CPU&Memory stats in JSON format:
	cmd := exec.Command("sadf", "-j", "--", "-r", "-u", "-q", "-n", "DEV")
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
func drawCpuMemLoadChart(sys System) {

	fmt.Println("TIME     | CPU                      | MEMORY                   | LOAD AVG")

	for i, val := range sys.Sysstat.Hosts[0].Statistics {

		// The first datapoint from "midnight" can contain strange or incorrect data, safer to skip it
		if i == 0 {
			continue
		}

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

		fmt.Println(val.Timestamp.Time + " |\033[31m" + barCpuSys + "\033[32m" + barCpuUsr + "\033[0m" +
			barCpuSpace + " |" + "\033[33m" + barMem + "\033[0m" + barMemSpace + " |\033[34m" + barLoad5 + "\033[0m")
	}
}

func drawNetChart(sys System) {
	//Determine the "peak" values first:
	highestTx := 0.0
	highestRx := 0.0
	for _, val := range sys.Sysstat.Hosts[0].Statistics {
		for _, iface := range val.Network.NetDev {
			if iface.Txkb > highestTx {
				highestTx = iface.Txkb
			}
			if iface.Rxkb > highestRx {
				highestRx = iface.Rxkb
			}
		}
	}

	fmt.Println("Max TX (Kb/s): ", highestTx, " Max RX (Kb/s): ", highestRx)

	fmt.Println("TIME     | DOWNLOAD                 | UPLOAD                   | IO (RunQ + Blocked)")

	for i, val := range sys.Sysstat.Hosts[0].Statistics {
		// The first datapoint from "midnight" can contain strange or incorrect data, safer to skip it
		if i == 0 {
			continue
		}
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

		fmt.Println(val.Timestamp.Time + " |\033[34m" + barRx + spacesRx + " \033[0m|\033[35m" + barTx + spacesTx + " \033[0m|\033[36m" + barRq + "\033[31m" + barBk + "\033[0m")
	}
}

func numCores() int {
	cores := runtime.NumCPU()
	return cores
}
