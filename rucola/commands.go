package main

import (
	"fmt"
)

func command(origin string, name string, arguments []string) {
	switch name {
	case "?", "h", "help", "wtf":
		showHelp(origin)
	case "ping":
		pingClient(origin, arguments)
	case "msg", "message":
		sendDirectMsg(origin, arguments)
	case "list", "ls":
		listClients(origin)
	}
}

func showHelp(origin string) {
	client, _ := findClient(origin)

	msg :=
		`------------
write a message to send it to all other users

commands:

/ls - list users
/msg [user] [message] - send a private message
/ping [user] - ping a user
------------
`

	client.conn.Write([]byte(msg))
}

func listClients(origin string) {
	client, _ := findClient(origin)

	var output string

	for i := 0; i < len(clients); i++ {
		output = output + " * " + clients[i].name + " " + clients[i].conn.RemoteAddr().String() + "\n"
	}

	client.conn.Write([]byte(output))
}

func pingClient(origin string, arguments []string) {
	targetClient, _ := findClient(arguments[0])

	writeClient(targetClient, origin+" pinged you")
}

func sendDirectMsg(origin string, args []string) {
	target := args[0]
	var message string

	for i := 1; i < len(args); i++ {
		message = message + " " + args[i]
	}

	originClient, _ := findClient(origin)
	client, _ := findClient(target)

	if client.id == "" {
		fmt.Printf("private message from %s failed: %s does not exist\n", origin, target)

		if originClient.id != "" {
			writeClient(originClient, "user does not exist: "+target)
		} else {
			fmt.Printf("Sender of message does not exist\n")
		}
	} else {
		writeClient(client, "["+origin+"]: "+message)
	}
}
