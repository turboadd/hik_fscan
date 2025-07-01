#include "../includes/alarm_parser.h"
#include <cstdio>
#include <cstring>

std::string FormatAcsEvent(int lCommand, const NET_DVR_ALARMER* pAlarmer, const NET_DVR_ACS_ALARM_INFO* acsInfo) {
    char json[1024] = {0};

    // 1. แปลงเวลา
    const NET_DVR_TIME& t = acsInfo->struTime;
    char timestamp[64];
    snprintf(timestamp, sizeof(timestamp), "%04d-%02d-%02dT%02d:%02d:%02d",
             t.dwYear, t.dwMonth, t.dwDay, t.dwHour, t.dwMinute, t.dwSecond);

    // 2. ข้อมูลจาก struAcsEventInfo
    const NET_DVR_ACS_EVENT_INFO& info = acsInfo->struAcsEventInfo;

    // cardNo เป็น string
    char cardNo[ACS_CARD_NO_LEN + 1] = {0};
    strncpy(cardNo, reinterpret_cast<const char*>(info.byCardNo), ACS_CARD_NO_LEN);

    // dwEmployeeNo เป็นตัวเลข → แปลงเป็น string
    char userID[32] = {0};
    snprintf(userID, sizeof(userID), "%u", info.dwEmployeeNo);
    
    snprintf(json, sizeof(json),
        "{\"cmd\":%d,\"type\":\"ACCESS_GRANTED\",\"deviceIP\":\"%s\","
        "\"userID\":\"%s\",\"cardNo\":\"%s\",\"SN\":%d,"
        "\"eventType\":%d,\"method\":%d,\"doorNo\":%d,"
        "\"timestamp\":\"%s\"}",
        lCommand,
        pAlarmer->sDeviceIP,
        userID,
        cardNo,
        info.dwSerialNo,
        acsInfo->dwMinor,
        info.dwVerifyNo,  // ใช้ verifyNo เป็น method
        info.dwDoorNo,
        timestamp
    );
    return std::string(json);
}