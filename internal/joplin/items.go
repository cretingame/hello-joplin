package joplin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// https://joplinapp.org/fr/help/api/references/rest_api/#properties-1
type PageResponse struct {
	Has_more bool
	Items    []ItemResponse
}

type ItemResponse struct {
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

type NoteResponse struct {
	Id                     string
	Parent_id              string // ID of the notebook that contains this note. Change this ID to move the note to a different notebook.
	Title                  string // The note title.
	Body                   string // The note body, in Markdown. May also contain HTML.
	Created_time           int    // When the note was created.
	Updated_time           int    // When the note was last updated.
	Is_conflict            int    // Tells whether the note is a conflict or not.
	Latitude               float64
	Longitude              float64
	Altitude               float64
	Author                 string
	Source_url             string // The full URL where the note comes from.
	Is_todo                int    // Tells whether this note is a todo or not.
	Todo_due               int    // When the todo is due. An alarm will be triggered on that date.
	Todo_completed         int    // Tells whether todo is completed or not. This is a timestamp in milliseconds.
	Source                 string
	Source_application     string
	Application_data       string
	Order                  float64
	User_created_time      int // When the note was created. It may differ from created_time as it can be manually set by the user.
	User_updated_time      int // When the note was last updated. It may differ from updated_time as it can be manually set by the user.
	Encryption_cipher_text string
	Encryption_applied     int
	Markup_language        int
	Is_shared              int
	Share_id               string
	Conflict_original_id   string
	Master_key_id          string
	User_data              string
	Deleted_time           int
	Body_html              string // Note body, in HTML format
	Base_url               string // If body_html is provided and contains relative URLs, provide the base_url parameter too so that all the URLs can be converted to absolute ones. The base URL is basically where the HTML was fetched from, minus the query (everything after the '?'). For example if the original page was https://stackoverflow.com/search?q=%5Bjava%5D+test, the base URL is https://stackoverflow.com/search.
	Image_data_url         string // An image to attach to the note, in Data URL format.
	Crop_rect              string // If an image is provided, you can also specify an optional rectangle that will be used to crop the image. In format { x: x, y: y, width: width, height: height }
}

type FolderResponse struct {
	Id                     string
	Title                  string // The folder title.
	Created_time           int    // When the folder was created.
	Updated_time           int    // When the folder was last updated.
	User_created_time      int    // When the folder was created. It may differ from created_time as it can be manually set by the user.
	User_updated_time      int    // When the folder was last updated. It may differ from updated_time as it can be manually set by the user.
	Encryption_cipher_text string
	Encryption_applied     int
	Parent_id              string
	Is_shared              int
	Share_id               string
	Master_key_id          string
	Icon                   string
	User_data              string
	Deleted_time           int
}

type RessourceResponse struct {
	Id                        string
	Title                     string // The resource title.
	Mime                      string
	Filename                  string
	Created_time              int // When the resource was created.
	Updated_time              int // When the resource was last updated.
	User_created_time         int // When the resource was created. It may differ from created_time as it can be manually set by the user.
	User_updated_time         int // When the resource was last updated. It may differ from updated_time as it can be manually set by the user.
	File_extension            string
	Encryption_cipher_text    string
	Encryption_applied        int
	Encryption_blob_encrypted int
	Size                      int
	Is_shared                 int
	Share_id                  string
	Master_key_id             string
	User_data                 string
	Blob_updated_time         int
	Ocr_text                  string
	Ocr_details               string
	Ocr_status                int
	Ocr_error                 string
}

func GetItems(host string, token string, joplinType string) (items []ItemResponse, err error) {
	hasMore := true
	page := 0

	for hasMore {
		req := fmt.Sprintf("%s/%s?token=%s&page=%d", host, joplinType, token, page)
		response, err := http.Get(req)
		if err != nil {
			return items, err
		}

		bs, err := io.ReadAll(response.Body)
		if err != nil {
			return items, err
		}

		var jPage PageResponse
		err = json.Unmarshal(bs, &jPage)
		if err != nil {
			return items, err
		}

		hasMore = jPage.Has_more

		items = append(items, jPage.Items...)
		page++
	}

	return items, err
}

func GetNote(host string, token string, id string) (note NoteResponse, err error) {
	req := fmt.Sprintf("%s/notes/%s?token=%s&fields=title,body", host, id, token)
	response, err := http.Get(req)
	if err != nil {
		return
	}

	bs, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bs, &note)
	if err != nil {
		return
	}

	return
}

func GetFolder(host string, token string, id string) (folder FolderResponse, err error) {
	req := fmt.Sprintf("%s/folders/%s?token=%s&fields=title,body", host, id, token)
	response, err := http.Get(req)
	if err != nil {
		return
	}

	bs, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bs, &folder)
	if err != nil {
		return
	}

	return
}

func BuildTree(nodes []Node) []*Node {
	nodeMap := make(map[string]*Node)
	var roots []*Node

	for i := range nodes {
		nodeMap[nodes[i].Base().Id] = &nodes[i]
	}

	for i := range nodes {
		node := nodeMap[nodes[i].Base().Id]
		if (*node).Base().Parent_id == "" {
			roots = append(roots, node)
		} else if parent, ok := nodeMap[(*node).Base().Parent_id]; ok {
			(*parent).AddChild(node)
		}
	}

	return roots
}

func PrintTree(nodes []*Node, level int) {
	for _, node := range nodes {
		out := ""
		for i := 0; i < level*2; i++ {
			out = out + " "
		}
		out = out + (*node).Base().Title
		fmt.Println(out)
		PrintTree((*node).Base().Children, level+1)
	}
}
