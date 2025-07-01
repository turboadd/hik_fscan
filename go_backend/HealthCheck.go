package main

import (
	"fmt"
	"net/http"
	"time"
)

func StartHealthCheck() {
	go func() {
		for {
			checkBackend()
			checkQueueStatus()
			time.Sleep(30 * time.Minute)
		}
	}()
}

func checkBackend() {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(AppConfig.EventEndpoint)
	if err != nil || resp.StatusCode >= 400 {
		Error(fmt.Sprintf("[HEALTH] Backend responded with status: %d %s", resp.StatusCode, resp.Status))
		return
	}
	Info("[HEALTH] Backend is reachable")
}

func checkQueueStatus() {
	size := GetQueueSize()
	if size == 0 {
		Warn("[HEALTH] Event queue is empty - no event detected")
	} else {
		Info("[HEALTH] Queue size: " + fmt.Sprint(size))
	}
}

func StartListenerMonitor() {
	go func() {
		for {
			time.Sleep(10 * time.Second) //Time to check Listener

			if GetQueueSize() < 0 {
				Error("[MONITOR] Listener not running. Restarting...")

				if err := StopListening(); err != nil {
					Error("[MONITOR] Failed to stop listener: " + err.Error())
				}
				time.Sleep(1 * time.Second)
				if err := StartListening(AppConfig.ListenPort); err != nil {
					Error("[MONITOR] Failed to restart listener: " + err.Error())
				} else {
					Info("[MONITOR] Listener restarted successfully")
				}
			}
		}
	}()
}
