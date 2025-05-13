package joplin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// curl -X POST "$ADDRESS/auth" | jq '.auth_token' | sed 's/\"//g'
func GetAuthToken(host string) (authToken string, err error) {
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

// https://joplinapp.org/fr/help/dev/spec/clipper_auth
func GetToken(host string, authToken string) (token string, err error) {
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
