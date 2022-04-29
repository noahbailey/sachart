package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

//Just returns the number of cores/threads on the system, that is all.
func numCores() int {
	cores := runtime.NumCPU()
	return cores
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
