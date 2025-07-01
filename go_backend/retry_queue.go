package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var retryQueue = make([]string, 0)
var retryMutex sync.Mutex

// postNow tries to POST the event immediately.
func postNow(json string) error {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Post(AppConfig.EventEndpoint, "application/json", bytes.NewBufferString(json))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("backend returned status code: %s", resp.Status)
	}
	fmt.Println("Sent event with status: " + resp.Status)
	return nil
}

// SendEventWithRetry tries to send, or pushes to retry queue on failure.
func SendEventWithRetry(json string) {
	err := postNow(json)
	if err != nil {
		Error("Failed to POST event, retrying later: " + err.Error())
		retryMutex.Lock()
		retryQueue = append(retryQueue, json)
		retryMutex.Unlock()
	}
}

// StartRetryWorker continuously retries failed events.
func StartRetryWorker() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			retryMutex.Lock()
			if len(retryQueue) == 0 {
				retryMutex.Unlock()
				continue
			}
			event := retryQueue[0]
			retryQueue = retryQueue[1:]
			retryMutex.Unlock()

			Error("Retrying event: " + event)
			SendEventWithRetry(event)
		}
	}()
}

func generateJWT()(string, error) {
	claims := jwt.MapClaims{
		"site": "siteA",
		"exp" : time.now().Add(24 * time.Hour).Unix(),
		"iat" : time.now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.signedString([]byte(AppConfig.AuthToken))
}
