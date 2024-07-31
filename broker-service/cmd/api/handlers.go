package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Broker service is running",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	if err := app.readJSON(w, r, &requestPayload); err != nil {
		if err := app.errorJSON(w, err); err != nil {
			log.Printf("error writing response %v", err)
		}
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		err := app.errorJSON(w, errors.New("invalid action"))
		if err != nil {
			log.Printf("error writing response %v", err)
		}

	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("error writing response %v", err)
		}
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("error writing response %v", err)
		}
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if response.StatusCode == http.StatusUnauthorized {
		err := app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		if err != nil {
			log.Printf("error writing response %v", err)
		}
		return
	} else if response.StatusCode != http.StatusAccepted {
		err := app.errorJSON(w, errors.New(fmt.Sprintf("error calling auth service %d", response.StatusCode)), http.StatusInternalServerError)
		if err != nil {
			log.Printf("error writing response %v", err)
		}
		return
	}

	var authResponse jsonResponse
	if err := json.NewDecoder(response.Body).Decode(&authResponse); err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			log.Printf("error writing response %v", err)
		}
		return
	}

	if authResponse.Error {
		err := app.errorJSON(w, err, http.StatusUnauthorized)
		if err != nil {
			log.Printf("error writing response %v", err)
		}
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    authResponse.Data,
	}

	err = app.writeJSON(w, http.StatusAccepted, payload)
	if err != nil {
		log.Printf("error writing response %v", err)
	}
}
