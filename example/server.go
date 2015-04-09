package main

import (
	"fmt"
	"github.com/mattn/gospel"
	"log"
)

func main() {
	l, err := gospel.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		for {
			var b [256]byte
			n, err := c.Read(b[:])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(b[:n]))
		}
	}
}
