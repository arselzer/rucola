package main

import (
	"fmt"
)

var clients []Client

func findClient(name string) (Client, int) {
	for i := 0; i < len(clients); i++ {
		if clients[i].name == name {
			return clients[i], i
		}
	}
	return Client{id: "", name: "", online: false}, 0
}

func addClient(client Client) {
	clients = append(clients, client)
}

func removeClient(name string) {
	_, i := findClient(name)

	clients = append(clients[:i], clients[1+i:]...)
}

func writeClient(client Client, msg string) {
	_, err := client.conn.Write([]byte("  " + msg + "\n"))

	if err != nil {
		fmt.Printf("Errow writing message: %s", err.Error())
	}
}
