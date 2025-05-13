package main

import (
	"fmt"
	"hello-joplin/internal/joplin"
	"os"
	"strings"
)

const (
	authTokenLocation = "./auth_token"
	tokenLocation     = "./token"
	host              = "http://localhost:41184"
)

func main() {
	_, err := os.Stat(authTokenLocation)
	if os.IsNotExist(err) {
		authToken, err := joplin.GetAuthToken(host)
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

		token, err := joplin.GetToken(host, authToken)
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

// NOTE: That's an useless abstraction
// OPTIM: I should use a parameter instead of a constant
func saveAuthToken(authToken string) error {
	err := os.WriteFile(authTokenLocation, []byte(authToken), 0644)
	if err != nil {
		return err
	}
	return nil
}

// NOTE: That's an useless abstraction
// OPTIM: I should use a parameter instead of a constant
func readAuthToken() (string, error) {
	bs, err := os.ReadFile(authTokenLocation)
	str := string(bs)
	str = strings.Trim(str, "\n")
	return str, err
}

// NOTE: That's an useless abstraction
// OPTIM: I should use a parameter instead of a constant
func saveToken(token string) error {
	err := os.WriteFile(tokenLocation, []byte(token), 0644)
	if err != nil {
		return err
	}
	return nil
}

// NOTE: That's an useless abstraction
// OPTIM: I should use a parameter instead of a constant
// getToken in the bash script
func readToken() (string, error) {
	bs, err := os.ReadFile(tokenLocation)
	str := string(bs)
	str = strings.Trim(str, "\n")
	return str, err
}
