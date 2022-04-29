package main 


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