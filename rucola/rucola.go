package rucola

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
)

func Rucola() {
	/* Arguments */
	port := flag.String("port", "8001", "the tcp port to listen on")
	datadir := flag.String("db", "./db", "the LevelDB directory")

	flag.Parse()

	/* Chat Database */
	err := initDb(*datadir)
	if err != nil {
		fmt.Printf("Could not init database: %s", err.Error())
		os.Exit(1)
	}

	/* TCP Listener */
	host := "localhost:"+*port
	listener, err := net.Listen("tcp", host)

	if err != nil {
		fmt.Printf("Could not listen: %s", err.Error())
		os.Exit(1)
	}

	/* Signals */
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt)

	go func() {
		for sig := range sigs {
			fmt.Printf("\nCaught %v, shutting down server...\n", sig)

			for _, client := range clients {
				shutdownClient(client)
			}

			closeDb()

			os.Exit(0)
		}
	}()

	fmt.Printf("Chat Server started\n")

	fmt.Printf("Listening on %s\n", host)

	/* TCP Connections */
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Printf(err.Error())
			continue
		}

		go handleClient(conn)
	}
}
