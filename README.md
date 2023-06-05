# GoBatteryWatcher

GoBatteryWatcher is a command-line interface (CLI) application for Linux built in Go that allows you to monitor your device's battery usage.

## Features

- Live tracking of battery usage for all devices or processes
- Graphical display of battery usage over time
- Display of top 10 devices by power consumption
- Easy-to-use, interactive CLI

## Prerequisites

- `powertop` must be installed on your system

## Development Status

**This application is still in development.** Although the graph and the top 10 list are functional. The sum of the power consumption of all devices and processes may not be accurate when I compare it directly to powertop so I am still working on that. And I'll add features as I go along.

## Todo

- [ ] Use the history data to calculate the average watt/hour consumption of each device/process
- [ ] Use the package battery to get the battery information and do something with it
- [ ] Add a feature to display the battery information
- [ ] Separate the agent and the CLI so it runs in the background (as service) and the CLI just queries the agent
- [ ] More useful stats
