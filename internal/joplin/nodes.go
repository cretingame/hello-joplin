package joplin

import "github.com/hanwen/go-fuse/v2/fs"

type Node interface {
	Base() NodeBase
	AddChild(n *Node)
}

type NodeBase struct {
	Id        string
	Parent_id string
	Name      string
	Children  []*Node
}

type FolderNode struct {
	Id        string
	Parent_id string
	Name      string
	Children  []*Node
}

func (fn FolderNode) Base() NodeBase {
	return NodeBase{
		Id:        fn.Id,
		Parent_id: fn.Parent_id,
		Name:      fn.Name,
		Children:  fn.Children,
	}
}

func (fn *FolderNode) AddChild(n *Node) {
	fn.Children = append(fn.Children, n)
}

type NoteNode struct {
	Id        string
	Parent_id string
	Name      string
	Children  []*Node

	File *fs.MemRegularFile
}

func (nn NoteNode) Base() NodeBase {
	return NodeBase{
		Id:        nn.Id,
		Parent_id: nn.Parent_id,
		Name:      nn.Name,
		Children:  nn.Children,
	}
}

func (nn *NoteNode) AddChild(n *Node) {
	nn.Children = append(nn.Children, n)
}
