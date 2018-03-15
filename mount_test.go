package gofuse

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestMount(t *testing.T) {
	dir := "/tmp/xx"
	fmt.Println(mount(dir))

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT,
		syscall.SIGTERM)
	<-c

	fmt.Println(unmount(dir))
}
