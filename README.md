# SaChart

This tool's objective is to answer the question, "Was the CPU load high last night when Foobar broke?" in a simple way that does not require complex monitoring systems. In conjunction with the lightweight `sysstat` daemon, this tool makes easy to read histograms in your shell. 

There are two modes, CPU and NET mode. 

CPU mode shows the data you're used to seeing in tools like `htop`, processor usage, memory usage, and per-core load average. NET mode displays total network throughput on all interfaces as a proportion of the highest measured value, and the system's run queue and blocked threads. 

Output: 

`./sachart -cpu`

```
TIME     | CPU                      | MEMORY                   | LOAD AVG
13:10:01 |@######                   |****                      ||||
13:15:01 |@######                   |****                      |||||
13:20:01 |@#######                  |****                      |||||
13:25:01 |@#######                  |****                      |||||
13:30:01 |@#######                  |****                      ||||||
13:35:01 |##                        |****                      ||||
13:40:01 |                          |****                      ||
13:45:01 |                          |****                      |
```

`./sachart -net`

```
TIME     | DOWNLOAD                 | UPLOAD                   | IO (RunQ + Blocked)
13:10:01 |===================       |========================  |
13:15:01 |=================         |========================  |--------------
13:20:01 |=====================     |========================= |--------
13:25:01 |========================= |========================  |-
13:30:01 |=======================   |========================  |---
13:35:01 |=====                     |=====                     |--
13:40:01 |                          |                          |
13:45:01 |                          |                          |
```

Previous days data can also be viewed by using the `-days` flag. For example, to see yesterday's CPU graph:

`./sachart -cpu -days 1`

```
TIME     | CPU                      | MEMORY                   | LOAD AVG
13:20:02 |                          |*                         |
13:25:01 |#                         |****                      |
13:30:01 |#                         |****                      |
13:35:01 |#                         |****                      ||
13:40:01 |@##                       |****                      |||
```

## Sysstat setup

Make sure sysstat is installed: 

    sudo apt install sysstat

The service should be running: 

    sudo systemctl enable --now sysstat.service

And data collection should be enabled: 

    sudo sed -i 's/false/true/g' /etc/default/sysstat

