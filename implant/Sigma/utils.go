package main

import (
	"bytes"
	"io"
	"net/http"
)

func http_get(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body) 
	if err != nil {
		return nil, err
	}

	return body, nil
}

func http_post_json(url string, json []byte ) ([]byte, error) {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return io.ReadAll(res.Body)

}
