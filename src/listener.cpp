#include "../includes/listener.h"
#include "../includes/HCNetSDK.h"
#include "../includes/alarm_parser.h"
#include <iostream>
#include <mutex>
#include <string>
#include <stdio.h>
#include <queue>
#include <cstring>


static LONG g_listenHandle = -1;

void CALLBACK MessageCallback(LONG lCommand, NET_DVR_ALARMER* pAlarmer, char* pAlarmInfo, DWORD dwBufLen, void* pUser) {
    std::string json;   

    if (lCommand == COMM_ALARM_ACS) {
        auto* acsInfo = reinterpret_cast<NET_DVR_ACS_ALARM_INFO*>(pAlarmInfo);
        json = FormatAcsEvent(lCommand, pAlarmer, acsInfo);
        
    } else {
        char fallback[256];
        snprintf(fallback, sizeof(fallback), "{\"cmd\":%d,\"ip\":\"%s\"}", lCommand, pAlarmer->sDeviceIP);
        json = fallback;
    }    
    hik_enqueue_event(json.c_str());
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
    //std::cout << "enqueue event" << std::endl;
    if (!json) return -1;
    
    std::lock_guard<std::mutex> lock(g_eventQueueMutex);
    g_eventQueue.push(std::string(json));
    return 0;
}

const char* hik_pop_event() {
    //std::cout << "pop event" << std::endl;
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
    //std::cout << "queue size" << std::endl;
    if(g_listenHandle < 0) return -1;
    std::lock_guard<std::mutex> lock(g_eventQueueMutex);
    return static_cast<int>(g_eventQueue.size());
}