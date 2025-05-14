package main

import (
	"context"
	"flag"
	"hello-joplin/internal/joplin"
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

// INFO: I don't why know it was added in the example
var _ = (fs.NodeGetattrer)((*JoplinRoot)(nil))
var _ = (fs.NodeOnAdder)((*JoplinRoot)(nil))

func main() {
	debug := flag.Bool("debug", false, "print debug data")
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n  hello MOUNTPOINT")
	}
	opts := &fs.Options{}
	opts.Debug = *debug
	server, err := fs.Mount(flag.Arg(0), &JoplinRoot{}, opts)
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

	// Joplin tree BEGIN

	var items []joplin.Item

	token, err := joplin.Authenticate(host, tokenLocation)
	if err != nil {
		panic(err)
	}

	folders, err := joplin.GetItems(host, token, "folders")
	if err != nil {
		panic(err)
	}

	notes, err := joplin.GetItems(host, token, "notes")
	if err != nil {
		panic(err)
	}

	items = append(items, folders...)
	items = append(items, notes...)

	tree := joplin.BuildTree(items)
	joplin.PrintTree(tree, 0)
}

type JoplinRoot struct {
	fs.Inode
}

func (r *JoplinRoot) OnAdd(ctx context.Context) {
	ch := r.NewPersistentInode(
		ctx, &fs.MemRegularFile{
			Data: []byte("Hello World in file.txt\n"),
			Attr: fuse.Attr{
				Mode: 0644,
				// TODO: Change the Owner
			},
		}, fs.StableAttr{Ino: 2})
	r.AddChild("file.txt", ch, false)
	r.AddChild("file2.txt", ch, false)

	ch2 := r.NewPersistentInode(
		ctx, &fs.Inode{}, fs.StableAttr{Mode: syscall.S_IFDIR})
	ch2.AddChild("fileInDirectory.txt", ch, false)
	r.AddChild("directory", ch2, false)
}

func (r *JoplinRoot) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = 0755
	return 0
}
