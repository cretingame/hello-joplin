package joplin

import "github.com/hanwen/go-fuse/v2/fs"

type Node interface {
	Header() ItemHeader
	AddChild(n *Node)
}

type FolderNode struct {
	Id        string
	Parent_id string
	Title     string
	Children  []*Node
}

func (fn FolderNode) Header() ItemHeader {
	return ItemHeader{
		Id:        fn.Id,
		Parent_id: fn.Parent_id,
		Title:     fn.Title,
		Children:  fn.Children,
	}
}

func (fn *FolderNode) AddChild(n *Node) {
	fn.Children = append(fn.Children, n)
}

type NoteNode struct {
	Id        string
	Parent_id string
	Title     string
	Children  []*Node

	File *fs.MemRegularFile
}

func (nn NoteNode) Header() ItemHeader {
	return ItemHeader{
		Id:        nn.Id,
		Parent_id: nn.Parent_id,
		Title:     nn.Title,
		Children:  nn.Children,
	}
}

func (nn *NoteNode) AddChild(n *Node) {
	nn.Children = append(nn.Children, n)
}
