package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
)

func ExtractEvents(message string) [][]string {
	re := regexp.MustCompile(`(?mU)Day (\d+), (\d+\:\d+\:\d+)\: (.+)!`)
	events := re.FindAllStringSubmatch(message, -1)
	return events
}

func PrintWebhook(http_request *http.Request) error {
	// Save a copy of this request for debugging.
	request_dump, err := httputil.DumpRequest(http_request, true)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(string(request_dump))
	return nil
}

func RelayWebhook(http_request *http.Request, webhook string) error {
	body, err := io.ReadAll(http_request.Body)
	if err != nil {
		return err
	}
	http_request.Body = io.NopCloser(bytes.NewReader(body))
	httpClient := &http.Client{}
	proxyReq, err := http.NewRequest(http_request.Method, webhook, bytes.NewReader(body))
	if err != nil {
		return err
	}
	proxyReq.Header = http_request.Header
	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
