package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const tokenLocation = "./token"

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
	page := 0

	for hasMore {
		// TODO: url constant
		req := fmt.Sprintf("http://localhost:41184/folders?token=%s&page=%d", token, page)
		fmt.Println(req)
		response, err := http.Get(req)
		if err != nil {
			return err
		}

		bs, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var v map[string]any
		var v2 JoplinPage
		json.Unmarshal(bs, &v)
		json.Unmarshal(bs, &v2)

		// fmt.Println(v)
		// fmt.Println(v["has_more"])
		// fmt.Println(v["items"])
		// item structure
		// items := []item{}
		// map[deleted_time:0 id:2c378d6176a446b5bfa4adfb89ffe27c parent_id:df457761ed9c422f9826266e881ea68e title:Shared Library]

		fmt.Println("v2", v2)

		items := v["items"].([]any)

		for i, item := range items {
			fmt.Println(i, item)
		}

		var ok bool
		hasMore, ok = v["has_more"].(bool)
		if !ok {
			hasMore = false
		}

		page++
	}

	return nil
}
