# SaChart

A useful tool for any system with sysstat data collection. Sachart visualizes your CPU, Memory, and 5 Minute load average data in an easy to read vertical histogram. 

Output: 

```
TIME     | CPU                      | MEMORY                   | LOAD AVG
13:10:01 |#                         |****                      | ||||||||
13:15:01 |#                         |*****                     | ||||||||||
13:20:02 |                          |*                         | |||||
13:25:01 |#                         |****                      | ||||||
13:30:01 |#                         |****                      | |||||||
13:35:01 |#                         |****                      | ||||||||||
13:40:01 |@##                       |****                      | ||||||||||||||||
13:45:01 |@##                       |****                      | ||||||||||||||
13:50:01 |@##                       |****                      | ||||||||||||||
13:55:01 |#                         |****                      | ||||||||||
```

## Sysstat setup

Make sure sysstat is installed: 

    sudo apt install sysstat

The service should be running: 

    sudo systemctl enable --now sysstat.service

And data collection should be enabled: 

    sudo sed -i 's/false/true/g' /etc/default/sysstat

