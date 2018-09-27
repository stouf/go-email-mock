package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Email struct {
	Sender          string
	SenderName      string
	Recipients      []string
	CarbonCopy      []string
	BlindCarbonCopy []string
	Subject         string
	Content         string
}

type HttpHandler struct{}

var emails = make([]Email, 0)

func notFound(response http.ResponseWriter) {
	response.WriteHeader(http.StatusNotFound)
	response.Write(nil)
}

func receiveEmail(response http.ResponseWriter, request *http.Request) {
	if request.ContentLength == -1 {
		log.Println("Unexpected content-length header")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(nil))
		return
	}
	if request.ContentLength == 0 {
		log.Println("Empty body (was expecting one)")
		response.WriteHeader(http.StatusUnprocessableEntity)
		response.Write([]byte("MISSING_BODY"))
		return
	}

	var body = make([]byte, request.ContentLength)
	var _, err = request.Body.Read(body)
	if err != nil && err != io.EOF {
		log.Printf("Error while reading the request body: %v", err)
		response.WriteHeader(http.StatusBadRequest)
		response.Write(nil)
		return
	}

	var email Email
	err = json.Unmarshal(body, &email)
	if err != nil {
		log.Printf("Could not parse the body as a JSON: %v", err)
		response.WriteHeader(http.StatusBadRequest)
		response.Write(nil)
		return
	}
	emails = append(emails, email)

	response.WriteHeader(http.StatusOK)
	response.Write(nil)
	log.Println("Successfully received email")
}

func returnEmails(response http.ResponseWriter, request *http.Request) {
	var asBytes, err = json.Marshal(emails)
	if err != nil {
		log.Printf("Could not json.Marshal the emails: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(nil)
		return
	}
	response.Header().Add("content-type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(asBytes)
	log.Println("Current list of all the received emails returned to client")
}

func flushEmails(response http.ResponseWriter, request *http.Request) {
	emails = make([]Email, 0)
	response.WriteHeader(http.StatusOK)
	response.Write(nil)
	log.Println("All the emails have been flushed")
}

func (handler HttpHandler) ServeHTTP(
	response http.ResponseWriter,
	request *http.Request) {
	defer request.Body.Close()
	if request.Method != "POST" {
		notFound(response)
	} else {
		switch request.URL.Path {
		case "/send":
			receiveEmail(response, request)
		case "/get":
			returnEmails(response, request)
		case "/flush":
			flushEmails(response, request)
		default:
			notFound(response)
		}
	}
}

func Make(port int) http.Server {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: HttpHandler{},
	}
}
