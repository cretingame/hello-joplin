package main

import (
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

	tree := joplin.BuildTree(items)
	joplin.PrintTree(tree, 0)
}
