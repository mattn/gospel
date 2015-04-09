package gospel

import (
	"net"
	"os"
	"syscall"
	"time"
	"unsafe"
)

var (
	ws2_32 = syscall.NewLazyDLL("ws2_32.dll")
	procAccept = ws2_32.NewProc("accept")
)

type Listener struct {
	fd uintptr
	sa syscall.Sockaddr
	addr *net.IPAddr
}

type Conn struct {
	fd uintptr
	addr *net.TCPAddr
}

func Listen(n, addr string) (*Listener, error) {
	ad, err := net.ResolveTCPAddr(n, addr)
	if err != nil {
		return nil, err
	}

	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}

	var sa syscall.Sockaddr
	switch n {
	case "tcp", "tcp4":
		sa4 := new(syscall.SockaddrInet4)
		for i := 0; i < len(ad.IP); i++ {
			sa4.Addr[i] = ad.IP[i]
		}
		sa4.Port = ad.Port
		sa = sa4
	case "tcp6":
		sa4 := new(syscall.SockaddrInet6)
		for i := 0; i < len(ad.IP); i++ {
			sa4.Addr[i] = ad.IP[i]
		}
		sa4.Port = ad.Port
		sa = sa4
	}

	err = syscall.Bind(s, sa)
	if err != nil {
		return nil, err
	}

	err = syscall.Listen(s, syscall.SOMAXCONN)
	if err != nil {
		return nil, err
	}

	ssa, err := syscall.Getsockname(syscall.Handle(s))
	if err != nil {
		return nil, err
	}
	ta := &net.IPAddr{IP: ssa.(*syscall.SockaddrInet4).Addr[0:]}
	return &Listener{uintptr(s), sa, ta}, nil
}

func (l *Listener) Addr() net.Addr {
	return l.addr

}

func (l *Listener) Close() error {
	return syscall.Closesocket(syscall.Handle(l.fd))
}

func (l *Listener) Accept() (net.Conn, error) {
	var sa syscall.SockaddrInet4
	sl := unsafe.Sizeof(sa)
	newfd, r1, err := procAccept.Call(uintptr(l.fd), uintptr(unsafe.Pointer(&sa)), uintptr(unsafe.Pointer(&sl)))
	if err != nil && r1 == 0 {
		return nil, err
	}
	//return &Conn{uintptr(newfd), nil}, nil
	return net.FileConn(os.NewFile(newfd, "sysfile"))
}

func (c *Conn) Read(b []byte) (n int, e error) {
	var buf syscall.WSABuf
	buf.Buf = &b[0]
	buf.Len = uint32(len(b))
	var qty, flags uint32
	err := syscall.WSARecv(syscall.Handle(c.fd), &buf, 1, &qty, &flags, nil, nil)
	return int(qty), err
}

func (c *Conn) Write(b []byte) (int, error) {
	return syscall.Write(syscall.Handle(c.fd), b)
}

func (c *Conn) Close() error {
	return syscall.Closesocket(syscall.Handle(c.fd))
}

func (c *Conn) LocalAddr() net.Addr {
	return c.addr
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.addr
}

func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}

