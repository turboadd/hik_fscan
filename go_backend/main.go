package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Start Logger
	err := InitLogger("log/app.log", LevelDebug)
	if err != nil {
		fmt.Println("Logging setup failed:", err)
		os.Exit(1)
	}

	// Load Config
	if err := LoadConfig("../config.json"); err != nil {
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

	// // For Test Event
	// go func() {
	// 	for {
	// 		InjectMockEvent("{\"cmd\":1000, \"ip\":\"192.168.1.200\", \"port\":8000, \"mock\":true}")
	// 		time.Sleep(3 * time.Second)
	// 	}
	// }()

	go PollEvents()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Press Ctl+C to stop...")
	<-quit

	Info("Shutting down...")
}
