package gospel

import (
	"errors"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"runtime"
	"strconv"
	"syscall"
	"unsafe"
)

type Gospel struct {
	cmd    *exec.Cmd
}

func New(cmd *exec.Cmd) *Gospel {
	return &Gospel{cmd}
}

func (g *Gospel) Listen(addr string) error {
	var l net.Listener
	var err error
	if runtime.GOOS == "windows" {
		l, err = Listen("tcp", addr)
	} else {
		l, err = net.Listen("tcp", addr)
	}
	if err != nil {
		return err
	}

	p, err := g.exec(l)
	if err != nil {
		return err
	}
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP)
	for {
		switch sig := <-c; sig {
		case syscall.SIGHUP:
			child, err := g.exec(l)
			if err != nil {
				return err
			}
			p.Signal(syscall.SIGINT)
			p.Wait()
			p = child
		case syscall.SIGINT:
			signal.Stop(c)
			l.Close()
			p.Signal(syscall.SIGINT)
			_, err := p.Wait()
			return err
		}
	}
}

func sysfd(l net.Listener) uintptr {
	if ll, ok := l.(*Listener); ok {
		return ll.fd
	}
	return *(*uintptr)(unsafe.Pointer(reflect.ValueOf(l).Elem().FieldByName("fd").Elem().FieldByName("sysfd").Addr().Pointer()))
}

func fd() (uintptr, error) {
	s := os.Getenv("GOSPEL_FD")
	if s == "" {
		return 0, errors.New("server not found")
	}
	fd, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return uintptr(fd), nil
}
