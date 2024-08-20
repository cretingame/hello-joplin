package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const tokenLocation = "./token"

func main() {
	token, err := readToken()
	if err != nil {
		panic(err)
	}
	fmt.Printf("token <%s>\n", token)

	err = listJoplinItems(token)
	if err != nil {
		panic(err)
	}
}

func saveToken() {
	// TODO: implemented in the bash script
}

// getToken in the bash script
func readToken() (string, error) {
	bs, err := os.ReadFile(tokenLocation)
	str := string(bs)
	str = strings.Trim(str, "\n")
	return str, err
}

func listJoplinItems(token string) error {
	hasMore := true

	for hasMore {
		// TODO: url constant
		req := fmt.Sprintf("http://localhost:41184/folders?token=%s&page=%d", token, 0)
		fmt.Println(req)
		response, err := http.Get(req)
		if err != nil {
			return err
		}

		fmt.Println(response)

		bs, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		str := string(bs)

		fmt.Println(str)

		// TODO: should get information from the request
		hasMore = false
	}

	return nil
}
