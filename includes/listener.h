#ifndef LISTENER_H
#define LISTENER_H

#ifdef __cplusplus
extern "C" {
#endif

int hik_start_listening(int port);
int hik_stop_listening();
const char* hik_get_last_event();

// Event Queue
int hik_enqueue_event(const char* json);
const char* hik_pop_event();
int hik_queue_size();

//For Test send event to GO
int hik_mock_event(const char* json);

#ifdef __cplusplus
}
#endif

#endif // LISTENER_H 

