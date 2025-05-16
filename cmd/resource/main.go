package main

import (
	"fmt"
	"hello-joplin/internal/joplin"
	"log"
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

	id := "73ca02a59f2741bcb98a77516b85b9a5"

	bs, err := joplin.GetRessourceFile(host, token, id)
	if err != nil {
		panic(err)
	}
	log.Println(bs)

	resp, err := joplin.GetRessource(host, token, id)
	if err != nil {
		panic(err)
	}
	log.Println(resp)
}
