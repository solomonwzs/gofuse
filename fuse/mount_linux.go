package fuse

import (
	"bytes"
	"fmt"
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

func getConnFromSocket(sock int) (fd int, err error) {
	buf := make([]byte, 32)
	oob := make([]byte, 32)
	_, oobn, _, _, err := syscall.Recvmsg(sock, buf, oob, 0)

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

func mount(mountpoint string, conf *MountConfig) (
	f *os.File, err error) {
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		return
	}
	defer syscall.Close(fds[0])
	defer syscall.Close(fds[1])

	argv := []string{_CMD_FUSERMOUNT, "-o", "rw,nodev,nosuid,sync",
		"--", mountpoint}
	syscall.Syscall(syscall.SYS_FCNTL, uintptr(fds[0]), uintptr(syscall.F_SETFD), 0)
	proc, err := os.StartProcess(_CMD_FUSERMOUNT, argv,
		&os.ProcAttr{
			Env: []string{fmt.Sprintf("_FUSE_COMMFD=%d", fds[0])},
		})
	if err != nil {
		return
	}

	if state, err := proc.Wait(); err != nil {
		return nil, err
	} else if !state.Success() {
		return nil, fmt.Errorf("fusermount exist: %v", state.Sys())
	}

	fd, err := getConnFromSocket(fds[1])
	if err != nil {
		return
	}

	syscall.CloseOnExec(fd)
	f = os.NewFile(uintptr(fd), "/dev/fuse")

	return
}

func umount(mountpoint string) (err error) {
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
