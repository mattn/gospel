package main

import (
	"flag"
	"github.com/mattn/gospel"
	"log"
	"os"
	"os/exec"
)

var addr = flag.String("a", "listen address", ":8888")

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	g := gospel.New(cmd)
	err := g.Listen(*addr)
	if err != nil {
		log.Fatal(err)
	}
}
