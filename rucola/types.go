package rucola

import (
	"net"
)

type Client struct {
	Id     string
	Name   string
	Online bool
	Conn   net.Conn
}

type User struct {
	Id   string
	Name string
}
