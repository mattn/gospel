// +build !windows

package gospel

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func (g *Gospel) exec(l net.Listener) (*os.Process, error) {
	fd := sysfd(l)

	g.cmd.Env = []string{fmt.Sprintf("GOSPEL_FD=%d", fd)}
	g.cmd.ExtraFiles = []*os.File{os.Stdin, os.Stdout, os.Stderr, os.NewFile(uintptr(fd), "sysfile")}
	err := g.cmd.Start()
	return g.cmd.Process, err
}

func ListenerFromEnv() (net.Listener, error) {
	fd, err := strconv.Atoi(os.Getenv("GOSPEL_FD"))
	if err != nil {
		return nil, err
	}
	return net.FileListener(os.NewFile(uintptr(fd), "sysfile"))
}
