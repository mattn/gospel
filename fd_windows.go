package gospel

import (
	"errors"
	"io/ioutil"
	"net"
	"os"
	"syscall"
	"time"
	"unsafe"
)

func (g *Gospel) exec(l net.Listener) (*os.Process, error) {
	fd := sysfd(l)

	f, err := ioutil.TempFile(os.TempDir(), "gospel-fd")
	if err != nil {
		return nil, err
	}
	os.Setenv("GOSPEL_FD", f.Name())
	err = g.cmd.Start()
	if err != nil {
		f.Close()
		os.Remove(f.Name())
		return nil, err
	}
	b := make([]byte, int(unsafe.Sizeof(syscall.WSAProtocolInfo{})))
	err = syscall.WSADuplicateSocket(syscall.Handle(fd), uint32(g.cmd.Process.Pid), (*syscall.WSAProtocolInfo)(unsafe.Pointer(&b[0])))
	if err != nil {
		f.Close()
		os.Remove(f.Name())
		return nil, err
	}
	f.Write(b)
	f.Close()
	return g.cmd.Process, err
}

func ListenerFromEnv() (net.Listener, error) {
	fn := os.Getenv("GOSPEL_FD")
	l := int(unsafe.Sizeof(syscall.WSAProtocolInfo{}))
	var b []byte
	var err error
	for n := 0; n < 3; n++ {
		b, err = ioutil.ReadFile(fn)
		if len(b) == l {
			break
		}
		time.Sleep(1e9)
	}
	if len(b) == 0 {
		return nil, errors.New("server not found")
	}
	if err != nil {
		return nil, err
	}
	pi := (*syscall.WSAProtocolInfo)(unsafe.Pointer(&b[0]))
	fd, err := syscall.WSASocket(-1, -1, -1, pi, 0, 0)
	if err != nil {
		return nil, err
	}

	/*
	syscall.SetNonblock(syscall.Handle(fd), true)

	sa, err := syscall.Getsockname(syscall.Handle(fd))
	if err != nil {
		return nil, err
	}
	ta := &net.IPAddr{IP: sa.(*syscall.SockaddrInet4).Addr[0:]}
	return &Listener{uintptr(fd), sa, ta}, nil
	*/
	return net.FileListener(os.NewFile(uintptr(fd), "sysfile"))
}
