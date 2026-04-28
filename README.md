# Tomato

An overengineered Pomodoro CLI

## Description

Tomato is an CLI to track time using [Pomodoro Technique](https://en.wikipedia.org/wiki/Pomodoro_Technique) it work by save current tomato phase(pomodoro, short break, long break) on a [BadgerDB](https://github.com/dgraph-io/badger) store

## Getting Started

### Dependencies

* Currently only work on linux with with ```paplay``` installed if you want its notifications to have sound

### Installing

Build from source(having go 1.26.2+ installed):
```
git clone https://github.com/zimlewis/tomato
cd tomato
go build -o tomato .
```

Install with go(having go installed)
```
go install github.com/zimlewis/tomato
```
Go will installed it in ```$GOPATH/bin```


### Executing program

Because the cli work by writing to BadgerDB, it is imposible to make it to run by itself without some sort of communication between each call, so, coming from a networking background, I moved to gRPC pattern.

The server will hold an connection to BadgerDB database(saved in /tmp/tomato/, might add a way to read config file in later version).

So first thing first is to start the server:
```
tomato ss
```
(ss stand for Start Server)

Client(other cli command) will request send request to the server

use
```
tomato help [COMMAND]
```
for more information

To start the session with current tomato phase:
```
tomato start
```

To stop the session:
```
tomato stop
```

To switch different tomato phase(this will automatically stop your session)
```
tomato switch [up|down]
```

To change to a specified tomato phase(this will automatically stop your session)
```
tomato switch [pomodoro|short|long]
```

To open the time tracker:
```
tomato current
```
Flags:
| option           | typeof  | default | description |
| ---------------- | ------- | ------- | ----------- |
| `-f`, `--format` | string  |`default`|The formatter that is used to print current time, currently accept 3 values: `default`, `waybar`, `basic`. Any other value will be conver to default, default is human readable format: eg. "Your pomodoro session has 05:26 remaining"| 

## Help

Run tomato help for more information
```
tomato help
```
