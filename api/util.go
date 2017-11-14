package api

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func execute(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", APIKey) // set api key
	req.Header.Set("Content-Type", "application/json")
	return execute(req)
}

func post(url string, data []byte) ([]byte, error) {
	// create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	// make request
	req.Header.Set("x-api-key", APIKey) // set api key
	req.Header.Set("Content-Type", "application/json")
	return execute(req)
}

func delete(url string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", APIKey) // set api key
	return execute(req)
}
