# Tomato
A Pomodoro CLI

## Description
Tomato is an CLI to track time using [Pomodoro Technique](https://en.wikipedia.org/wiki/Pomodoro_Technique) it work by save current tomato phase(pomodoro, short break, long break) on a [BadgerDB](https://github.com/dgraph-io/badger) store

### Motivation
I made this to use in my desktop environment where I can manage time using pomodoro technique in multiple application and devices

## Getting Started

### Installing
Make sure you have go installed

Build from source:
```bash
git clone https://github.com/zimlewis/tomato
cd tomato
go build -o tomato .
```

Install with go
```bash
go install github.com/zimlewis/tomato
```
Go will installed it in ```$GOPATH/bin```

### Quick Start
Because the cli work by writing to BadgerDB, it is imposible to make it to run by itself without some sort of communication between each call, so, coming from a networking background, I moved to gRPC pattern.

The server will hold an connection to BadgerDB database(saved in /tmp/tomato/, might add a way to read config file in later version).

So first thing first is to start the server:
```bash
tomato ss
```
(ss stand for Start Server)

### Usage 
Client(other cli command) will request send request to the server

use
```bash
tomato help [COMMAND]
```
for more information

To start the session with current tomato phase:
```bash
tomato start
```

To stop the session:
```bash
tomato stop
```

To switch different tomato phase(this will automatically stop your session)
```bash
tomato switch [up|down]
```

To change to a specified tomato phase(this will automatically stop your session)
```bash
tomato change [pomodoro|short|long]
```

To open the time tracker:
```bash
tomato current
```
Flags:
| option           | typeof  | default | description |
| ---------------- | ------- | ------- | ----------- |
| `-f`, `--format` | string  |`default`|The formatter that is used to print current time, currently accept 3 values: `default`, `waybar`, `basic`. Any other value will be conver to default, default is human readable format: eg. "Your pomodoro session has 05:26 remaining"| 

## Help
Run tomato help for more information
```bash
tomato help
```

## Contributing

### Clone the repo
```bash
git clone https://github.com/zimlewis/tomato
cd tomato
```

### Build the compiled binary
```bash
go build
```

### Submit a pull request
If you'd like to contribute, please fork the repository and open a pull request to the `dev` branch.
