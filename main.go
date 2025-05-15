package main

import (
	"context"
	"flag"
	"hello-joplin/internal/joplin"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
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
		name := folders[i].Title
		name = removeSpecialCharacters(name)
		name = sanitizeFilename(name)
		folderNode := joplin.FolderNode{
			Id:        folders[i].Id,
			Parent_id: folders[i].Parent_id,
			Name:      name,
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

		name := notes[i].Title
		name = removeSpecialCharacters(name)
		name = sanitizeFilename(name)
		name = name + ".md"

		noteNode := joplin.NoteNode{
			Id:        notes[i].Id,
			Parent_id: notes[i].Parent_id,
			Name:      name,
			File: &fs.MemRegularFile{
				Data: []byte(noteResponse.Body),
				Attr: fuse.Attr{
					Mode:  0444,
					Owner: *fuse.CurrentOwner(),
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

// sanitizeFilename removes or replaces characters not allowed in filenames
func sanitizeFilename(name string) string {
	// Replace forbidden characters with underscore
	// Forbidden: \ / : * ? " < > | (Windows)
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	name = invalidChars.ReplaceAllString(name, "_")

	// Trim spaces and dots (Windows does not allow filenames ending with them)
	name = strings.Trim(name, " .")

	// Avoid reserved filenames on Windows (case-insensitive)
	reserved := map[string]bool{
		"CON": true, "PRN": true, "AUX": true, "NUL": true,
		"COM1": true, "COM2": true, "COM3": true, "COM4": true,
		"COM5": true, "COM6": true, "COM7": true, "COM8": true, "COM9": true,
		"LPT1": true, "LPT2": true, "LPT3": true, "LPT4": true,
		"LPT5": true, "LPT6": true, "LPT7": true, "LPT8": true, "LPT9": true,
	}
	upper := strings.ToUpper(name)
	if reserved[upper] {
		name = "_" + name
	}

	// Limit length (255 bytes is a common limit)
	if len(name) > 255 {
		name = name[:255]
	}

	if name == "" {
		return "unnamed"
	}

	return name
}

func removeSpecialCharacters(input string) string {
	var sb strings.Builder
	for _, r := range input {
		if r < 0xFFFF {
			sb.WriteRune(r)
		}
	}
	return sb.String()
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

			parentInode.AddChild(v.Name, childInode, false)
			addNode(ctx, childInode, v.Children)
		case *joplin.NoteNode:
			childInode := parentInode.NewPersistentInode(
				ctx, v.File, v.File.StableAttr())

			parentInode.AddChild(v.Name, childInode, false)
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
