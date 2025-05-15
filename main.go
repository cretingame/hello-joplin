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
// I think to make `JoplinRoot` compatible with the `InoEmbedded` interface
var _ = (fs.NodeGetattrer)((*JoplinRoot)(nil))
var _ = (fs.NodeOnAdder)((*JoplinRoot)(nil))

func main() {
	var items []joplin.Node

	token, err := joplin.Authenticate(host, tokenLocation)
	if err != nil {
		panic(err)
	}

	folders, err := joplin.GetItems(host, token, "folders")
	if err != nil {
		panic(err)
	}
	for i := range folders {
		folderNode := joplin.FolderNode{
			Id:        folders[i].Id,
			Parent_id: folders[i].Parent_id,
			Title:     folders[i].Title,
		}
		items = append(items, &folderNode)
	}

	notes, err := joplin.GetItems(host, token, "notes")
	if err != nil {
		panic(err)
	}
	for i := range notes {
		noteResponse, err := joplin.GetNote(host, token, notes[i].Id)
		if err != nil {
			panic(err)
		}

		noteNode := joplin.NoteNode{
			Id:        notes[i].Id,
			Parent_id: notes[i].Parent_id,
			Title:     notes[i].Title,
			File: &fs.MemRegularFile{
				Data: []byte(noteResponse.Body),
				Attr: fuse.Attr{
					Mode: 0444,
					// TODO: Change the Owner
				},
			},
		}

		items = append(items, &noteNode)
	}

	root := JoplinRoot{
		items: items,
	}

	debug := flag.Bool("debug", false, "print debug data")
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n  hello-joplin MOUNTPOINT")
	}
	opts := &fs.Options{}
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

type JoplinRoot struct {
	fs.Inode
	items []joplin.Node
}

func (r *JoplinRoot) OnAdd(ctx context.Context) {
	tree := joplin.BuildTree(r.items)

	addNode(ctx, &r.Inode, tree)

	log.Println("Add finished")
}

func addNode(ctx context.Context, parentInode *fs.Inode, items []*joplin.Node) {
	for i := range items {
		child := items[i]

		switch v := (*child).(type) {
		case *joplin.FolderNode:
			childInode := parentInode.NewPersistentInode(
				ctx, &fs.Inode{}, fs.StableAttr{Mode: syscall.S_IFDIR})

			parentInode.AddChild(v.Title, childInode, false)
			addNode(ctx, childInode, v.Children)
		case *joplin.NoteNode:
			childInode := parentInode.NewPersistentInode(
				ctx, v.File, v.File.StableAttr())

			parentInode.AddChild(v.Title+".md", childInode, false)
			addNode(ctx, childInode, v.Children)
		default:
			panic("not handled")
		}
	}
}

func (r *JoplinRoot) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = 0755
	return 0
}
