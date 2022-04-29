package main 

import (
	"fmt"
	"strings"
)


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
