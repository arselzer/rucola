package main

import (
	"net"
)

type Client struct {
	id     string
	name   string
	online bool
	conn   net.Conn
}

type User struct {
	id   string
	name string
}
