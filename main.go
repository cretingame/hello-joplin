package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	tokenLocation = "./token"
	host          = "http://localhost:41184"
)

// https://joplinapp.org/fr/help/api/references/rest_api/#properties-1
type JoplinPage struct {
	Has_more bool
	Items    []JoplinItem
}

type JoplinItem struct {
	Id                     string
	Parent_id              string
	Title                  string
	Created_time           int // When the folder was created.
	Updated_time           int // When the folder was last updated.
	User_created_time      int // When the folder was created. It may differ from created_time as it can be manually set by the user.
	User_updated_time      int // When the folder was last updated. It may differ from updated_time as it can be manually set by the user.
	Encryption_cipher_text string
	Encryption_applied     int
	Is_shared              int
	Share_id               string
	Master_key_id          string
	Icon                   string
	User_data              string
	Deleted_time           int
}

func main() {
	token, err := readToken()
	if err != nil {
		panic(err)
	}
	fmt.Printf("token <%s>\n", token)

	items, err := listJoplinItems(token)
	if err != nil {
		panic(err)
	}
	for i, item := range items {
		fmt.Println(i, item)
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

func listJoplinItems(token string) (items []JoplinItem, err error) {
	hasMore := true
	page := 0

	for hasMore {
		// TODO: url constant
		req := fmt.Sprintf("%s/folders?token=%s&page=%d", host, token, page)
		response, err := http.Get(req)
		if err != nil {
			return items, err
		}

		bs, err := io.ReadAll(response.Body)
		if err != nil {
			return items, err
		}

		var jPage JoplinPage
		json.Unmarshal(bs, &jPage)

		hasMore = jPage.Has_more

		items = append(items, jPage.Items...)
		page++
	}

	return items, err
}

// TODO: to be tested
func readJoplinNote(token string, id string) (content []byte, err error) {
	req := fmt.Sprintf("%s/notes/%s?token=%s", host, id, token)
	response, err := http.Get(req)
	if err != nil {
		return
	}

	bs, err := io.ReadAll(response.Body)

	var v any

	json.Unmarshal(bs, &v)

	fmt.Println(v)

	return
}
