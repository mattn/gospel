package main

import (
	"log"
	"net"
	"os"
	"strconv"
)

func main() {
	s, err := strconv.Atoi(os.Getenv("FD"))
	if err != nil {
		log.Fatal(err)
	}

	f := os.NewFile(uintptr(s), "sysfile")
	l, err := net.FileListener(f)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}
	var b [100]byte
	n, err := conn.Read(b[:])
	if err != nil {
		log.Fatal(err)
	}
	println(n, b[:n])
}
