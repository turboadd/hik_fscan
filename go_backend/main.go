package main

import (
	"fmt"
)

func main() {
	if err := InitHik(); err != nil {
		fmt.Println("Init Error:", err)
		return
	}
	defer CleanupHik()

	fmt.Println("Hikvision SDK initialized successfully")

	if err := StartListening(8000); err != nil {
		fmt.Println("Start Listening Error:", err)
	}
	defer StopListening()

	// // For Test Event
	// go func() {
	// 	for {
	// 		InjectMockEvent("{\"cmd\":1000, \"ip\":\"192.168.1.200\", \"port\":8000, \"mock\":true}")
	// 		time.Sleep(3 * time.Second)
	// 	}
	// }()

	go PollEvents()

	fmt.Println("Press Ctl+C to stop...")
	select {}

	// ต่อไปอาจ startListen หรือ run logic อื่น ๆ
}
