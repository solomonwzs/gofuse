package gofuse

import (
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/solomonwzs/gofuse/filesystem/simplefs"
	"github.com/solomonwzs/gofuse/fuse"
)

func TestMount(t *testing.T) {
	dir := "/tmp/xx"
	// st := syscall.Stat_t{}
	// f, _ := os.Open(dir)
	// syscall.Fstat(int(f.Fd()), &st)
	// fmt.Printf("%+v\n", st)
	// f.Close()

	fs, err := fuse.NewFuseServer(dir, nil, simplefs.NewExampleSimpleFS())
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
