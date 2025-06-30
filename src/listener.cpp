#include "../includes/listener.h"
#include "../includes/HCNetSDK.h"
#include <iostream>
#include <mutex>
#include <string>
#include <stdio.h>
#include <queue>


static LONG g_listenHandle = -1;
static std::string g_lastEvent;
static std::mutex g_eventMutex;

void CALLBACK MessageCallback(LONG lCommand, NET_DVR_ALARMER* pAlarmer, char* pAlarmInfo, DWORD dwBufLen, void* pUser) {
    std::lock_guard<std::mutex> lock(g_eventMutex);
    char buffer[512];
    snprintf(buffer, sizeof(buffer), "{\"cmd\":%d,\"ip\":\"%s\"}", lCommand, pAlarmer->sDeviceIP);
    hik_enqueue_event(buffer);
}

int hik_start_listening(int port) {
    g_listenHandle = NET_DVR_StartListen_V30(NULL, (WORD)port, MessageCallback, NULL);
    if (g_listenHandle < 0) {
        printf("StartListen failed. Error: %d\n", NET_DVR_GetLastError());
        return -1;
    }    
    return 0;
}

int hik_stop_listening() {
    if (g_listenHandle >= 0) {
        if (!NET_DVR_StopListen_V30(g_listenHandle)) {
            printf("StopListen failed. Error: %d\n", NET_DVR_GetLastError());
            return -1;
        }
        g_listenHandle = -1;
    }
    return 0;
}

// Event Queue
static std::queue<std::string> g_eventQueue;
static std::mutex g_eventQueueMutex;

int hik_enqueue_event(const char* json) {
    if (!json) return -1;

    std::lock_guard<std::mutex> lock(g_eventQueueMutex);
    g_eventQueue.push(std::string(json));
    return 0;
}

const char* hik_pop_event() {
    static std::string lastEvent;
    std::lock_guard<std::mutex> lock(g_eventQueueMutex);

    if (g_eventQueue.empty()) {
        return "";
    }

    lastEvent = g_eventQueue.front();
    g_eventQueue.pop();
    return lastEvent.c_str();
}

int hik_queue_size() {
    std::lock_guard<std::mutex> lock(g_eventQueueMutex);
    return static_cast<int>(g_eventQueue.size());
}

// For Testing Event to GO.
int hik_mock_event(const char* json) {
    std::lock_guard<std::mutex> lock(g_eventMutex);
    g_lastEvent = json;
    return 0;
}