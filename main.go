package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	authTokenLocation = "./auth_token"
	tokenLocation     = "./token"
	host              = "http://localhost:41184"
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
	_, err := os.Stat(authTokenLocation)
	if os.IsNotExist(err) {
		authToken, err := getAuthToken()
		if err != nil {
			panic(err)
		}
		fmt.Println("create authToken file with token:", authToken)
		err = saveAuthToken(authToken)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	_, err = os.Stat(tokenLocation)
	if os.IsNotExist(err) {
		authToken, err := readAuthToken()
		if err != nil {
			panic(err)
		}

		token, err := getToken(authToken)
		if err != nil {
			panic(err)
		}
		fmt.Println("create token file with token:", token)
		err = saveToken(token)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	token, err := readToken()
	if err != nil {
		panic(err)
	}
	fmt.Printf("token <%s>\n", token)

	folders, err := listJoplinItems(token, "folders")
	if err != nil {
		panic(err)
	}
	for i, folder := range folders {
		fmt.Println(i, folder)
	}

	notes, err := listJoplinItems(token, "notes")
	if err != nil {
		panic(err)
	}
	for i, item := range notes {
		fmt.Println(i, item)
	}

	note, err := readJoplinNote(token, notes[0].Id)
	if err != nil {
		panic(err)
	}
	fmt.Println(note)
}

// curl -X POST "$ADDRESS/auth" | jq '.auth_token' | sed 's/\"//g'
func getAuthToken() (authToken string, err error) {
	var body io.Reader
	var v map[string]string
	var ok bool

	resp, err := http.Post(host+"/auth", "application/json", body)
	if err != nil {
		return
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bs, &v)
	if err != nil {
		return
	}

	authToken, ok = v["auth_token"]
	if !ok {
		err = errors.New("parsing auth JSON failed")
		return
	}

	return
}

func saveAuthToken(authToken string) error {
	err := os.WriteFile(authTokenLocation, []byte(authToken), 0644)
	if err != nil {
		return err
	}
	return nil
}

func readAuthToken() (string, error) {
	bs, err := os.ReadFile(authTokenLocation)
	str := string(bs)
	str = strings.Trim(str, "\n")
	return str, err
}

// https://joplinapp.org/fr/help/dev/spec/clipper_auth
func getToken(authToken string) (token string, err error) {
	var v map[string]string

	req := fmt.Sprintf("%s/auth/check?auth_token=%s", host, authToken)
	resp, err := http.Get(req)
	if err != nil {
		return
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bs, &v)
	if err != nil {
		return
	}

	status, ok := v["status"]
	if !ok {
		err = errors.New("parsing status from token JSON failed")
		return
	}
	if status != "accepted" {
		err = fmt.Errorf("getToken status: %s", status)
		return
	}

	token, ok = v["token"]
	if !ok {
		err = errors.New("parsing token from token JSON failed")
		return
	}

	return
}

func saveToken(token string) error {
	err := os.WriteFile(tokenLocation, []byte(token), 0644)
	if err != nil {
		return err
	}
	return nil
}

// getToken in the bash script
func readToken() (string, error) {
	bs, err := os.ReadFile(tokenLocation)
	str := string(bs)
	str = strings.Trim(str, "\n")
	return str, err
}

func listJoplinItems(token string, joplinType string) (items []JoplinItem, err error) {
	hasMore := true
	page := 0

	for hasMore {
		// TODO: url constant
		req := fmt.Sprintf("%s/%s?token=%s&page=%d", host, joplinType, token, page)
		response, err := http.Get(req)
		if err != nil {
			return items, err
		}

		bs, err := io.ReadAll(response.Body)
		if err != nil {
			return items, err
		}

		var jPage JoplinPage
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

type JoplinNote struct {
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

func readJoplinNote(token string, id string) (note JoplinNote, err error) {
	req := fmt.Sprintf("%s/notes/%s?token=%s&fields=title,body", host, id, token)
	fmt.Println("req", req)
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

type JoplinFolder struct {
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

// TODO: to be tested
func readJoplinFolder(token string, id string) (folder JoplinFolder, err error) {
	req := fmt.Sprintf("%s/folders/%s?token=%s&fields=title,body", host, id, token)
	fmt.Println("req", req)
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

type JoplinRessource struct {
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
