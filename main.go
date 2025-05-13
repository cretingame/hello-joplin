package main

import (
	"fmt"
	"hello-joplin/internal/joplin"
	"os"
	"strings"
	"time"
)

const (
	tokenLocation = "./token"
	host          = "http://localhost:41184"
)

func main() {
	authToken, err := joplin.GetAuthToken(host)
	if err != nil {
		panic(err)
	}

	_, err = os.Stat(tokenLocation)
	if os.IsNotExist(err) {
		token, err := joplin.GetToken(host, authToken)
		for err == joplin.ErrCheckJoplin {
			fmt.Println("Please check joplin application to grant access")
			time.Sleep(1000 * time.Millisecond)
			token, err = joplin.GetToken(host, authToken)
		}
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(tokenLocation, []byte(token), 0644)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	bs, err := os.ReadFile(tokenLocation)
	if err != nil {
		panic(err)
	}
	token := strings.Trim(string(bs), "\n")

	fmt.Printf("token <%s>\n", token)

	// end of authentification

	folders, err := joplin.GetJoplinItems(host, token, "folders")
	if err != nil {
		panic(err)
	}
	for i, folder := range folders {
		fmt.Println(i, folder)
	}

	notes, err := joplin.GetJoplinItems(host, token, "notes")
	if err != nil {
		panic(err)
	}
	for i, item := range notes {
		fmt.Println(i, item)
	}

	note, err := joplin.GetJoplinNote(host, token, notes[0].Id)
	if err != nil {
		panic(err)
	}
	fmt.Println(note)
}
