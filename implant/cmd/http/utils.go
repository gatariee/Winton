package http

import (
    "bytes"
    "io"
    "net/http"
)

type HTTPClient struct{}

func NewHTTPClient() *HTTPClient {
    return &HTTPClient{}
}

func (c *HTTPClient) Get(url string) ([]byte, error) {
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

func (c *HTTPClient) PostJSON(url string, json []byte) ([]byte, error) {
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