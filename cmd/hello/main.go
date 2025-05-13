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
	token, err := joplin.Authenticate(host, tokenLocation)
	if err != nil {
		panic(err)
	}
	fmt.Printf("token <%s>\n", token)

	folders, err := joplin.GetItems(host, token, "folders")
	if err != nil {
		panic(err)
	}
	for i, folder := range folders {
		fmt.Println(i, folder)
	}

	notes, err := joplin.GetItems(host, token, "notes")
	if err != nil {
		panic(err)
	}
	for i, item := range notes {
		fmt.Println(i, item)
	}

	note, err := joplin.GetNote(host, token, notes[0].Id)
	if err != nil {
		panic(err)
	}
	fmt.Println(note)
}
