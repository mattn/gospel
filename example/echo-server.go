package main

import (
	"fmt"
	"github.com/mattn/gospel"
	"log"
	"net"
)

func echo_handler(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 128)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Received: ", string(buf[:n]))
		conn.Write(buf[:n])
	}
}

func main() {
	l, err := gospel.ListenerFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("CONNECTED")
	for {
		conn, e := l.Accept()
		if e != nil {
			log.Fatal(e)
			return
		}
		fmt.Println("ACCEPTED")
		go echo_handler(conn)
	}
}
