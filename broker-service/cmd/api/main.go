package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const webPort = "80"

type Config struct{}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	mailServiceURL := "http://mail-service/send"

	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error sending mail with code"+strconv.Itoa(response.StatusCode)))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Mail sent to " + msg.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define the http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the http server
	if err := srv.ListenAndServe(); err != nil {
		log.Panic("server failed to start: %v", err)
	}
}
