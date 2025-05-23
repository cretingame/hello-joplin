package main

import (
	"flag"
	"joplin-fuse/internal/joplin"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

const (
	tokenLocation = "./token"
	host          = "http://localhost:41184"
)

func main() {
	root, err := joplin.NewRoot(host, tokenLocation)
	if err != nil {
		log.Fatal(err)
	}

	debug := flag.Bool("debug", false, "print debug data")
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n  hello-joplin MOUNTPOINT")
	}
	opts := &fs.Options{
		UID: fuse.CurrentOwner().Uid,
		GID: fuse.CurrentOwner().Gid,
	}
	opts.Debug = *debug
	server, err := fs.Mount(flag.Arg(0), &root, opts)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		server.Unmount()
	}()

	server.Wait()
}
