#ifndef LISTENER_H
#define LISTENER_H

#ifdef __cplusplus
extern "C" {
#endif

int hik_start_listening(int port);
int hik_stop_listening();
const char* hik_get_last_event();

//For Test send event to GO
int hik_mock_event(const char* json);

#ifdef __cplusplus
}
#endif

#endif // LISTENER_H

