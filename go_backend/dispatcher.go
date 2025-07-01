package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func StartEventDispatcher() {
	go func() {
		for {
			event := PopEvent()
			if event != "" {
				//Info("[DISPATCH] Sending event to backend: " + event)
				SendEventWithRetry(enrichEvent(event))
				fmt.Println("[DISPATCH] Sent event to backend: " + event)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func enrichEvent(jsonStr string) string {
	var event map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		return jsonStr
	}
	ip := event["ip"].(string)
	if name, found := AppConfig.Devices[ip]; found {
		event["deviceName"] = name
	} else {
		event["deviceName"] = "Unknown"
	}

	event["eventID"] = uuid.New().String()

	bytes, err := json.Marshal(event)
	if err != nil {
		Error("Failed to re-marshal JSON: " + err.Error())
		return jsonStr
	}
	return string(bytes)
}
