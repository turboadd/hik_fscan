// bridge.go
package main

/*
#cgo CFLAGS: -I../includes
#cgo LDFLAGS: -L../libs -lhik_adapter -lHCNetSDK -lstdc++
#include "../includes/hik_adapter.h"
#include "../includes/listener.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

func InitHik() error {
	if C.hik_init() != 0 {
		Error(fmt.Sprintln("hikvision SDK initialization failed"))
		os.Exit(1)
	}
	return nil
}

func CleanupHik() {
	C.hik_cleanup()
}

func StartListening(port int) error {
	if C.hik_start_listening(C.int(port)) != 0 {
		Error(fmt.Sprintln("listener failed to start"))
		os.Exit(1)
	}
	return nil
}

func StopListening() error {
	if C.hik_stop_listening() != 0 {
		return fmt.Errorf("listener failed to stop")
	}
	return nil
}

func PopEvent() string {
	return C.GoString(C.hik_pop_event())
}

func GetQueueSize() int {
	return int(C.hik_queue_size())
}

// For Test Event
func InjectMockEvent(json string) {
	cstr := C.CString(json)
	defer C.free(unsafe.Pointer(cstr))
	C.hik_enqueue_event(cstr)
}
