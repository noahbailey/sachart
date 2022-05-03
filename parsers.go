package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

//Captures the output of the `sadf` command
// Essentially, just gets data from /var/log/sysstat/saXX
// and exports it in a machine-readable format...
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

//Parses the JSON from getFile() into a System object
func parseJson(rawJson string) (system System) {
	json.Unmarshal([]byte(rawJson), &system)
	return
}

//Converts the UTC formatted time to local time for readability
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
