package main

import (
	"github.com/mattn/gospel"
	"log"
	"os"
	"os/exec"
)

func main() {
	//cmd := exec.Command(`c:\dev\listen-socket\child.exe`)
	cmd := exec.Command(`c:\dev\go\src\github.com\mattn\gospel\example\echo-server.exe`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	g := gospel.New(cmd)
	err := g.Listen(":8888")
	if err != nil {
		log.Fatal(err)
	}
}
