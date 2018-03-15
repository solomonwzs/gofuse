package gofuse

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
)

func newSocketpair() (f0, f1 *os.File, err error) {
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, nil, os.NewSyscallError("socketpair", err)
	}
	f0 = os.NewFile(uintptr(fds[0]), "socketpair-0")
	f1 = os.NewFile(uintptr(fds[1]), "socketpair-1")
	return
}

func getConnFromSocket(socket *os.File) (fd int, err error) {
	conn, err := net.FileConn(socket)
	if err != nil {
		return
	}
	defer conn.Close()

	uconn, ok := conn.(*net.UnixConn)
	if !ok {
		err = fmt.Errorf("gofuse: expected UnixConn, got %T", conn)
		return
	}

	buf := make([]byte, 32)
	oob := make([]byte, 32)
	_, oobn, _, _, err := uconn.ReadMsgUnix(buf, oob)
	if err != nil {
		return
	}

	scMsgs, err := syscall.ParseSocketControlMessage(oob[:oobn])
	if err != nil {
		return
	} else if len(scMsgs) != 1 {
		err = fmt.Errorf(
			"gofuse: expected got 1 SocketControlMessage, got: %#v", scMsgs)
		return
	}

	fds, err := syscall.ParseUnixRights(&scMsgs[0])
	if err != nil {
		return
	} else if len(fds) != 1 {
		err = fmt.Errorf("gofuse: expected got 1 fd, got: %#v", fds)
		return
	}
	fd = fds[0]

	return
}

func mount(mountpoint string) (f *os.File, err error) {
	local, remote, err := newSocketpair()
	if err != nil {
		return
	}
	defer local.Close()
	defer remote.Close()

	argv := []string{_CMD_FUSERMOUNT, mountpoint}
	proc, err := os.StartProcess(_CMD_FUSERMOUNT, argv,
		&os.ProcAttr{
			Env:   []string{"_FUSE_COMMFD=3"},
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr, remote},
		})
	if err != nil {
		return
	}

	if state, err := proc.Wait(); err != nil {
		return nil, err
	} else if !state.Success() {
		return nil, fmt.Errorf("fusermount exist: %v", state.Sys())
	}

	fd, err := getConnFromSocket(local)
	if err != nil {
		return
	}

	syscall.CloseOnExec(fd)
	f = os.NewFile(uintptr(fd), "/dev/fuse")

	return
}

func unmount(mountpoint string) (err error) {
	cmd := exec.Command(_CMD_FUSERMOUNT, "-u", mountpoint)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) > 0 {
			output = bytes.TrimRight(output, "\n")
			err = fmt.Errorf("%v: %s", err, output)
		}
		return
	}
	return
}
