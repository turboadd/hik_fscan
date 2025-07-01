#ifndef ALARM_PARSER_H
#define ALARM_PARSER_H

#include "../includes/HCNetSDK.h"
#include <string>

std::string FormatAcsEvent(int lCommand, const NET_DVR_ALARMER* pAlarmer, const NET_DVR_ACS_ALARM_INFO* acsInfo);

#endif // ALARM_PARSER_H