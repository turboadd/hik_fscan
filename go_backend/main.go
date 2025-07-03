package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Start Logger
	err := InitLogger("log/edgeservice.log", LevelDebug)
	if err != nil {
		fmt.Println("Logging setup failed:", err)
		os.Exit(1)
	}

	// Load Config
	if err := LoadConfig("config.json"); err != nil {
		Error(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initilization Hikvision SDK
	if err := InitHik(); err != nil {
		Error(fmt.Sprintln("Hikvision SDK initialization failed"))
		return
	}
	defer func() {
		CleanupHik()
		Info("Hikvision SDK cleaned up successfully")
	}()

	Info("Hikvision SDK initialized successfully")

	if err := StartListening(AppConfig.ListenPort); err != nil {
		Error(fmt.Sprintf("Failed to start listener on port %d: %v", AppConfig.ListenPort, err))
	} else {
		Info("Listener started successfully")
		Info(fmt.Sprint("Listening on port: ", AppConfig.ListenPort))
	}
	defer func() {
		StopListening()
		Info("Listener stopped successfully")
	}()

	go func() {
		for {
			fmt.Println("Press enter for mock event..")
			fmt.Scanln()
			currentTime := time.Now().Unix()
			event := fmt.Sprintf(`{"site_id":1234,"ip":"192.168.1.101", "device_id":"0001", "client_id":"0001", "time": "%d"}`, currentTime)
			InjectMockEvent(event)
			//time.Sleep(1 * time.Second)
		}
	}()

	// Pull Evetns from C++
	//ctx, cancel := context.WithCancel(context.Background())
	//go PollEvents(ctx)

	// Start retry worker
	StartRetryWorker()

	// Send event to backend
	StartEventDispatcher()

	// Health check system
	StartHealthCheck()

	// Wait for signal Ctl+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Press Ctl+C to stop...")
	<-quit

	Info("Shutting down...")
	//cancel()
	time.Sleep(300 * time.Millisecond)
}
