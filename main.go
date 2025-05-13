package main

import (
	"fmt"
	"hello-joplin/internal/joplin"
)

const (
	tokenLocation = "./token"
	host          = "http://localhost:41184"
)

func main() {
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

	tree := buildTree(items)
	printTree(tree, 0)
}

func buildTree(nodes []joplin.Item) []*joplin.Item {
	nodeMap := make(map[string]*joplin.Item)
	var roots []*joplin.Item

	for i := range nodes {
		nodeCopy := nodes[i]
		nodeMap[nodeCopy.Id] = &nodeCopy
		// Why the node copy ?
	}

	for i := range nodes {
		node := nodeMap[nodes[i].Id]
		if node.Parent_id == "" {
			roots = append(roots, node)
		} else if parent, ok := nodeMap[node.Parent_id]; ok {
			parent.Children = append(parent.Children, node)
		}
	}

	return roots
}

func printTree(nodes []*joplin.Item, level int) {
	for _, node := range nodes {
		out := ""
		for i := 0; i < level*2; i++ {
			out = out + " "
		}
		out = out + node.Title
		fmt.Println(out)
		printTree(node.Children, level+1)
	}
}
