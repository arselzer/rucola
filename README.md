Chat Server in Go
======

This is a very minimal TCP chat server, which can be used with netcat.

## Requirements

Go 1.2 or higher (`go version`)
[goleveldb](https://github.com/syndtr/goleveldb) (`go get github.com/syndtr/goleveldb`)

## Obtaining

to build:
`git clone https://github.com/AlexanderSelzer/goleveldb`

to install:
*warning: might not yet work, I am still new to Go package organization :)*
`go get github.com/AlexanderSelzer/rucola`
`go install github.com/AlexanderSelzer/rucola`

## Building

### Build
`make build`
### Run
`make run`
### Test (very simple, just connects, and tells the server a name)
`make test`

## Usage

```
  -db="./db": the LevelDB directory
  -port="8001": the tcp port to listen on
```

Use netcat (`nc`) to connect to the server, but not telnet.
Telnet will screw up the input because it has its own kind of transfer format (less raw).

### Example

```bash
$nc localhost 8001
-- select a name:
alex
  your id:
  69a2c22213bb2e8b
  remember it to log in again.
/ls
 * alex 127.0.0.1:55431
/help
------------
commands:

/ls - list users
/msg [user] [message] - send a private message
/ping [user] - ping a user
------------
-- koen joined the chat
-- [koen]: hello
/msg koen hi koen
```

```bash
$nc localhost 8001
-- select a name:
koen
  your id:
  84724bd8e0f29095
  remember it to log in again.
hello
-- [alex]: hi koen
```
