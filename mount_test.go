package gofuse

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestMount(t *testing.T) {
	dir := "/tmp/xx"
	// st := syscall.Stat_t{}
	// f, _ := os.Open(dir)
	// syscall.Fstat(int(f.Fd()), &st)
	// fmt.Printf("%+v\n", st)
	// f.Close()

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
