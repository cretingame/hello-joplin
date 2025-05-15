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
			Children:  folders[i].Children,
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
			Children:  notes[i].Children,
			Body:      noteResponse.Body,
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

	tree := joplin.BuildTree(r.items)

	addFolder(ctx, &r.Inode, tree)

	log.Println("Add finished")
}

func addFolder(ctx context.Context, parentInode *fs.Inode, items []*joplin.Node) {
	for i := range items {
		child := items[i]

		// TODO: differenciate files and folder
		childInode := parentInode.NewPersistentInode(
			ctx, &fs.Inode{}, fs.StableAttr{Mode: syscall.S_IFDIR})

		parentInode.AddChild((*child).Header().Title, childInode, false)

		addFolder(ctx, childInode, (*child).Header().Children)
	}
}

// TODO: for later, I have to dowload the file content first
func addFile(ctx context.Context, parentInode *fs.Inode, items []*joplin.Node) {
	for i := range items {
		child := items[i]

		// TODO: differenciate files and folder
		childInode := parentInode.NewPersistentInode(
			ctx, &fs.Inode{}, fs.StableAttr{Mode: syscall.S_IFDIR})

		parentInode.AddChild((*child).Header().Title, childInode, false)

		addFile(ctx, childInode, (*child).Header().Children)
	}
}

func (r *JoplinRoot) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = 0755
	return 0
}
