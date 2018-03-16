package gofuse

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestMount(t *testing.T) {
	dir := "/tmp/xx"
	fs, err := NewFileSystem(dir, nil)
	if err != nil {
		t.Fatal(err)
	} else {
		defer fs.Close()
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT,
		syscall.SIGTERM)
	<-c
}
