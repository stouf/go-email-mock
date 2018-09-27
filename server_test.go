package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

type requestBody struct {
	reader *strings.Reader
	length int
}

func (r requestBody) Close() error {
	return nil
}

func (r requestBody) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

func makeBody(e *Email) (requestBody, error) {
	if e == nil {
		return requestBody{
			reader: strings.NewReader(""),
			length: 0,
		}, nil
	}
	var asJson, err = json.Marshal(e)
	if err != nil {
		return requestBody{}, err
	}
	return requestBody{
		reader: strings.NewReader(string(asJson)),
		length: len(asJson),
	}, nil
}

func getUrl(path string) (url.URL, error) {
	var port, err = readPort()
	if err != nil {
		return url.URL{}, err
	}
	return url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("127.0.0.1:%d", port),
		Path:   path,
	}, nil
}

func send(path string, content *Email) (*http.Response, error) {
	var client = http.Client{}
	var headers = http.Header{}
	headers.Add("Content-Type", "application/json")
	var url, err = getUrl(path)
	if err != nil {
		return nil, err
	}
	body, err := makeBody(content)
	if err != nil {
		return nil, err
	}
	var request = http.Request{
		Method:        "POST",
		URL:           &url,
		Header:        headers,
		Body:          body,
		ContentLength: int64(body.length),
	}
	return client.Do(&request)
}

func TestSendingEmail(T *testing.T) {
	var now = time.Now().UnixNano()
	var sender = fmt.Sprintf("me-%d@test.net", now)
	var senderName = fmt.Sprintf("me-%d", now)
	var recipients = []string{"you@test.net"}
	var subject = fmt.Sprintf("[%d] subject", now)
	var content = fmt.Sprintf("[%d] content", now)
	var requestBody = Email{
		Sender:     sender,
		SenderName: senderName,
		Recipients: recipients,
		Subject:    subject,
		Content:    content,
	}
	response, err := send("/send", &requestBody)
	if err != nil {
		T.Fatal(err)
	}
	if response.StatusCode != 200 {
		T.Fatalf("Received a status code that is not 200: %d", response.StatusCode)
	}

	response, err = send("/get", nil)
	if err != nil {
		T.Fatal(err)
	}
	if response.StatusCode != 200 {
		T.Fatalf("Received a status code that is not 200: %d", response.StatusCode)
	}
	if response.ContentLength <= 0 {
		T.Fatal("No body in the response from the server")
	}
	var rawBody = make([]byte, response.ContentLength)
	_, err = response.Body.Read(rawBody)
	if err != nil && err != io.EOF {
		T.Fatal(err)
	}
	var emails []Email
	err = json.Unmarshal(rawBody, &emails)
	if err != nil {
		T.Fatal(err)
	}

	var found = false
	for _, e := range emails {
		var recipientsMatch = len(e.Recipients) > 0
		for i, r := range e.Recipients {
			found = found && r == recipients[i]
		}
		if !recipientsMatch {
			break
		}
		if e.Sender == sender &&
			e.SenderName == senderName &&
			recipientsMatch &&
			e.Subject == subject &&
			e.Content == content {
			found = true
			break
		}
	}
	if !found {
		T.Fatal("Could not find the email on the server")
	}
}

func TestFlushingEmails(T *testing.T) {
	var now = time.Now().UnixNano()
	var requestBody = Email{
		Sender:     fmt.Sprintf("me-%d@test.net", now),
		SenderName: fmt.Sprintf("me-%d", now),
		Recipients: []string{"you@test.net"},
		Subject:    fmt.Sprintf("A topic [%d]", now),
		Content:    fmt.Sprintf("A word for you [%d]", now),
	}
	response, err := send("/send", &requestBody)
	if err != nil {
		T.Fatal(err)
	}
	if response.StatusCode != 200 {
		T.Fatalf("Received a status code that is not 200: %d", response.StatusCode)
	}

	response, err = send("/flush", nil)
	if err != nil {
		T.Fatal(err)
	}
	if response.StatusCode != 200 {
		T.Fatalf("Received a status code that is not 200: %d", response.StatusCode)
	}

	response, err = send("/get", nil)
	if err != nil {
		T.Fatal(err)
	}
	if response.StatusCode != 200 {
		T.Fatalf("Received a status code that is not 200: %d", response.StatusCode)
	}

	if response.ContentLength <= 0 {
		T.Fatal("No body in the response from the server")
	}
	var rawBody = make([]byte, response.ContentLength)
	_, err = response.Body.Read(rawBody)
	if err != nil && err != io.EOF {
		T.Fatal(err)
	}
	var emails []Email
	err = json.Unmarshal(rawBody, &emails)
	if err != nil {
		T.Fatal(err)
	}

	if len(emails) > 0 {
		T.Fatal("Found emails on the server after flushing")
	}
}

func TestNotFoundRoute(T *testing.T) {
	response, err := send("/foo", nil)
	if err != nil {
		T.Fatal(err)
	}
	if response.StatusCode != 404 {
		T.Fatalf(
			"Was expected a 404 status code but got a %d instead",
			response.StatusCode)
	}
}
