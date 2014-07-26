package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
)

func broadcastMsg(origin string, msg string, showOrigin bool) {
	for i := 0; i < len(clients); i++ {
		if clients[i].name != origin {
			if showOrigin {
				writeClient(clients[i], "["+origin+"]: "+msg)
			} else {
				writeClient(clients[i], msg)
			}
		}
	}
}

func handleClient(conn net.Conn) {
	fmt.Printf("client connected: %s\n", conn.RemoteAddr().String())

	/* Set up client object */

	// Get name of user
	_, err := conn.Write([]byte("  select a name...\n"))

	if err != nil {
		fmt.Printf(err.Error())
	}

	nameBuf := make([]byte, 64)

	_, err = conn.Read(nameBuf)

	if err != nil {
		fmt.Printf(err.Error())
	}

	name := strings.TrimRight(string(nameBuf), "\x00")
	name = strings.Replace(name, "\n", "", -1)

	// client id
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)

	id := hex.EncodeToString(randomBytes)

	users := getUsers()

	fmt.Println(users)

	// Add client to shared list
	client := Client{id: id, name: name, conn: conn}

	clients = append(clients, client)

	broadcastMsg(client.name, client.name+" joined the chat", false)
	writeClient(client, "your id: ")
	writeClient(client, client.id)
	writeClient(client, "remember it to log in again.")

	writeClient(client, "welcome "+name)

	if err != nil {
		fmt.Printf("Could not write to client: %s\n", err.Error())
	}

	// Chat loop
	for {
		buff := make([]byte, 1024)
		_, err := conn.Read(buff)

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Printf("%s quit the chat", name)
				conn.Close()
				return
			}
		}

		msg := string(buff)
		msg = strings.Replace(msg, "\n", "", -1)

		isCommand, err := regexp.MatchString("/\\w+", msg)

		if isCommand {
			reCmd, _ := regexp.Compile("(\\w+)")
			matches := reCmd.FindAllStringSubmatch(msg, -1)
			cmd := matches[0][0]

			var args []string

			for i := 1; i < len(matches); i++ {
				args = append(args, matches[i][0])
			}

			fmt.Printf(name+" issued command: %s\n", cmd)

			command(name, string(cmd), args)
		} else {
			fmt.Printf(name + " wrote '" + msg + "'" + "\n")

			broadcastMsg(name, msg, true)
		}
	}
}

func main() {
	port := flag.String("port", "8001", "the tcp port to listen on")

	flag.Parse()

	listener, err := net.Listen("tcp", "localhost:"+*port)

	if err != nil {
		fmt.Printf("Could not listen: %s", err.Error())
	}

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			fmt.Printf("Caught %v, shutting down server...", sig)

			for _, client := range clients {
				writeClient(client, "server is shutting down")
				client.conn.Close()
			}

			os.Exit(0)
		}
	}()

	fmt.Printf("Chat Server started\n")

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Printf(err.Error())
			continue
		}

		go handleClient(conn)
	}
}
