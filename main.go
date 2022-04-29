package main

import "flag"

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
