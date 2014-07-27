package rucola

import (
	"fmt"
	"net"
	"strings"
	"regexp"
)

var clients []Client

func findClient(name string) (Client, int) {
	for i := 0; i < len(clients); i++ {
		if clients[i].Name == name {
			return clients[i], i
		}
	}
	return Client{Id: "", Name: "", Online: false}, -1
}

func addClient(client Client) {
	fmt.Printf("%#s", clients)
	clients = append(clients, client)
}

func removeClient(client Client) {
	_, i := findClient(client.Name)

	clients = append(clients[:i], clients[1+i:]...)
}

func writeClient(client Client, msg string) {
	_, err := client.Conn.Write([]byte("-- " + msg + "\n"))

	if err != nil {
		fmt.Printf("Error writing message: %s", err.Error())
	}
}

func writeClientMultiline(client Client, messages []string) {
	var text string

	for _, msg := range messages {
		text = text + "  " + msg + "\n"
	}

	_, err := client.Conn.Write([]byte(text))

	if err != nil {
		fmt.Printf("Error writing message: %s", err.Error())
	}
}

func broadcastMsg(origin string, msg string, showOrigin bool) {
	for i := 0; i < len(clients); i++ {
		if clients[i].Name != origin {
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

	_, err := conn.Write([]byte("-- select a name:\n"))

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

	id := generateId()

	client := Client{Id: id, Name: name, Conn: conn}

	justAskedForId := false

	if userExists(client) {

		_, clientIndex := findClient(client.Name)

		if clientIndex != -1 {
			writeClient(client, client.Name + " is already in the chat")
			client.Conn.Close()
			return
		}

		fmt.Println("User exists: %s", client.Name)
		writeClient(client, "id:")

		idBuf := make([]byte, 32)
		_, err = client.Conn.Read(idBuf)

		client.Id = strings.Trim(strings.Replace(string(idBuf), "\n", "", -1), "\x00")

		authenticated, _ := authUser(client)

		justAskedForId = true

		if !authenticated {
			writeClient(client, "wrong id")
			client.Conn.Close()
			return
		} else {
			writeClient(client, "success")
		}
	}

	user := User{Name: client.Name, Id: client.Id}

	saveUser(user)

	addClient(client)

	/* Chat notifications */

	broadcastMsg(client.Name, client.Name+" joined the chat", false)

	if !justAskedForId {
		writeClientMultiline(client, []string{"your id: ", client.Id, "remember it to log in again."})
	}

	/* Chat Loop */

	for {
		buff := make([]byte, 1024)
		_, err := conn.Read(buff)

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Printf("%s quit the chat\n", name)
				conn.Close()
				removeClient(client)
				return
			}
		}

		msg := string(buff)
		msg = strings.Replace(msg, "\n", "", -1)

		isCommand, _ := regexp.MatchString("/\\w+", msg)

		if isCommand {
			reCmd, _ := regexp.Compile("(\\w+)")
			matches := reCmd.FindAllStringSubmatch(msg, -1)

			/* Stuff like `/msg alwin hello you`
				will become []string{'msg' 'alwin' 'hello' 'you'}
			*/

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

func shutdownClient(client Client) {
	writeClient(client, "server is shutting down")
	client.Conn.Close()
	// Don't have to remove because the program was closed anyway
}
