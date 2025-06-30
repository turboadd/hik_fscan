#include "../includes/hik_adapter.h"
#include "../includes/HCNetSDK.h"

int hik_init() {
    if (!NET_DVR_Init()) {
        return -1;
    }
    return 0;
}

int hik_cleanup() {
    NET_DVR_Cleanup();
    return 0;
}